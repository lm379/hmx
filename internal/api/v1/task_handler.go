package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/queue"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleGetBatchTaskStatus (GET /api/v1/admin/tasks/batch/:id)
func HandleGetBatchTaskStatus(c *gin.Context) {
	batchID := c.Param("id")

	batchTask, err := queue.GetBatchTask(batchID)
	if err != nil {
		resp.NotFound(c, "Batch task not found")
		return
	}

	resp.Success(c, batchTask)
}

// HandleGetTaskStatus (GET /api/v1/admin/tasks/:id)
func HandleGetTaskStatus(c *gin.Context) {
	taskID := c.Param("id")

	task, err := queue.GetTask(taskID)
	if err != nil {
		resp.NotFound(c, "Task not found")
		return
	}

	resp.Success(c, task)
}

// HandleGetQueueStatus (GET /api/v1/admin/tasks/queue/status)
func HandleGetQueueStatus(c *gin.Context) {
	embeddingLength, _ := queue.GetQueueLength(queue.TaskTypeEmbedding)
	summaryLength, _ := queue.GetQueueLength(queue.TaskTypeSummary)
	transcodeLength, _ := queue.GetQueueLength(queue.TaskTypeTranscode)
	subtitleLength, _ := queue.GetQueueLength(queue.TaskTypeSubtitle)
	coverLength, _ := queue.GetQueueLength(queue.TaskTypeCover)

	// 获取等待中的任务
	pendingEmbedding, _ := queue.GetPendingTasks(queue.TaskTypeEmbedding)
	pendingSummary, _ := queue.GetPendingTasks(queue.TaskTypeSummary)
	pendingTranscode, _ := queue.GetPendingTasks(queue.TaskTypeTranscode)
	pendingSubtitle, _ := queue.GetPendingTasks(queue.TaskTypeSubtitle)
	pendingCover, _ := queue.GetPendingTasks(queue.TaskTypeCover)

	// 获取处理中的任务
	processingEmbedding, _ := queue.GetProcessingTasks(queue.TaskTypeEmbedding)
	processingSummary, _ := queue.GetProcessingTasks(queue.TaskTypeSummary)
	processingTranscode, _ := queue.GetProcessingTasks(queue.TaskTypeTranscode)
	processingSubtitle, _ := queue.GetProcessingTasks(queue.TaskTypeSubtitle)
	processingCover, _ := queue.GetProcessingTasks(queue.TaskTypeCover)

	// 获取已完成的任务（最近100个）
	completedEmbedding, _ := queue.GetCompletedTasks(queue.TaskTypeEmbedding)
	completedSummary, _ := queue.GetCompletedTasks(queue.TaskTypeSummary)
	completedTranscode, _ := queue.GetCompletedTasks(queue.TaskTypeTranscode)
	completedSubtitle, _ := queue.GetCompletedTasks(queue.TaskTypeSubtitle)
	completedCover, _ := queue.GetCompletedTasks(queue.TaskTypeCover)

	resp.Success(c, gin.H{
		"embedding_queue_length":     embeddingLength,
		"summary_queue_length":       summaryLength,
		"transcode_queue_length":     transcodeLength,
		"subtitle_queue_length":      subtitleLength,
		"cover_queue_length":         coverLength,
		"pending_embedding_tasks":    pendingEmbedding,
		"pending_summary_tasks":      pendingSummary,
		"pending_transcode_tasks":    pendingTranscode,
		"pending_subtitle_tasks":     pendingSubtitle,
		"pending_cover_tasks":        pendingCover,
		"processing_embedding_tasks": processingEmbedding,
		"processing_summary_tasks":   processingSummary,
		"processing_transcode_tasks": processingTranscode,
		"processing_subtitle_tasks":  processingSubtitle,
		"processing_cover_tasks":     processingCover,
		"completed_embedding_tasks":  completedEmbedding,
		"completed_summary_tasks":    completedSummary,
		"completed_transcode_tasks":  completedTranscode,
		"completed_subtitle_tasks":   completedSubtitle,
		"completed_cover_tasks":      completedCover,
	})
}
