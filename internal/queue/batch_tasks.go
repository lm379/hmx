package queue

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// BatchGenerateEmbeddingsAsync 批量加入向量生成任务到队列
func BatchGenerateEmbeddingsAsync(operaIDs []uint, force bool) (string, error) {
	if len(operaIDs) == 0 {
		return "", fmt.Errorf("no opera IDs provided")
	}

	// 创建批量任务
	batchID := uuid.New().String()
	batchTask := &BatchTask{
		ID:        batchID,
		Type:      TaskTypeEmbedding,
		Total:     len(operaIDs),
		Success:   0,
		Failed:    0,
		Pending:   len(operaIDs),
		Status:    TaskStatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := CreateBatchTask(batchTask); err != nil {
		return "", err
	}

	// 将每个作品加入队列
	for _, operaID := range operaIDs {
		task := &Task{
			ID:      fmt.Sprintf("%s-%d", batchID, operaID),
			Type:    TaskTypeEmbedding,
			OperaID: operaID,
			Force:   force,
		}

		if err := EnqueueTask(task); err != nil {
			log.Printf("[Batch Embedding] Failed to enqueue task for Opera ID %d: %v", operaID, err)
			continue
		}
	}

	log.Printf("[Batch Embedding] Batch task %s created with %d operas (force=%v)", batchID, len(operaIDs), force)
	return batchID, nil
}

// BatchGenerateOperaSummariesAsync 批量加入摘要生成任务到队列
func BatchGenerateOperaSummariesAsync(operaIDs []uint, force bool) (string, error) {
	if len(operaIDs) == 0 {
		return "", fmt.Errorf("no opera IDs provided")
	}

	// 创建批量任务
	batchID := uuid.New().String()
	batchTask := &BatchTask{
		ID:        batchID,
		Type:      TaskTypeSummary,
		Total:     len(operaIDs),
		Success:   0,
		Failed:    0,
		Pending:   len(operaIDs),
		Status:    TaskStatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := CreateBatchTask(batchTask); err != nil {
		return "", err
	}

	// 将每个作品加入队列
	for _, operaID := range operaIDs {
		task := &Task{
			ID:      fmt.Sprintf("%s-%d", batchID, operaID),
			Type:    TaskTypeSummary,
			OperaID: operaID,
			Force:   force,
		}

		if err := EnqueueTask(task); err != nil {
			log.Printf("[Batch Summary] Failed to enqueue task for Opera ID %d: %v", operaID, err)
			continue
		}
	}

	log.Printf("[Batch Summary] Batch task %s created with %d operas (force=%v)", batchID, len(operaIDs), force)
	return batchID, nil
}
