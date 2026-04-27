package models

import "time"

// EducationBook 黄梅教育板块的 PDF 书籍资源。
type EducationBook struct {
	BookID      uint      `gorm:"column:book_id;primaryKey" json:"book_id"`
	Title       string    `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	PdfPath     string    `gorm:"column:pdf_path;type:varchar(255);not null" json:"pdf_path"`
	CoverPath   string    `gorm:"column:cover_path;type:varchar(255)" json:"cover_path"`
	IsPublished bool      `gorm:"column:is_published;default:false;not null" json:"is_published"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (EducationBook) TableName() string {
	return "education_books"
}

type EducationBookListResponse struct {
	BookID      uint      `json:"book_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PdfPath     string    `json:"pdf_path"`
	PdfURL      string    `json:"pdf_url"`
	CoverPath   string    `json:"cover_path"`
	CoverURL    string    `json:"cover_url,omitempty"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EducationBookDetailResponse struct {
	BookID      uint      `json:"book_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PdfPath     string    `json:"pdf_path"`
	PdfURL      string    `json:"pdf_url"`
	CoverPath   string    `json:"cover_path"`
	CoverURL    string    `json:"cover_url,omitempty"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateEducationBookRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	PdfPath     string `json:"pdf_path" binding:"required"`
	CoverPath   string `json:"cover_path"`
	IsPublished bool   `json:"is_published"`
}

type UpdateEducationBookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PdfPath     string `json:"pdf_path"`
	CoverPath   string `json:"cover_path"`
	IsPublished *bool  `json:"is_published"`
}
