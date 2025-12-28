package services

import (
	"time"

	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/pagination"
	"gorm.io/gorm"
)

// GetUserLikes 获取用户点赞的视频列表
func GetUserLikes(userID uint, pagination *pagination.Pagination) ([]models.Opera, int64, error) {
	var operas []models.Opera
	var total int64
	db := database.DB

	var likedOperaIDs []uint
	db.Model(&models.Like{}).Where("user_id = ?", userID).Pluck("opera_id", &likedOperaIDs)

	if len(likedOperaIDs) == 0 {
		return operas, 0, nil
	}

	operaDB := db.Model(&models.Opera{}).Where("opera_id IN ?", likedOperaIDs)
	operaDB.Count(&total)

	err := operaDB.Scopes(pagination.Paginate()).Preload("Artists").Find(&operas).Error
	return operas, total, err
}

// GetUserFavorites 获取用户收藏的视频列表
func GetUserFavorites(userID uint, pagination *pagination.Pagination) ([]models.Opera, int64, error) {
	var operas []models.Opera
	var total int64
	db := database.DB

	var favOperaIDs []uint
	db.Model(&models.Favorite{}).Where("user_id = ?", userID).Pluck("opera_id", &favOperaIDs)

	if len(favOperaIDs) == 0 {
		return operas, 0, nil
	}

	operaDB := db.Model(&models.Opera{}).Where("opera_id IN ?", favOperaIDs)
	operaDB.Count(&total)
	err := operaDB.Scopes(pagination.Paginate()).Preload("Artists").Find(&operas).Error

	return operas, total, err
}

// GetUserHistory 获取用户观看历史
func GetUserHistory(userID uint, pagination *pagination.Pagination) ([]models.Opera, int64, error) {
	var operas []models.Opera
	var total int64
	db := database.DB

	// 直接查询 PlayHistory 表并关联 Opera
	// 由于 RecordPlayHistory 已经做了去重，这里直接分页查询即可
	historyDB := db.Table("play_history").
		Select("opera.*").
		Joins("join opera on opera.opera_id = play_history.opera_id").
		Where("play_history.user_id = ?", userID).
		Order("play_history.created_at desc")

	// 统计总数
	historyDB.Count(&total)

	// 分页并加载关联的 Artists
	err := historyDB.Scopes(pagination.Paginate()).
		Preload("Artists").
		Find(&operas).Error

	// 如果 Preload 失败（因为使用了 Table("play_history")），可以改回使用模型查询
	if err != nil {
		// 回退方案：先查 ID，再查模型
		var operaIDs []uint
		db.Model(&models.PlayHistory{}).
			Where("user_id = ?", userID).
			Order("created_at desc").
			Scopes(pagination.Paginate()).
			Pluck("opera_id", &operaIDs)

		if len(operaIDs) == 0 {
			return []models.Opera{}, total, nil
		}

		err = db.Preload("Artists").Where("opera_id IN ?", operaIDs).Find(&operas).Error
		// 再次手动排序
		operaMap := make(map[uint]models.Opera)
		for _, op := range operas {
			operaMap[op.OperaID] = op
		}

		operas = make([]models.Opera, 0, len(operaIDs))
		for _, id := range operaIDs {
			if op, ok := operaMap[id]; ok {
				operas = append(operas, op)
			}
		}
	}

	return operas, total, err
}

// RecordPlayHistory 记录播放历史
func RecordPlayHistory(userID *uint, operaID uint) error {
	db := database.DB

	// 检查作品是否存在
	var opera models.Opera
	if err := db.First(&opera, operaID).Error; err != nil {
		return err
	}

	// 如果是登录用户，尝试更新已存在的记录（去重，只保留最新）
	if userID != nil {
		var existing models.PlayHistory
		err := db.Where("user_id = ? AND opera_id = ?", *userID, operaID).First(&existing).Error
		if err == nil {
			// 找到了 -> 增加播放次数并更新时间到当前
			return db.Model(&existing).Updates(map[string]interface{}{
				"count":      existing.Count + 1,
				"created_at": time.Now(),
			}).Error
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	// 没找到或为游客 -> 创建新记录
	history := models.PlayHistory{
		UserID:  userID,
		OperaID: operaID,
		Count:   1, // 初始播放次数为 1
	}

	return db.Create(&history).Error
}
