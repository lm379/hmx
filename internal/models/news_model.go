package models

import "time"

// News 黄梅戏新闻资讯表
type News struct {
	NewsID      uint       `gorm:"column:news_id;primaryKey" json:"news_id"`
	Title       string     `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Summary     string     `gorm:"column:summary;type:varchar(500)" json:"summary"`
	Content     string     `gorm:"column:content;type:text;not null" json:"content"`
	Cover       string     `gorm:"column:cover;type:varchar(255)" json:"cover"`
	Source      string     `gorm:"column:source;type:varchar(100)" json:"source"`
	Author      string     `gorm:"column:author;type:varchar(100)" json:"author"`
	IsPublished bool       `gorm:"column:is_published;default:false;not null" json:"is_published"`
	PublishedAt *time.Time `gorm:"column:published_at" json:"published_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (News) TableName() string {
	return "news"
}

type NewsListResponse struct {
	NewsID      uint       `json:"news_id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Cover       string     `json:"cover,omitempty"`
	Source      string     `json:"source,omitempty"`
	Author      string     `json:"author,omitempty"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type NewsDetailResponse struct {
	NewsID      uint       `json:"news_id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content"`
	Cover       string     `json:"cover,omitempty"`
	Source      string     `json:"source,omitempty"`
	Author      string     `json:"author,omitempty"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateNewsRequest struct {
	Title       string     `json:"title" binding:"required"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content" binding:"required"`
	Cover       string     `json:"cover"`
	Source      string     `json:"source"`
	Author      string     `json:"author"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at"`
}

type UpdateNewsRequest struct {
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content"`
	Cover       string     `json:"cover"`
	Source      string     `json:"source"`
	Author      string     `json:"author"`
	IsPublished *bool      `json:"is_published"`
	PublishedAt *time.Time `json:"published_at"`
}
