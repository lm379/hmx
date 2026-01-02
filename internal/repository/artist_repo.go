package repository

import (
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/pagination"
	"gorm.io/gorm"
)

// ArtistRepo 艺术家相关的数据库操作
type ArtistRepo struct{}

// NewArtistRepo 创建艺术家仓库实例
func NewArtistRepo() *ArtistRepo {
	return &ArtistRepo{}
}

// getDB 获取数据库实例
func (r *ArtistRepo) getDB() *gorm.DB {
	return database.DB
}

// GetAll 获取所有艺术家 (带分页)
func (r *ArtistRepo) GetAll(pagination *pagination.Pagination) ([]models.Artist, int64, error) {
	var artists []models.Artist
	var total int64
	db := r.getDB().Model(&models.Artist{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Scopes(pagination.Paginate()).Find(&artists).Error
	return artists, total, err
}

// GetByID 根据ID获取艺术家
func (r *ArtistRepo) GetByID(id uint) (*models.Artist, error) {
	var artist models.Artist
	err := r.getDB().Preload("Operas").Where("artist_id = ?", id).First(&artist).Error
	if err != nil {
		return nil, err
	}
	return &artist, nil
}

// Create 创建艺术家
func (r *ArtistRepo) Create(artist *models.Artist) error {
	return r.getDB().Create(artist).Error
}

// Update 更新艺术家
func (r *ArtistRepo) Update(id uint, updates map[string]interface{}) error {
	return r.getDB().Model(&models.Artist{}).Where("artist_id = ?", id).Updates(updates).Error
}

// Delete 删除艺术家
func (r *ArtistRepo) Delete(id uint) error {
	return r.getDB().Delete(&models.Artist{}, id).Error
}

// Exists 检查艺术家是否存在
func (r *ArtistRepo) Exists(id uint) (bool, error) {
	var count int64
	err := r.getDB().Model(&models.Artist{}).Where("artist_id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByName 根据名字查找艺术家
func (r *ArtistRepo) FindByName(name string) (*models.Artist, error) {
	var artist models.Artist
	err := r.getDB().Where("name = ?", name).First(&artist).Error
	if err != nil {
		return nil, err
	}
	return &artist, nil
}
