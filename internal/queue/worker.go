package queue

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/ai"
)

// Worker 任务处理器
type Worker struct {
	taskType TaskType
	stopChan chan struct{}
	wg       *sync.WaitGroup
}

// NewWorker 创建新的Worker
func NewWorker(taskType TaskType) *Worker {
	return &Worker{
		taskType: taskType,
		stopChan: make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}
}

// Start 启动Worker
func (w *Worker) Start() {
	w.wg.Add(1)
	go w.run()
	log.Printf("[Worker] %s worker started", w.taskType)
}

// Stop 停止Worker
func (w *Worker) Stop() {
	close(w.stopChan)
	w.wg.Wait()
	log.Printf("[Worker] %s worker stopped", w.taskType)
}

// run Worker运行循环
func (w *Worker) run() {
	defer w.wg.Done()

	for {
		select {
		case <-w.stopChan:
			return
		default:
			w.processTask()
		}
	}
}

// processTask 处理单个任务
func (w *Worker) processTask() {
	task, err := DequeueTask(w.taskType)
	if err != nil {
		log.Printf("[Worker] Failed to dequeue task: %v", err)
		return
	}

	if task == nil {
		return // 队列为空，继续轮询
	}

	// 更新任务状态为处理中
	task.Status = TaskStatusProcessing
	if err := UpdateTask(task); err != nil {
		log.Printf("[Worker] Failed to update task status: %v", err)
		return
	}

	log.Printf("[Worker] Processing %s task for Opera ID %d", task.Type, task.OperaID)

	// 执行任务
	var taskErr error
	switch task.Type {
	case TaskTypeEmbedding:
		taskErr = w.processEmbeddingTask(task)
	case TaskTypeSummary:
		taskErr = w.processSummaryTask(task)
	default:
		taskErr = nil
	}

	// 更新任务状态
	if taskErr != nil {
		task.Status = TaskStatusFailed
		task.Error = taskErr.Error()
		log.Printf("[Worker] Task failed for Opera ID %d: %v", task.OperaID, taskErr)
	} else {
		task.Status = TaskStatusCompleted
		log.Printf("[Worker] Task completed for Opera ID %d", task.OperaID)
	}

	if err := UpdateTask(task); err != nil {
		log.Printf("[Worker] Failed to update task final status: %v", err)
	}

	// 更新批量任务统计
	w.updateBatchTask(task)
}

// processEmbeddingTask 处理向量生成任务
func (w *Worker) processEmbeddingTask(task *Task) error {
	operaRepo := repository.NewOperaRepo()
	opera, err := operaRepo.GetByID(task.OperaID)
	if err != nil {
		return err
	}

	// 过滤：如果已有向量且不是强制生成，直接跳过
	slice := opera.Embedding.Slice()
	if len(slice) > 0 && !task.Force {
		log.Printf("[Worker] Opera ID %d already has embedding, skipping", task.OperaID)
		return nil // 不算错误，让任务成功完成
	}

	// 构建文本：标题 + 描述 + AI摘要
	text := opera.OperaTitle
	if opera.Description != "" {
		text += " " + opera.Description
	}
	if opera.AiSummary != "" {
		text += " " + opera.AiSummary
	}

	// 调用 Embedding API 生成向量
	embedding, err := ai.GenerateEmbedding(text)
	if err != nil {
		return err
	}

	// 保存向量
	return operaRepo.UpdateEmbedding(opera.OperaID, embedding)
}

// processSummaryTask 处理摘要生成任务
func (w *Worker) processSummaryTask(task *Task) error {
	operaRepo := repository.NewOperaRepo()
	opera, err := operaRepo.GetByID(task.OperaID)
	if err != nil {
		return err
	}

	// 过滤：如果已有摘要且不是强制生成，直接跳过
	if opera.AiSummary != "" && !task.Force {
		log.Printf("[Worker] Opera ID %d already has summary, skipping", task.OperaID)
		return nil // 不算错误，让任务成功完成
	}

	// 检查是否有字幕文件
	if !opera.SrtPath.Valid || opera.SrtPath.String == "" {
		log.Printf("[Worker] Opera ID %d has no subtitle, skipping", task.OperaID)
		return nil // 没有字幕文件，跳过
	}

	// 从CDN读取字幕
	subtitleContent, err := readSubtitleFromCDN(opera.SrtPath.String)
	if err != nil {
		return err
	}

	// 生成摘要
	summary, err := ai.GenerateSubtitleSummary(opera.OperaTitle, subtitleContent)
	if err != nil {
		return err
	}

	// 保存摘要
	return operaRepo.UpdateSummary(opera.OperaID, summary)
}

// readSubtitleFromCDN 从CDN/对象存储读取字幕文件内容
func readSubtitleFromCDN(srtPath string) (string, error) {
	// 构建完整的CDN URL
	cdnURL := config.AppConfig.S3CustomDomain
	if cdnURL == "" {
		// 如果没有自定义域名，使用S3 endpoint
		cdnURL = config.AppConfig.S3Endpoint + "/" + config.AppConfig.S3Bucket
	}

	// 移除srtPath开头的斜杠（如果有）
	if len(srtPath) > 0 && srtPath[0] == '/' {
		srtPath = srtPath[1:]
	}

	fullURL := cdnURL + "/" + srtPath

	// 发起HTTP GET请求读取字幕文件
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch subtitle file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch subtitle, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read subtitle content: %w", err)
	}

	return string(body), nil
}

// updateBatchTask 更新批量任务统计
func (w *Worker) updateBatchTask(task *Task) {
	// 从任务ID中提取批量任务ID（格式：batchID-operaID）
	batchID := ""
	for i := len(task.ID) - 1; i >= 0; i-- {
		if task.ID[i] == '-' {
			batchID = task.ID[:i]
			break
		}
	}

	if batchID == "" {
		return // 不是批量任务
	}

	// 获取批量任务
	batchTask, err := GetBatchTask(batchID)
	if err != nil {
		return
	}

	// 更新统计
	batchTask.Pending--
	switch task.Status {
	case TaskStatusCompleted:
		batchTask.Success++
	case TaskStatusFailed:
		batchTask.Failed++
	}

	// 保存批量任务
	if err := UpdateBatchTask(batchTask); err != nil {
		log.Printf("[Worker] Failed to update batch task: %v", err)
	}
}

// StartWorkers 启动所有Worker
func StartWorkers() []*Worker {
	workers := []*Worker{
		NewWorker(TaskTypeEmbedding),
		NewWorker(TaskTypeSummary),
	}

	for _, worker := range workers {
		worker.Start()
	}

	return workers
}

// StopWorkers 停止所有Worker
func StopWorkers(workers []*Worker) {
	for _, worker := range workers {
		worker.Stop()
	}
}
