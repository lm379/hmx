package models

import (
	"database/sql"
	"time"

	"github.com/pgvector/pgvector-go"
)

// Artist (艺术家表)
type Artist struct {
	ArtistID     uint           `gorm:"column:artist_id;primaryKey"`
	Name         string         `gorm:"column:name;type:varchar(100);unique;not null"`
	Bio          string         `gorm:"column:bio;type:text"`
	ArtistAvatar sql.NullString `gorm:"column:artist_avatar;type:varchar(255)"`
	CreatedAt    time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	Operas       []*Opera       `gorm:"many2many:opera_artists;joinForeignKey:artist_id;joinReferences:opera_id"`
}

// TableName 指定表名
func (Artist) TableName() string {
	return "artists"
}

// Opera (黄梅戏作品表)
type Opera struct {
	OperaID     uint            `gorm:"column:opera_id;primaryKey"`
	OperaTitle  string          `gorm:"column:opera_title;type:varchar(100);not null"`
	ReleaseDate sql.NullTime    `gorm:"column:release_date"`
	Duration    sql.NullString  `gorm:"column:duration;type:time"` // 时长
	MusicPath   sql.NullString  `gorm:"column:music_path;type:varchar(255)"`
	VideoPath   string          `gorm:"column:video_path;type:varchar(255);not null"`
	SrtPath     sql.NullString  `gorm:"column:srt_path;type:varchar(255)"`
	Description string          `gorm:"column:description;type:text"`
	Avatar      sql.NullString  `gorm:"column:avatar;type:varchar(255)"`
	AiSummary   string          `gorm:"column:ai_summary;type:text"`
	Embedding   pgvector.Vector `gorm:"column:embedding;type:vector(1536)"` // AI向量表示，使用pgvector
	IsHidden    bool            `gorm:"column:is_hidden;default:false;not null"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
	Artists     []*Artist       `gorm:"many2many:opera_artists;joinForeignKey:opera_id;joinReferences:artist_id"`
}

// TableName 指定表名
func (Opera) TableName() string {
	return "opera"
}

// SimpleArtist 简化的艺术家信息用于 Opera 响应
type SimpleArtist struct {
	ArtistID uint   `json:"artist_id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
}

// SimpleOpera 简化的作品信息用于 Artist 响应（不包含艺术家信息，避免冗余）
type SimpleOpera struct {
	OperaID    uint      `json:"opera_id"`
	OperaTitle string    `json:"opera_title"`
	Avatar     *string   `json:"avatar,omitempty"`
	Duration   *string   `json:"duration,omitempty"`
	PlayCount  int64     `json:"play_count"`
	LikeCount  int64     `json:"like_count"`
	CreatedAt  time.Time `json:"created_at"`
}

// OperaListResponse 用于列表页的简化响应
type OperaListResponse struct {
	OperaID    uint           `json:"opera_id"`
	OperaTitle string         `json:"opera_title"`
	Artists    []SimpleArtist `json:"artists"`
	Avatar     *string        `json:"avatar"`
	Duration   *string        `json:"duration,omitempty"`
	IsHidden   bool           `json:"is_hidden"`
	PlayCount  int64          `json:"play_count"`
	LikeCount  int64          `json:"like_count"`
	CreatedAt  time.Time      `json:"created_at"`
}

// OperaDetailResponse 用于详情页的完整响应
type OperaDetailResponse struct {
	OperaID       uint           `json:"opera_id"`
	OperaTitle    string         `json:"opera_title"`
	Artists       []SimpleArtist `json:"artists"`
	ReleaseDate   *string        `json:"release_date,omitempty"`
	Duration      *string        `json:"duration,omitempty"`
	MusicPath     *string        `json:"music_path,omitempty"`
	VideoPath     string         `json:"video_path"`
	SrtPath       *string        `json:"srt_path,omitempty"`
	Description   string         `json:"description"`
	Avatar        *string        `json:"avatar"`
	AiSummary     string         `json:"ai_summary,omitempty"`
	IsHidden      bool           `json:"is_hidden"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	LikeCount     int64          `json:"like_count"`
	FavoriteCount int64          `json:"favorite_count"`
	ShareCount    int64          `json:"share_count"`
	PlayCount     int64          `json:"play_count"`
	Liked         bool           `json:"liked"`
	Favorited     bool           `json:"favorited"`
}

// OperaResponse 通用响应（向后兼容）
type OperaResponse struct {
	OperaID       uint           `json:"opera_id"`
	OperaTitle    string         `json:"opera_title"`
	Artists       []SimpleArtist `json:"artists"`
	ReleaseDate   *string        `json:"release_date"`
	Duration      *string        `json:"duration"`
	MusicPath     *string        `json:"music_path"` // 完整 URL
	VideoPath     string         `json:"video_path"` // 完整 URL
	SrtPath       *string        `json:"srt_path"`   // 字幕 URL
	Description   string         `json:"description"`
	Avatar        *string        `json:"avatar"` // 完整 URL
	AiSummary     string         `json:"ai_summary"`
	IsHidden      bool           `json:"is_hidden"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	LikeCount     int64          `json:"like_count"`
	FavoriteCount int64          `json:"favorite_count"`
	ShareCount    int64          `json:"share_count"`
	PlayCount     int64          `json:"play_count"`
	Liked         bool           `json:"liked"`
	Favorited     bool           `json:"favorited"`
}

// ArtistListResponse 用于艺术家列表的简化响应
type ArtistListResponse struct {
	ArtistID  uint      `json:"artist_id"`
	Name      string    `json:"name"`
	Avatar    *string   `json:"avatar"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ArtistDetailResponse 用于艺术家详情的完整响应
type ArtistDetailResponse struct {
	ArtistID  uint          `json:"artist_id"`
	Name      string        `json:"name"`
	Bio       string        `json:"bio"`
	Avatar    *string       `json:"avatar"`
	Operas    []SimpleOpera `json:"operas,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// ArtistResponse 通用响应（向后兼容）
type ArtistResponse struct {
	ArtistID  uint             `json:"artist_id"`
	Name      string           `json:"name"`
	Bio       string           `json:"bio"`
	Avatar    *string          `json:"avatar"`           // 艺术家头像 URL
	Operas    []*OperaResponse `json:"operas,omitempty"` // 包含的作品列表
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// PresignRequest DTO (POST /uploads/presign)
type PresignRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	UploadType  string `json:"upload_type" binding:"required"` // 'user_avatar', 'artist_avatar', 'opera_cover', 'video_upload'
	TargetID    uint   `json:"target_id"`                      // 目标ID：用户ID、艺术家ID或曲目ID (video_upload时可为0)
}

// PresignResponse DTO
type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

// CreateOperaRequest DTO (POST /operas)
type CreateOperaRequest struct {
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description"`
	VideoPath      string   `json:"video_path" binding:"required"` // 必须是 PresignResponse.ObjectKey
	AvatarPath     string   `json:"avatar_path"`                   // 封面的 ObjectKey
	Duration       string   `json:"duration"`                      // 时长 (HH:MM:SS)
	ArtistIDs      []uint   `json:"artist_ids"`                    // 关联的艺术家ID
	NewArtistNames []string `json:"new_artist_names"`              // 新增艺术家名字（自动创建）
}

// UpdateOperaRequest DTO (PUT /admin/operas/:id)
type UpdateOperaRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	VideoPath      string   `json:"video_path"`
	AvatarPath     string   `json:"avatar_path"`
	Duration       string   `json:"duration"`
	ArtistIDs      []uint   `json:"artist_ids"`
	NewArtistNames []string `json:"new_artist_names"`
	IsHidden       *bool    `json:"is_hidden"`
}

// CreateArtistRequest DTO (POST /admin/artists)
type CreateArtistRequest struct {
	Name   string `json:"name" binding:"required"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
}

// UpdateArtistRequest DTO (PUT /admin/artists/:id)
type UpdateArtistRequest struct {
	Name   string `json:"name"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
}

// UpdateUserRoleRequest DTO (PUT /admin/users/:id/role)
type UpdateUserRoleRequest struct {
	Role UserRole `json:"role" binding:"required,oneof=Administrator User"`
}

// UpdateStatusRequest DTO (PATCH /admin/operas/:id/status)
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

// TencentCloudCallbackRequest 腾讯云数据万象任务回调请求
type TencentCloudCallbackRequest struct {
	EventName  string      `json:"EventName"`  // 事件名称：TaskFinish
	JobsDetail []JobDetail `json:"JobsDetail"` // 任务详情列表
}

// JobDetail 任务详情
type JobDetail struct {
	Code         string       `json:"Code"`         // 状态码：Success
	Message      string       `json:"Message"`      // 状态消息
	JobId        string       `json:"JobId"`        // 任务ID
	Tag          string       `json:"Tag"`          // 任务类型：Transcode, SmartCover, SpeechRecognition
	State        string       `json:"State"`        // 任务状态：Success, Failed
	CreationTime string       `json:"CreationTime"` // 创建时间
	StartTime    string       `json:"StartTime"`    // 开始时间
	EndTime      string       `json:"EndTime"`      // 结束时间
	QueueId      string       `json:"QueueId"`      // 队列ID
	Input        JobInput     `json:"Input"`        // 输入信息
	Operation    JobOperation `json:"Operation"`    // 操作信息
	Workflow     WorkflowInfo `json:"Workflow"`     // 工作流信息
}

// JobInput 输入信息
type JobInput struct {
	BucketId string `json:"BucketId"` // Bucket名称
	Object   string `json:"Object"`   // 输入文件路径，如：tmp/黄梅戏.mp4
	Region   string `json:"Region"`   // 地域
}

// JobOperation 操作信息
type JobOperation struct {
	TemplateId              string        `json:"TemplateId"`              // 模板ID
	TemplateName            string        `json:"TemplateName"`            // 模板名称
	Output                  OutputConfig  `json:"Output"`                  // 输出配置
	MediaResult             *MediaResult  `json:"MediaResult"`             // 媒体处理结果
	MediaInfo               *MediaInfo    `json:"MediaInfo"`               // 媒体信息（转码任务）
	SpeechRecognitionResult *SpeechResult `json:"SpeechRecognitionResult"` // 语音识别结果
}

// OutputConfig 输出配置
type OutputConfig struct {
	Region string `json:"Region"` // 输出地域
	Bucket string `json:"Bucket"` // 输出Bucket
	Object string `json:"Object"` // 输出对象路径模板
}

// MediaResult 媒体处理结果
type MediaResult struct {
	OutputFile OutputFile `json:"OutputFile"` // 输出文件信息
}

// OutputFile 输出文件信息
type OutputFile struct {
	Bucket     string   `json:"Bucket"`     // Bucket名称
	Region     string   `json:"Region"`     // 地域
	ObjectName []string `json:"ObjectName"` // 输出文件路径列表
	ObjectUrl  []string `json:"ObjectUrl"`  // 输出文件URL列表
}

// MediaInfo 媒体信息
type MediaInfo struct {
	Format MediaFormat `json:"Format"` // 格式信息
}

// MediaFormat 格式信息
type MediaFormat struct {
	Duration string `json:"Duration"` // 时长（秒），如："6774.499999"
}

// SpeechResult 语音识别结果
type SpeechResult struct {
	ObjectName string  `json:"ObjectName"` // 字幕文件路径
	AudioTime  float64 `json:"AudioTime"`  // 音频时长（秒）
}

// WorkflowInfo 工作流信息
type WorkflowInfo struct {
	WorkflowId   string `json:"WorkflowId"`   // 工作流ID
	WorkflowName string `json:"WorkflowName"` // 工作流名称
	RunId        string `json:"RunId"`        // 执行ID
}
