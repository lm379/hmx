package services

import (
	"strings"

	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/pagination"
	"github.com/lm379/hmx/pkg/urlutils"
)

type NewsService struct {
	newsRepo *repository.NewsRepo
}

func NewNewsService() *NewsService {
	return &NewsService{newsRepo: repository.NewNewsRepo()}
}

// GetPublishedNews 获取前台新闻列表。
func (s *NewsService) GetPublishedNews(p *pagination.Pagination) ([]models.NewsListResponse, int64, error) {
	news, total, err := s.newsRepo.GetPublished(p)
	if err != nil {
		return nil, 0, err
	}
	return toNewsList(news), total, nil
}

// GetAllNews 获取后台新闻列表。
func (s *NewsService) GetAllNews(p *pagination.Pagination) ([]models.NewsListResponse, int64, error) {
	news, total, err := s.newsRepo.GetAll(p)
	if err != nil {
		return nil, 0, err
	}
	return toNewsList(news), total, nil
}

// GetPublishedNewsDetail 获取前台新闻详情。
func (s *NewsService) GetPublishedNewsDetail(id uint) (*models.NewsDetailResponse, error) {
	news, err := s.newsRepo.GetPublishedByID(id)
	if err != nil {
		return nil, err
	}
	return toNewsDetail(*news), nil
}

// GetAdminNewsDetail 获取后台编辑详情。
func (s *NewsService) GetAdminNewsDetail(id uint) (*models.NewsDetailResponse, error) {
	news, err := s.newsRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toNewsDetail(*news), nil
}

// CreateNews 创建新闻内容。
func (s *NewsService) CreateNews(input models.CreateNewsRequest) (*models.NewsDetailResponse, error) {
	news := &models.News{
		Title:       strings.TrimSpace(input.Title),
		Summary:     strings.TrimSpace(input.Summary),
		Content:     strings.TrimSpace(input.Content),
		Cover:       strings.TrimSpace(input.Cover),
		Source:      strings.TrimSpace(input.Source),
		Author:      strings.TrimSpace(input.Author),
		IsPublished: input.IsPublished,
		PublishedAt: input.PublishedAt,
	}
	if err := s.newsRepo.Create(news); err != nil {
		return nil, err
	}
	return toNewsDetail(*news), nil
}

// UpdateNews 按后台完整表单更新新闻。
func (s *NewsService) UpdateNews(id uint, input models.UpdateNewsRequest) (*models.NewsDetailResponse, error) {
	news, err := s.newsRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 管理端提交的是完整表单，允许清空摘要、封面、来源等可选字段。
	updates := map[string]interface{}{
		"title":   strings.TrimSpace(input.Title),
		"summary": strings.TrimSpace(input.Summary),
		"content": strings.TrimSpace(input.Content),
		"cover":   strings.TrimSpace(input.Cover),
		"source":  strings.TrimSpace(input.Source),
		"author":  strings.TrimSpace(input.Author),
	}
	if input.IsPublished != nil {
		updates["is_published"] = *input.IsPublished
	}
	// 发布时间由管理员显式设置；未设置且首次发布时由仓库层兜底补当前时间。
	updates["published_at"] = input.PublishedAt

	if len(updates) > 0 {
		if err := s.newsRepo.Update(news, updates); err != nil {
			return nil, err
		}
	}

	updated, err := s.newsRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toNewsDetail(*updated), nil
}

// DeleteNews 删除新闻内容。
func (s *NewsService) DeleteNews(id uint) error {
	return s.newsRepo.Delete(id)
}

func toNewsList(news []models.News) []models.NewsListResponse {
	result := make([]models.NewsListResponse, 0, len(news))
	for _, item := range news {
		result = append(result, models.NewsListResponse{
			NewsID:      item.NewsID,
			Title:       item.Title,
			Summary:     item.Summary,
			Cover:       fullNewsURL(item.Cover),
			Source:      item.Source,
			Author:      item.Author,
			IsPublished: item.IsPublished,
			PublishedAt: item.PublishedAt,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}
	return result
}

func toNewsDetail(news models.News) *models.NewsDetailResponse {
	return &models.NewsDetailResponse{
		NewsID:      news.NewsID,
		Title:       news.Title,
		Summary:     news.Summary,
		Content:     news.Content,
		Cover:       fullNewsURL(news.Cover),
		Source:      news.Source,
		Author:      news.Author,
		IsPublished: news.IsPublished,
		PublishedAt: news.PublishedAt,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	}
}

func fullNewsURL(value string) string {
	if value == "" {
		return ""
	}
	return urlutils.GetFullURL(value)
}
