package repository

import (
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/pagination"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

// OperaRepo 作品相关的数据库操作
type OperaRepo struct{}

// NewOperaRepo 创建作品仓库实例
func NewOperaRepo() *OperaRepo {
	return &OperaRepo{}
}

// getDB 获取数据库实例
func (r *OperaRepo) getDB() *gorm.DB {
	return database.DB
}

// GetAll 获取所有作品 (带分页和过滤)
func (r *OperaRepo) GetAll(pagination *pagination.Pagination, includeHidden bool) ([]models.Opera, int64, error) {
	var operas []models.Opera
	var total int64
	db := r.getDB().Model(&models.Opera{})

	if !includeHidden {
		db = db.Where("is_hidden = ?", false)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 默认按播放量从高到低排序
	offset := (pagination.Page - 1) * pagination.PageSize

	err := db.Select("opera.*").
		Joins("LEFT JOIN play_history ON opera.opera_id = play_history.opera_id").
		Group("opera.opera_id").
		Order("COALESCE(SUM(play_history.count), 0) DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Preload("Artists").
		Find(&operas).Error

	return operas, total, err
}

// GetByID 根据ID获取作品
func (r *OperaRepo) GetByID(id uint) (*models.Opera, error) {
	var opera models.Opera
	err := r.getDB().Preload("Artists").First(&opera, id).Error
	if err != nil {
		return nil, err
	}
	return &opera, nil
}

// Create 创建作品
func (r *OperaRepo) Create(opera *models.Opera) error {
	return r.getDB().Omit("embedding").Create(opera).Error
}

// Update 更新作品基本信息
func (r *OperaRepo) Update(id uint, updates map[string]interface{}) error {
	return r.getDB().Model(&models.Opera{}).Where("opera_id = ?", id).Updates(updates).Error
}

// UpdateEmbedding 更新作品的向量字段
func (r *OperaRepo) UpdateEmbedding(id uint, embedding []float64) error {
	// 转换为float32数组
	vec32 := make([]float32, len(embedding))
	for i, v := range embedding {
		vec32[i] = float32(v)
	}
	// 转换为pgvector.Vector类型
	vec := pgvector.NewVector(vec32)
	return r.getDB().Model(&models.Opera{}).Where("opera_id = ?", id).Update("embedding", vec).Error
}

// UpdateSummary 更新作品的AI摘要字段
func (r *OperaRepo) UpdateSummary(id uint, summary string) error {
	return r.getDB().Model(&models.Opera{}).Where("opera_id = ?", id).Update("ai_summary", summary).Error
}

// UpdateArtists 更新作品的艺术家关联
func (r *OperaRepo) UpdateArtists(opera *models.Opera, artists []*models.Artist) error {
	return r.getDB().Model(opera).Association("Artists").Replace(artists)
}

// Delete 删除作品
func (r *OperaRepo) Delete(id uint) error {
	return r.getDB().Delete(&models.Opera{}, id).Error
}

// Exists 检查作品是否存在
func (r *OperaRepo) Exists(id uint) (bool, error) {
	var count int64
	err := r.getDB().Model(&models.Opera{}).Where("opera_id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindArtistsByIDs 根据ID列表查找艺术家
func (r *OperaRepo) FindArtistsByIDs(ids []uint) ([]*models.Artist, error) {
	var artists []*models.Artist
	err := r.getDB().Find(&artists, ids).Error
	return artists, err
}

// GetByIDs 根据ID列表获取作品（带分页）
func (r *OperaRepo) GetByIDs(ids []uint, pagination *pagination.Pagination) ([]models.Opera, int64, error) {
	var operas []models.Opera
	var total int64

	db := r.getDB().Model(&models.Opera{}).Where("opera_id IN ?", ids)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Scopes(pagination.Paginate()).Preload("Artists").Find(&operas).Error
	return operas, total, err
}

// GetByIDsInOrder 根据ID列表获取作品（保持顺序）
func (r *OperaRepo) GetByIDsInOrder(ids []uint) ([]models.Opera, error) {
	var operas []models.Opera
	err := r.getDB().Preload("Artists").Where("opera_id IN ?", ids).Find(&operas).Error
	if err != nil {
		return nil, err
	}

	// 按照输入的ID顺序排序
	operaMap := make(map[uint]models.Opera)
	for _, op := range operas {
		operaMap[op.OperaID] = op
	}

	orderedOperas := make([]models.Opera, 0, len(ids))
	for _, id := range ids {
		if op, ok := operaMap[id]; ok {
			orderedOperas = append(orderedOperas, op)
		}
	}

	return orderedOperas, nil
}

// FindByTitle 根据作品标题查找作品（用于转码回调）
func (r *OperaRepo) FindByTitle(title string) (*models.Opera, error) {
	var opera models.Opera
	err := r.getDB().Preload("Artists").Where("opera_title = ?", title).First(&opera).Error
	if err != nil {
		return nil, err
	}
	return &opera, nil
}

// GetPopularOperas 获取热门作品（基于兴趣分）
// excludeIDs: 需要排除的作品ID列表
// limit: 返回数量
func (r *OperaRepo) GetPopularOperas(excludeIDs []uint, limit int) ([]uint, error) {
	var operaIDs []uint

	db := r.getDB().Table("opera op").
		Select("op.opera_id").
		Joins("LEFT JOIN (SELECT opera_id, SUM(total_score) as interest_score FROM user_interest_scores GROUP BY opera_id) scores ON op.opera_id = scores.opera_id").
		Where("op.is_hidden = ?", false)

	if len(excludeIDs) > 0 {
		db = db.Where("op.opera_id NOT IN ?", excludeIDs)
	}

	err := db.Order("COALESCE(scores.interest_score, 0) DESC, op.opera_id ASC").
		Limit(limit).
		Pluck("opera_id", &operaIDs).Error

	return operaIDs, err
}

// GetOperasNeedingSummary 获取所有有字幕但没有摘要的作品
func (r *OperaRepo) GetOperasNeedingSummary() ([]models.Opera, error) {
	var operas []models.Opera
	err := r.getDB().Where("srt_path != '' AND srt_path IS NOT NULL AND (ai_summary = '' OR ai_summary IS NULL)").Find(&operas).Error
	return operas, err
}

// CountVisible 获取可见（非隐藏）作品的总数
func (r *OperaRepo) CountVisible() (int64, error) {
	var count int64
	err := r.getDB().Model(&models.Opera{}).Where("is_hidden = ?", false).Count(&count).Error
	return count, err
}
