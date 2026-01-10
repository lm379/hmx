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

	// 获取等待中的任务
	pendingEmbedding, _ := queue.GetPendingTasks(queue.TaskTypeEmbedding)
	pendingSummary, _ := queue.GetPendingTasks(queue.TaskTypeSummary)

	// 获取处理中的任务
	processingEmbedding, _ := queue.GetProcessingTasks(queue.TaskTypeEmbedding)
	processingSummary, _ := queue.GetProcessingTasks(queue.TaskTypeSummary)

	// 获取已完成的任务（最近100个）
	completedEmbedding, _ := queue.GetCompletedTasks(queue.TaskTypeEmbedding)
	completedSummary, _ := queue.GetCompletedTasks(queue.TaskTypeSummary)

	resp.Success(c, gin.H{
		"embedding_queue_length":     embeddingLength,
		"summary_queue_length":       summaryLength,
		"pending_embedding_tasks":    pendingEmbedding,
		"pending_summary_tasks":      pendingSummary,
		"processing_embedding_tasks": processingEmbedding,
		"processing_summary_tasks":   processingSummary,
		"completed_embedding_tasks":  completedEmbedding,
		"completed_summary_tasks":    completedSummary,
	})
}
