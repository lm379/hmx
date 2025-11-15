package models

import (
	"database/sql"
	"time"
)

// Artist (艺术家表)
type Artist struct {
	ArtistID  uint      `gorm:"column:artist_id;primaryKey"`
	Name      string    `gorm:"column:name;type:varchar(100);unique;not null"`
	Bio       string    `gorm:"column:bio;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	Operas    []*Opera  `gorm:"many2many:opera_artists;"`
}

// TableName 指定表名
func (Artist) TableName() string {
	return "artists"
}

// Opera (黄梅戏作品表)
type Opera struct {
	OperaID     uint           `gorm:"column:opera_id;primaryKey"`
	OperaTitle  string         `gorm:"column:opera_title;type:varchar(100);not null"`
	ReleaseDate sql.NullTime   `gorm:"column:release_date"`
	Duration    sql.NullString `gorm:"column:duration;type:time"` // 时长
	MusicPath   sql.NullString `gorm:"column:music_path;type:varchar(255)"`
	VideoPath   string         `gorm:"column:video_path;type:varchar(255);not null"`
	Description string         `gorm:"column:description;type:text"`
	Avatar      sql.NullString `gorm:"column:avatar;type:varchar(255)"`
	AiSummary   string         `gorm:"column:ai_summary;type:text"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
	Artists     []*Artist      `gorm:"many2many:opera_artists;"`
}

// TableName 指定表名
func (Opera) TableName() string {
	return "opera"
}

// OperaResponse 用于 API 返回的 DTO，自动处理 URL 拼接
type OperaResponse struct {
	OperaID     uint      `json:"opera_id"`
	OperaTitle  string    `json:"opera_title"`
	ReleaseDate *string   `json:"release_date"`
	Duration    *string   `json:"duration"`
	MusicPath   *string   `json:"music_path"` // 完整 URL
	VideoPath   string    `json:"video_path"` // 完整 URL
	Description string    `json:"description"`
	Avatar      *string   `json:"avatar"` // 完整 URL
	AiSummary   string    `json:"ai_summary"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PresignRequest DTO (POST /uploads/presign)
type PresignRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	UploadType  string `json:"upload_type" binding:"required"` // 'videos', 'avatars'
}

// PresignResponse DTO
type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

// CreateOperaRequest DTO (POST /operas)
type CreateOperaRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	VideoPath   string `json:"video_path" binding:"required"` // 必须是 PresignResponse.ObjectKey
	AvatarPath  string `json:"avatar_path"`                   // 封面的 ObjectKey
	ArtistIDs   []uint `json:"artist_ids"`                    // 关联的艺术家ID
}

// UpdateStatusRequest DTO (PATCH /admin/operas/:id/status)
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}
