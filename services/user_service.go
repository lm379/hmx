package services

import (
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/models"
	"github.com/lm379/hmx/utils"
)

// GetUserLikes 获取用户点赞的视频列表
func GetUserLikes(userID uint, pagination *utils.Pagination) ([]models.Opera, int64, error) {
	var operas []models.Opera
	var total int64
	db := database.DB

	// 找到用户
	// user := models.Users{UserID: userID}
	// 找到该用户点赞的所有 Opera ID
	var likedOperaIDs []uint
	db.Model(&models.Like{}).Where("user_id = ?", userID).Pluck("opera_id", &likedOperaIDs)

	if len(likedOperaIDs) == 0 {
		return operas, 0, nil
	}

	// 基于 ID 查询 Operas
	operaDB := db.Model(&models.Opera{}).Where("opera_id IN ?", likedOperaIDs)

	// 统计总数 (已批准的)
	operaDB.Where("status = ?", "approved").Count(&total)

	// 分页查询
	err := operaDB.Scopes(pagination.Paginate()).Preload("Artists").Find(&operas).Error

	return operas, total, err
}

// GetUserFavorites 获取用户收藏的视频列表
func GetUserFavorites(userID uint, pagination *utils.Pagination) ([]models.Opera, int64, error) {
	var operas []models.Opera
	var total int64
	db := database.DB

	// 逻辑与 GetUserLikes 完全相同，只是表不同
	var favOperaIDs []uint
	db.Model(&models.Favorite{}).Where("user_id = ?", userID).Pluck("opera_id", &favOperaIDs)

	if len(favOperaIDs) == 0 {
		return operas, 0, nil
	}

	operaDB := db.Model(&models.Opera{}).Where("opera_id IN ?", favOperaIDs)
	operaDB.Where("status = ?", "approved").Count(&total)
	err := operaDB.Where("status = ?", "approved").Scopes(pagination.Paginate()).Preload("Artists").Find(&operas).Error

	return operas, total, err
}
