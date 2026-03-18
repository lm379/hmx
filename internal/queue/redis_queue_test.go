package queue

import (
	"testing"
	"time"
)

func TestTaskType(t *testing.T) {
	tests := []struct {
		name  string
		typ   TaskType
		valid bool
	}{
		{"Embedding type", TaskTypeEmbedding, true},
		{"Summary type", TaskTypeSummary, true},
		{"Transcode type", TaskTypeTranscode, true},
		{"Subtitle type", TaskTypeSubtitle, true},
		{"Cover type", TaskTypeCover, true},
		{"DocumentEmbedding type", TaskTypeDocumentEmbedding, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid && tt.typ == "" {
				t.Error("Expected non-empty TaskType")
			}
		})
	}
}

func TestTaskStatus(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		valid  bool
	}{
		{"Pending status", TaskStatusPending, true},
		{"Processing status", TaskStatusProcessing, true},
		{"Completed status", TaskStatusCompleted, true},
		{"Failed status", TaskStatusFailed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid && tt.status == "" {
				t.Error("Expected non-empty TaskStatus")
			}
		})
	}
}

func TestTaskStruct(t *testing.T) {
	now := time.Now()

	task := Task{
		ID:        "task-123",
		Type:      TaskTypeEmbedding,
		OperaID:   1,
		DocID:     "doc-456",
		Force:     true,
		Status:    TaskStatusPending,
		Error:     "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if task.ID != "task-123" {
		t.Errorf("ID = %s, want task-123", task.ID)
	}
	if task.Type != TaskTypeEmbedding {
		t.Errorf("Type = %v, want TaskTypeEmbedding", task.Type)
	}
	if task.OperaID != 1 {
		t.Errorf("OperaID = %d, want 1", task.OperaID)
	}
	if task.DocID != "doc-456" {
		t.Errorf("DocID = %s, want doc-456", task.DocID)
	}
	if !task.Force {
		t.Error("Expected Force to be true")
	}
	if task.Status != TaskStatusPending {
		t.Errorf("Status = %v, want TaskStatusPending", task.Status)
	}
	if task.Error != "" {
		t.Errorf("Expected empty error, got %s", task.Error)
	}
}

func TestBatchTaskStruct(t *testing.T) {
	now := time.Now()

	batchTask := BatchTask{
		ID:        "batch-123",
		Type:      TaskTypeEmbedding,
		Total:     100,
		Success:   95,
		Failed:    5,
		Pending:   0,
		Status:    TaskStatusCompleted,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if batchTask.ID != "batch-123" {
		t.Errorf("ID = %s, want batch-123", batchTask.ID)
	}
	if batchTask.Type != TaskTypeEmbedding {
		t.Errorf("Type = %v, want TaskTypeEmbedding", batchTask.Type)
	}
	if batchTask.Total != 100 {
		t.Errorf("Total = %d, want 100", batchTask.Total)
	}
	if batchTask.Success != 95 {
		t.Errorf("Success = %d, want 95", batchTask.Success)
	}
	if batchTask.Failed != 5 {
		t.Errorf("Failed = %d, want 5", batchTask.Failed)
	}
	if batchTask.Pending != 0 {
		t.Errorf("Pending = %d, want 0", batchTask.Pending)
	}
	if batchTask.Status != TaskStatusCompleted {
		t.Errorf("Status = %v, want TaskStatusCompleted", batchTask.Status)
	}
}

func TestQueueConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Embedding queue", QueueEmbedding, "queue:embedding"},
		{"Summary queue", QueueSummary, "queue:summary"},
		{"Transcode queue", QueueTranscode, "queue:transcode"},
		{"Subtitle queue", QueueSubtitle, "queue:subtitle"},
		{"Cover queue", QueueCover, "queue:cover"},
		{"DocumentEmbedding queue", QueueDocumentEmbedding, "queue:document_embedding"},
		{"Task key prefix", TaskKeyPrefix, "task:"},
		{"BatchTask key prefix", BatchTaskKeyPrefix, "batch_task:"},
		{"TempTask key prefix", TempTaskKeyPrefix, "temp_task:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant = %s, want %s", tt.constant, tt.expected)
			}
		})
	}
}

func TestGetQueueName(t *testing.T) {
	tests := []struct {
		name     string
		taskType TaskType
		expected string
	}{
		{"Embedding", TaskTypeEmbedding, QueueEmbedding},
		{"Summary", TaskTypeSummary, QueueSummary},
		{"Transcode", TaskTypeTranscode, QueueTranscode},
		{"Subtitle", TaskTypeSubtitle, QueueSubtitle},
		{"Cover", TaskTypeCover, QueueCover},
		{"DocumentEmbedding", TaskTypeDocumentEmbedding, QueueDocumentEmbedding},
		{"Unknown type", TaskType("unknown"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getQueueName(tt.taskType)
			if result != tt.expected {
				t.Errorf("getQueueName(%v) = %s, want %s", tt.taskType, result, tt.expected)
			}
		})
	}
}

func TestTaskStatusTransitions(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		fromStatus  TaskStatus
		toStatus    TaskStatus
		description string
	}{
		{"Pending to Processing", TaskStatusPending, TaskStatusProcessing, "Task starts processing"},
		{"Processing to Completed", TaskStatusProcessing, TaskStatusCompleted, "Task completes successfully"},
		{"Processing to Failed", TaskStatusProcessing, TaskStatusFailed, "Task fails during processing"},
		{"Pending to Failed", TaskStatusPending, TaskStatusFailed, "Task fails before starting"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				ID:     string(tt.fromStatus),
				Status: tt.fromStatus,
			}

			// Transition to new status
			task.Status = tt.toStatus
			task.UpdatedAt = now

			if task.Status != tt.toStatus {
				t.Errorf("Expected status %v, got %v", tt.toStatus, task.Status)
			}
		})
	}
}

func TestBatchTaskStatusTransitions(t *testing.T) {
	tests := []struct {
		name           string
		pending        int
		success        int
		failed         int
		expectedStatus TaskStatus
	}{
		{
			name:           "All completed",
			pending:        0,
			success:        100,
			failed:         0,
			expectedStatus: TaskStatusCompleted,
		},
		{
			name:           "All failed",
			pending:        0,
			success:        0,
			failed:         100,
			expectedStatus: TaskStatusFailed,
		},
		{
			name:           "Mixed with success majority",
			pending:        0,
			success:        60,
			failed:         40,
			expectedStatus: TaskStatusCompleted,
		},
		{
			name:           "Mixed with no success",
			pending:        0,
			success:        0,
			failed:         100,
			expectedStatus: TaskStatusFailed,
		},
		{
			name:           "Still pending",
			pending:        50,
			success:        25,
			failed:         25,
			expectedStatus: TaskStatusProcessing,
		},
		{
			name:           "All pending",
			pending:        100,
			success:        0,
			failed:         0,
			expectedStatus: TaskStatusProcessing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchTask := BatchTask{
				ID:      "batch-123",
				Type:    TaskTypeEmbedding,
				Pending: tt.pending,
				Success: tt.success,
				Failed:  tt.failed,
			}

			// Simulate UpdateBatchTask logic
			if batchTask.Pending == 0 {
				if batchTask.Failed == 0 {
					batchTask.Status = TaskStatusCompleted
				} else if batchTask.Success == 0 {
					batchTask.Status = TaskStatusFailed
				} else {
					batchTask.Status = TaskStatusCompleted
				}
			} else {
				batchTask.Status = TaskStatusProcessing
			}

			if batchTask.Status != tt.expectedStatus {
				t.Errorf("Expected status %v, got %v", tt.expectedStatus, batchTask.Status)
			}
		})
	}
}

func TestTaskWithOptionalFields(t *testing.T) {
	// Test with optional DocID field
	taskWithDocID := Task{
		ID:     "task-123",
		Type:   TaskTypeDocumentEmbedding,
		DocID:  "doc-456",
		Status: TaskStatusPending,
	}

	if taskWithDocID.DocID != "doc-456" {
		t.Errorf("DocID = %s, want doc-456", taskWithDocID.DocID)
	}

	// Test without DocID
	taskWithoutDocID := Task{
		ID:     "task-124",
		Type:   TaskTypeEmbedding,
		Status: TaskStatusPending,
	}

	if taskWithoutDocID.DocID != "" {
		t.Errorf("Expected empty DocID, got %s", taskWithoutDocID.DocID)
	}

	// Test with optional Error field
	taskWithError := Task{
		ID:     "task-125",
		Type:   TaskTypeEmbedding,
		Status: TaskStatusFailed,
		Error:  "API timeout",
	}

	if taskWithError.Error != "API timeout" {
		t.Errorf("Error = %s, want 'API timeout'", taskWithError.Error)
	}

	// Test without Error field
	taskWithoutError := Task{
		ID:     "task-126",
		Type:   TaskTypeEmbedding,
		Status: TaskStatusPending,
	}

	if taskWithoutError.Error != "" {
		t.Errorf("Expected empty Error, got %s", taskWithoutError.Error)
	}
}

func TestTaskTimestamps(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	task := Task{
		ID:        "task-123",
		Status:    TaskStatusCompleted,
		CreatedAt: oneHourAgo,
		UpdatedAt: now,
	}

	if task.CreatedAt.Equal(oneHourAgo) {
		// Timestamps should be set
	} else {
		t.Error("CreatedAt should be set to provided value")
	}

	if task.UpdatedAt.Equal(now) {
		// Timestamps should be set
	} else {
		t.Error("UpdatedAt should be set to provided value")
	}
}

func TestBatchTaskWithAllStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
	}{
		{"Pending batch", TaskStatusPending},
		{"Processing batch", TaskStatusProcessing},
		{"Completed batch", TaskStatusCompleted},
		{"Failed batch", TaskStatusFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch := BatchTask{
				ID:     string(tt.status),
				Status: tt.status,
				Total:  100,
			}

			if batch.Status != tt.status {
				t.Errorf("Expected status %v, got %v", tt.status, batch.Status)
			}
		})
	}
}

func TestTaskForceFlag(t *testing.T) {
	tests := []struct {
		name  string
		force bool
	}{
		{"Force enabled", true},
		{"Force disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				ID:     "task-123",
				Force:  tt.force,
				Status: TaskStatusPending,
			}

			if task.Force != tt.force {
				t.Errorf("Expected Force to be %v, got %v", tt.force, task.Force)
			}
		})
	}
}
