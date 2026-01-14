package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lm379/hmx/database"
)

var (
	ctx = context.Background()
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeEmbedding TaskType = "embedding"
	TaskTypeSummary   TaskType = "summary"
	TaskTypeTranscode TaskType = "transcode" // 视频转码任务
	TaskTypeSubtitle  TaskType = "subtitle"  // 字幕生成任务
	TaskTypeCover     TaskType = "cover"     // 封面生成任务
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

// Task 任务结构
type Task struct {
	ID        string     `json:"id"`
	Type      TaskType   `json:"type"`
	OperaID   uint       `json:"opera_id"`
	Force     bool       `json:"force"` // 是否强制重新生成
	Status    TaskStatus `json:"status"`
	Error     string     `json:"error,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// BatchTask 批量任务
type BatchTask struct {
	ID        string     `json:"id"`
	Type      TaskType   `json:"type"`
	Total     int        `json:"total"`
	Success   int        `json:"success"`
	Failed    int        `json:"failed"`
	Pending   int        `json:"pending"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

const (
	// 队列名称
	QueueEmbedding = "queue:embedding"
	QueueSummary   = "queue:summary"
	QueueTranscode = "queue:transcode"
	QueueSubtitle  = "queue:subtitle"
	QueueCover     = "queue:cover"

	// 任务状态键前缀
	TaskKeyPrefix      = "task:"
	BatchTaskKeyPrefix = "batch_task:"
	TempTaskKeyPrefix  = "temp_task:" // 临时任务前缀（用于暂存提前到达的回调）
)

// EnqueueTask 将任务加入队列
func EnqueueTask(task *Task) error {
	task.Status = TaskStatusPending
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// 保存任务详情
	taskKey := TaskKeyPrefix + task.ID
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if err := database.RDBQueue.Set(ctx, taskKey, taskJSON, 7*24*time.Hour).Err(); err != nil {
		return err
	}

	// 加入队列
	queueName := getQueueName(task.Type)
	return database.RDBQueue.LPush(ctx, queueName, task.ID).Err()
}

// DequeueTask 从队列取出任务
func DequeueTask(taskType TaskType) (*Task, error) {
	queueName := getQueueName(taskType)

	// 阻塞式弹出（超时5秒）
	result, err := database.RDBQueue.BRPop(ctx, 5*time.Second, queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 队列为空
		}
		return nil, err
	}

	if len(result) < 2 {
		return nil, nil
	}

	taskID := result[1]
	return GetTask(taskID)
}

// GetTask 获取任务详情
func GetTask(taskID string) (*Task, error) {
	taskKey := TaskKeyPrefix + taskID
	taskJSON, err := database.RDBQueue.Get(ctx, taskKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	var task Task
	if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTask 更新任务状态
func UpdateTask(task *Task) error {
	task.UpdatedAt = time.Now()
	taskKey := TaskKeyPrefix + task.ID
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return database.RDBQueue.Set(ctx, taskKey, taskJSON, 7*24*time.Hour).Err()
}

// CreateBatchTask 创建批量任务
func CreateBatchTask(batchTask *BatchTask) error {
	batchTask.Status = TaskStatusProcessing
	batchTask.CreatedAt = time.Now()
	batchTask.UpdatedAt = time.Now()

	batchKey := BatchTaskKeyPrefix + batchTask.ID
	batchJSON, err := json.Marshal(batchTask)
	if err != nil {
		return err
	}

	return database.RDBQueue.Set(ctx, batchKey, batchJSON, 7*24*time.Hour).Err()
}

// GetBatchTask 获取批量任务详情
func GetBatchTask(batchID string) (*BatchTask, error) {
	batchKey := BatchTaskKeyPrefix + batchID
	batchJSON, err := database.RDBQueue.Get(ctx, batchKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("batch task not found")
		}
		return nil, err
	}

	var batchTask BatchTask
	if err := json.Unmarshal([]byte(batchJSON), &batchTask); err != nil {
		return nil, err
	}

	return &batchTask, nil
}

// UpdateBatchTask 更新批量任务
func UpdateBatchTask(batchTask *BatchTask) error {
	batchTask.UpdatedAt = time.Now()

	// 更新状态
	if batchTask.Pending == 0 {
		if batchTask.Failed == 0 {
			batchTask.Status = TaskStatusCompleted
		} else if batchTask.Success == 0 {
			batchTask.Status = TaskStatusFailed
		} else {
			batchTask.Status = TaskStatusCompleted
		}
	}

	batchKey := BatchTaskKeyPrefix + batchTask.ID
	batchJSON, err := json.Marshal(batchTask)
	if err != nil {
		return err
	}

	return database.RDBQueue.Set(ctx, batchKey, batchJSON, 7*24*time.Hour).Err()
}

// getQueueName 根据任务类型获取队列名称
func getQueueName(taskType TaskType) string {
	switch taskType {
	case TaskTypeEmbedding:
		return QueueEmbedding
	case TaskTypeSummary:
		return QueueSummary
	case TaskTypeTranscode:
		return QueueTranscode
	case TaskTypeSubtitle:
		return QueueSubtitle
	case TaskTypeCover:
		return QueueCover
	default:
		return ""
	}
}

// GetQueueLength 获取队列长度
func GetQueueLength(taskType TaskType) (int64, error) {
	queueName := getQueueName(taskType)
	return database.RDBQueue.LLen(ctx, queueName).Result()
}

// GetProcessingTasks 获取所有正在处理的任务（通过扫描任务键）
func GetProcessingTasks(taskType TaskType) ([]*Task, error) {
	pattern := "task:*"

	var tasks []*Task
	iter := database.RDBQueue.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		taskJSON, err := database.RDBQueue.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var task Task
		if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
			continue
		}

		// 过滤任务类型和状态
		if task.Type == taskType && task.Status == TaskStatusProcessing {
			tasks = append(tasks, &task)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetPendingTasks 获取队列中等待的任务（但不出队）
func GetPendingTasks(taskType TaskType) ([]*Task, error) {
	queueName := getQueueName(taskType)

	// 使用 LRANGE 获取所有任务（0到-1表示全部）
	taskIDs, err := database.RDBQueue.LRange(ctx, queueName, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for _, taskID := range taskIDs {
		task, err := GetTask(taskID)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetCompletedTasks 获取已完成的任务（最近100个）
func GetCompletedTasks(taskType TaskType) ([]*Task, error) {
	pattern := "task:*"

	var tasks []*Task
	iter := database.RDBQueue.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		taskJSON, err := database.RDBQueue.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var task Task
		if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
			continue
		}

		// 过滤任务类型和状态
		if task.Type == taskType && (task.Status == TaskStatusCompleted || task.Status == TaskStatusFailed) {
			tasks = append(tasks, &task)
		}

		// 限制最多返回100个
		if len(tasks) >= 100 {
			break
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// SaveTempTask 保存临时任务（用于暂存提前到达的回调，TTL 5分钟）
func SaveTempTask(taskID string) error {
	tempKey := TempTaskKeyPrefix + taskID
	// 只需要标记存在即可，值可以是简单的时间戳
	return database.RDBQueue.Set(ctx, tempKey, time.Now().Unix(), 5*time.Minute).Err()
}

// GetTempTask 检查临时任务是否存在
func GetTempTask(taskID string) (bool, error) {
	tempKey := TempTaskKeyPrefix + taskID
	_, err := database.RDBQueue.Get(ctx, tempKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 不存在
		}
		return false, err
	}
	return true, nil // 存在
}

// DeleteTempTask 删除临时任务
func DeleteTempTask(taskID string) error {
	tempKey := TempTaskKeyPrefix + taskID
	return database.RDBQueue.Del(ctx, tempKey).Err()
}
