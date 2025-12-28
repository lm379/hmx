package models

import (
	"database/sql"
	"time"
)

// --- Like (点赞) ---
type Like struct {
	UserID    uint      `gorm:"column:user_id;primaryKey"`  // 复合主键
	OperaID   uint      `gorm:"column:opera_id;primaryKey"` // 复合主键
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}

// --- Favorite (收藏) ---
type Favorite struct {
	UserID    uint      `gorm:"column:user_id;primaryKey"`  // 复合主键
	OperaID   uint      `gorm:"column:opera_id;primaryKey"` // 复合主键
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

// --- Comment (评论) ---
type Comment struct {
	CommentID       uint          `gorm:"column:comment_id;primaryKey"`
	UserID          sql.NullInt64 `gorm:"column:user_id"` // 可为空，用户删除时设置为 NULL
	OperaID         uint          `gorm:"column:opera_id;not null"`
	ParentCommentID *uint         `gorm:"column:parent_comment_id"` // 用于回复 (NULLable)
	CommentText     string        `gorm:"column:comment_text;type:text;not null"`
	TsvComment      string        `gorm:"column:tsv_comment;type:tsvector"` // FTS (只读)
	CreatedAt       time.Time     `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time     `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// --- PlayHistory (播放历史) ---
type PlayHistory struct {
	PlayID    uint      `gorm:"column:play_id;primaryKey"`
	UserID    *uint     `gorm:"column:user_id"` // 可为空，支持游客
	OperaID   uint      `gorm:"column:opera_id"`
	Count     int       `gorm:"column:count;default:1"` // 播放次数
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (PlayHistory) TableName() string {
	return "play_history"
}

// CreateCommentRequest DTO (POST /operas/:id/comments)
type CreateCommentRequest struct {
	CommentText     string `json:"comment_text" binding:"required,max=1000"`
	ParentCommentID *uint  `json:"parent_comment_id"` // 可选，用于回复
}
