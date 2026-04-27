package repository

import (
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/pagination"
	"gorm.io/gorm"
)

type EducationRepo struct{}

func NewEducationRepo() *EducationRepo {
	return &EducationRepo{}
}

func (r *EducationRepo) getDB() *gorm.DB {
	return database.DB
}

// GetPublished 返回前台可见的 PDF 书籍。
func (r *EducationRepo) GetPublished(p *pagination.Pagination) ([]models.EducationBook, int64, error) {
	var list []models.EducationBook
	var total int64
	db := r.getDB().Model(&models.EducationBook{}).Where("is_published = ?", true)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("created_at DESC, book_id DESC").
		Offset((p.Page - 1) * p.PageSize).
		Limit(p.PageSize).
		Find(&list).Error
	return list, total, err
}

// GetAll 返回后台管理列表，包含草稿资源。
func (r *EducationRepo) GetAll(p *pagination.Pagination) ([]models.EducationBook, int64, error) {
	var list []models.EducationBook
	var total int64
	db := r.getDB().Model(&models.EducationBook{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("created_at DESC, book_id DESC").
		Offset((p.Page - 1) * p.PageSize).
		Limit(p.PageSize).
		Find(&list).Error
	return list, total, err
}

// GetPublishedByID 用于前台 PDF 阅读页，避免草稿被直接访问。
func (r *EducationRepo) GetPublishedByID(id uint) (*models.EducationBook, error) {
	var book models.EducationBook
	err := r.getDB().Where("book_id = ? AND is_published = ?", id, true).First(&book).Error
	return &book, err
}

// GetByID 用于后台编辑，不限制发布状态。
func (r *EducationRepo) GetByID(id uint) (*models.EducationBook, error) {
	var book models.EducationBook
	err := r.getDB().Where("book_id = ?", id).First(&book).Error
	return &book, err
}

// Create 保存教育 PDF 书籍。
func (r *EducationRepo) Create(book *models.EducationBook) error {
	return r.getDB().Create(book).Error
}

// Update 更新教育 PDF 书籍。
func (r *EducationRepo) Update(book *models.EducationBook, updates map[string]interface{}) error {
	return r.getDB().Model(book).Updates(updates).Error
}

// Delete 删除教育 PDF 书籍记录。
func (r *EducationRepo) Delete(id uint) error {
	return r.getDB().Where("book_id = ?", id).Delete(&models.EducationBook{}).Error
}
