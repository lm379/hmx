package repository

import (
	"time"

	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/pagination"
	"gorm.io/gorm"
)

type NewsRepo struct{}

func NewNewsRepo() *NewsRepo {
	return &NewsRepo{}
}

func (r *NewsRepo) getDB() *gorm.DB {
	return database.DB
}

// GetPublished 只返回前台可见的新闻，并按发布时间倒序排列。
func (r *NewsRepo) GetPublished(p *pagination.Pagination) ([]models.News, int64, error) {
	var list []models.News
	var total int64
	db := r.getDB().Model(&models.News{}).Where("is_published = ?", true)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("COALESCE(published_at, created_at) DESC, news_id DESC").
		Offset((p.Page - 1) * p.PageSize).
		Limit(p.PageSize).
		Find(&list).Error
	return list, total, err
}

// GetAll 返回后台管理列表，包含草稿和已发布新闻。
func (r *NewsRepo) GetAll(p *pagination.Pagination) ([]models.News, int64, error) {
	var list []models.News
	var total int64
	db := r.getDB().Model(&models.News{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("created_at DESC, news_id DESC").
		Offset((p.Page - 1) * p.PageSize).
		Limit(p.PageSize).
		Find(&list).Error
	return list, total, err
}

// GetPublishedByID 用于前台详情页，避免未发布新闻被直接访问。
func (r *NewsRepo) GetPublishedByID(id uint) (*models.News, error) {
	var news models.News
	err := r.getDB().Where("news_id = ? AND is_published = ?", id, true).First(&news).Error
	return &news, err
}

// GetByID 用于后台管理，不限制发布状态。
func (r *NewsRepo) GetByID(id uint) (*models.News, error) {
	var news models.News
	err := r.getDB().Where("news_id = ?", id).First(&news).Error
	return &news, err
}

// Create 保存新闻；直接发布时补齐 published_at。
func (r *NewsRepo) Create(news *models.News) error {
	if news.IsPublished && news.PublishedAt == nil {
		now := time.Now()
		news.PublishedAt = &now
	}
	return r.getDB().Create(news).Error
}

// Update 更新新闻；首次从草稿改为发布时补齐 published_at。
func (r *NewsRepo) Update(news *models.News, updates map[string]interface{}) error {
	if published, ok := updates["is_published"].(bool); ok && published && news.PublishedAt == nil {
		updates["published_at"] = time.Now()
	}
	return r.getDB().Model(news).Updates(updates).Error
}

// Delete 删除新闻记录。
func (r *NewsRepo) Delete(id uint) error {
	return r.getDB().Where("news_id = ?", id).Delete(&models.News{}).Error
}
