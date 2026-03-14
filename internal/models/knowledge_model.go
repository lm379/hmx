package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// SourceType 知识库文档来源类型
type SourceType string

const (
	SourceTypeOpera        SourceType = "opera"        // 来自黄梅戏作品字幕
	SourceTypeProfessional SourceType = "professional" // 来自管理员上传的专业知识文档
)

// EmbeddingStatus 向量化状态
type EmbeddingStatus string

const (
	EmbeddingStatusPending    EmbeddingStatus = "pending"
	EmbeddingStatusProcessing EmbeddingStatus = "processing"
	EmbeddingStatusCompleted  EmbeddingStatus = "completed"
	EmbeddingStatusFailed     EmbeddingStatus = "failed"
)

// KnowledgeDocument 知识库文档表
type KnowledgeDocument struct {
	DocID           uuid.UUID       `gorm:"column:doc_id;type:uuid;primaryKey;default:gen_random_uuid()" json:"doc_id"`
	Title           string          `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Content         string          `gorm:"column:content;type:text;not null" json:"content"`
	SourceType      SourceType      `gorm:"column:source_type;type:varchar(50);not null" json:"source_type"`
	ChunksCount     int             `gorm:"column:chunks_count;not null;default:0" json:"chunks_count"`
	EmbeddingStatus EmbeddingStatus `gorm:"column:embedding_status;type:varchar(50);not null;default:'pending'" json:"embedding_status"`
	CreatedBy       uint            `gorm:"column:created_by;type:integer;not null" json:"created_by"`
	CreatedAt       time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsActive        bool            `gorm:"column:is_active;not null;default:true" json:"is_active"`
}

// TableName 指定表名
func (KnowledgeDocument) TableName() string {
	return "knowledge_documents"
}

// KnowledgeEmbedding 文档向量块表（分块存储）
type KnowledgeEmbedding struct {
	ChunkID    uuid.UUID       `gorm:"column:chunk_id;type:uuid;primaryKey;default:gen_random_uuid()"`
	DocID      uuid.UUID       `gorm:"column:doc_id;type:uuid;not null"`
	ChunkIndex int             `gorm:"column:chunk_index;not null"`
	ChunkText  string          `gorm:"column:chunk_text;type:text;not null"`
	Embedding  pgvector.Vector `gorm:"column:embedding;type:vector(4096);not null"`
	CreatedAt  time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (KnowledgeEmbedding) TableName() string {
	return "knowledge_embeddings"
}

// QAHistory 问答历史表
type QAHistory struct {
	QAID            uuid.UUID   `gorm:"column:qa_id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID          *uuid.UUID  `gorm:"column:user_id;type:uuid"`            // nullable：游客为 nil
	SessionID       string      `gorm:"column:session_id;type:varchar(255)"` // 游客 session
	Question        string      `gorm:"column:question;type:text;not null"`
	Answer          string      `gorm:"column:answer;type:text;not null"`
	RelatedDocIDs   UUIDArray   `gorm:"column:related_doc_ids;type:uuid[];not null;default:'{}'"`
	RelatedOperaIDs Uint32Array `gorm:"column:related_opera_ids;type:integer[];not null;default:'{}'"`
	CreatedAt       time.Time   `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time   `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (QAHistory) TableName() string {
	return "qa_history"
}

// QAFeedback 问答反馈表
type QAFeedback struct {
	FeedbackID uuid.UUID  `gorm:"column:feedback_id;type:uuid;primaryKey;default:gen_random_uuid()"`
	QAID       uuid.UUID  `gorm:"column:qa_id;type:uuid;not null"`
	UserID     *uuid.UUID `gorm:"column:user_id;type:uuid"`
	Rating     int        `gorm:"column:rating;check:rating >= 1 AND rating <= 5"`
	Comments   string     `gorm:"column:comments;type:text"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (QAFeedback) TableName() string {
	return "qa_feedback"
}

// =============================================
// 自定义 PostgreSQL 数组类型（GORM Scanner/Valuer）
// =============================================

// UUIDArray UUID 数组，映射 PostgreSQL uuid[]
type UUIDArray []uuid.UUID

func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	parts := make([]string, len(a))
	for i, u := range a {
		parts[i] = u.String()
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *UUIDArray) Scan(src interface{}) error {
	if src == nil {
		*a = UUIDArray{}
		return nil
	}
	s, ok := src.(string)
	if !ok {
		b, ok := src.([]byte)
		if !ok {
			return fmt.Errorf("UUIDArray.Scan: unsupported type %T", src)
		}
		s = string(b)
	}
	s = strings.Trim(s, "{}")
	if s == "" {
		*a = UUIDArray{}
		return nil
	}
	parts := strings.Split(s, ",")
	arr := make(UUIDArray, len(parts))
	for i, p := range parts {
		u, err := uuid.Parse(strings.TrimSpace(p))
		if err != nil {
			return fmt.Errorf("UUIDArray.Scan: invalid uuid %q: %w", p, err)
		}
		arr[i] = u
	}
	*a = arr
	return nil
}

// Uint32Array uint 数组，映射 PostgreSQL integer[]
type Uint32Array []uint

func (a Uint32Array) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *Uint32Array) Scan(src interface{}) error {
	if src == nil {
		*a = Uint32Array{}
		return nil
	}
	s, ok := src.(string)
	if !ok {
		b, ok := src.([]byte)
		if !ok {
			return fmt.Errorf("Uint32Array.Scan: unsupported type %T", src)
		}
		s = string(b)
	}
	s = strings.Trim(s, "{}")
	if s == "" {
		*a = Uint32Array{}
		return nil
	}
	parts := strings.Split(s, ",")
	arr := make(Uint32Array, len(parts))
	for i, p := range parts {
		var v uint
		_, err := fmt.Sscanf(strings.TrimSpace(p), "%d", &v)
		if err != nil {
			return fmt.Errorf("Uint32Array.Scan: invalid uint %q: %w", p, err)
		}
		arr[i] = v
	}
	*a = arr
	return nil
}

// =============================================
// DTO
// =============================================

// UploadDocumentRequest 管理员上传知识库文档请求（JSON body）
type UploadDocumentRequest struct {
	Title      string     `json:"title" binding:"required"`
	SourceType SourceType `json:"source_type" binding:"required,oneof=opera professional"`
	Content    string     `json:"content" binding:"required,min=10"`
}

// DocumentListItem 文档列表条目响应
type DocumentListItem struct {
	DocID           uuid.UUID       `json:"doc_id"`
	Title           string          `json:"title"`
	SourceType      SourceType      `json:"source_type"`
	ChunksCount     int             `json:"chunks_count"`
	EmbeddingStatus EmbeddingStatus `json:"embedding_status"`
	IsActive        bool            `json:"is_active"`
	CreatedAt       time.Time       `json:"created_at"`
}

// ChunkItem 单个分块信息（用于详情接口）
type ChunkItem struct {
	ChunkID    uuid.UUID `json:"chunk_id"`
	ChunkIndex int       `json:"chunk_index"`
	ChunkText  string    `json:"chunk_text"`
}

// DocumentDetailResponse 文档详情响应（含分块列表）
type DocumentDetailResponse struct {
	DocID           uuid.UUID       `json:"doc_id"`
	Title           string          `json:"title"`
	Content         string          `json:"content"`
	SourceType      SourceType      `json:"source_type"`
	ChunksCount     int             `json:"chunks_count"`
	EmbeddingStatus EmbeddingStatus `json:"embedding_status"`
	IsActive        bool            `json:"is_active"`
	CreatedBy       uint            `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Chunks          []ChunkItem     `json:"chunks"`
}

// AskRequest 问答请求 DTO
type AskRequest struct {
	Query     string `json:"query" binding:"required"`
	SessionID string `json:"session_id"` // 游客 session ID
	OperaID   *uint  `json:"opera_id"`   // 可选：当前作品上下文
}

// ChunkSource 检索到的文档块来源信息
type ChunkSource struct {
	ChunkID   uuid.UUID `json:"chunk_id"`
	DocID     uuid.UUID `json:"doc_id"`
	DocTitle  string    `json:"doc_title"`
	ChunkText string    `json:"chunk_text"`
	Score     float64   `json:"score"`
}

// AskResponse 问答响应 DTO
type AskResponse struct {
	QAID          uuid.UUID     `json:"qa_id"`
	Answer        string        `json:"answer"`
	Sources       []ChunkSource `json:"sources"`
	RelatedOperas []uint        `json:"related_operas"`
}

// QAHistoryItem 问答历史条目响应
type QAHistoryItem struct {
	QAID            uuid.UUID `json:"qa_id"`
	Question        string    `json:"question"`
	Answer          string    `json:"answer"`
	RelatedDocIDs   []string  `json:"related_doc_ids"`
	RelatedOperaIDs []uint    `json:"related_opera_ids"`
	CreatedAt       time.Time `json:"created_at"`
}

// FeedbackRequest 反馈提交请求 DTO
type FeedbackRequest struct {
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Comments string `json:"comments"`
}
