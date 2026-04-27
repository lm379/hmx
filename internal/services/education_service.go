package services

import (
	"strings"

	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/pagination"
	"github.com/lm379/hmx/pkg/urlutils"
)

type EducationService struct {
	educationRepo *repository.EducationRepo
}

func NewEducationService() *EducationService {
	return &EducationService{educationRepo: repository.NewEducationRepo()}
}

// GetPublishedBooks 获取前台教育资源列表。
func (s *EducationService) GetPublishedBooks(p *pagination.Pagination) ([]models.EducationBookListResponse, int64, error) {
	books, total, err := s.educationRepo.GetPublished(p)
	if err != nil {
		return nil, 0, err
	}
	return toEducationBookList(books), total, nil
}

// GetAllBooks 获取后台教育资源列表。
func (s *EducationService) GetAllBooks(p *pagination.Pagination) ([]models.EducationBookListResponse, int64, error) {
	books, total, err := s.educationRepo.GetAll(p)
	if err != nil {
		return nil, 0, err
	}
	return toEducationBookList(books), total, nil
}

// GetPublishedBookDetail 获取前台 PDF 阅读详情。
func (s *EducationService) GetPublishedBookDetail(id uint) (*models.EducationBookDetailResponse, error) {
	book, err := s.educationRepo.GetPublishedByID(id)
	if err != nil {
		return nil, err
	}
	return toEducationBookDetail(*book), nil
}

// GetAdminBookDetail 获取后台编辑详情。
func (s *EducationService) GetAdminBookDetail(id uint) (*models.EducationBookDetailResponse, error) {
	book, err := s.educationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toEducationBookDetail(*book), nil
}

// CreateBook 创建教育 PDF 资源。
func (s *EducationService) CreateBook(input models.CreateEducationBookRequest) (*models.EducationBookDetailResponse, error) {
	book := &models.EducationBook{
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		PdfPath:     strings.TrimSpace(input.PdfPath),
		CoverPath:   strings.TrimSpace(input.CoverPath),
		IsPublished: input.IsPublished,
	}
	if err := s.educationRepo.Create(book); err != nil {
		return nil, err
	}
	return toEducationBookDetail(*book), nil
}

// UpdateBook 按后台完整表单更新教育 PDF 资源。
func (s *EducationService) UpdateBook(id uint, input models.UpdateEducationBookRequest) (*models.EducationBookDetailResponse, error) {
	book, err := s.educationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"title":       strings.TrimSpace(input.Title),
		"description": strings.TrimSpace(input.Description),
		"pdf_path":    strings.TrimSpace(input.PdfPath),
		"cover_path":  strings.TrimSpace(input.CoverPath),
	}
	if input.IsPublished != nil {
		updates["is_published"] = *input.IsPublished
	}
	if err := s.educationRepo.Update(book, updates); err != nil {
		return nil, err
	}
	updated, err := s.educationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toEducationBookDetail(*updated), nil
}

// DeleteBook 删除教育 PDF 资源记录。
func (s *EducationService) DeleteBook(id uint) error {
	return s.educationRepo.Delete(id)
}

func toEducationBookList(books []models.EducationBook) []models.EducationBookListResponse {
	result := make([]models.EducationBookListResponse, 0, len(books))
	for _, book := range books {
		result = append(result, models.EducationBookListResponse{
			BookID:      book.BookID,
			Title:       book.Title,
			Description: book.Description,
			PdfPath:     book.PdfPath,
			PdfURL:      urlutils.GetFullURL(book.PdfPath),
			CoverPath:   book.CoverPath,
			CoverURL:    urlutils.GetFullURL(book.CoverPath),
			IsPublished: book.IsPublished,
			CreatedAt:   book.CreatedAt,
			UpdatedAt:   book.UpdatedAt,
		})
	}
	return result
}

func toEducationBookDetail(book models.EducationBook) *models.EducationBookDetailResponse {
	return &models.EducationBookDetailResponse{
		BookID:      book.BookID,
		Title:       book.Title,
		Description: book.Description,
		PdfPath:     book.PdfPath,
		PdfURL:      urlutils.GetFullURL(book.PdfPath),
		CoverPath:   book.CoverPath,
		CoverURL:    urlutils.GetFullURL(book.CoverPath),
		IsPublished: book.IsPublished,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}
