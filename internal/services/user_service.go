package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/hashutils"
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

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(userID uint, input models.UpdateUserProfileRequest) error {
	db := database.DB
	var user models.Users
	if err := db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	updates := make(map[string]interface{})

	// Username
	if input.Username != "" && input.Username != user.Username {
		var count int64
		db.Model(&models.Users{}).Where("username = ?", input.Username).Count(&count)
		if count > 0 {
			return &ServiceError{Code: http.StatusConflict, Message: "Username already taken"}
		}
		updates["username"] = input.Username
	}

	// Phone
	if input.Phone != "" && input.Phone != user.Phone {
		var count int64
		db.Model(&models.Users{}).Where("phone = ?", input.Phone).Count(&count)
		if count > 0 {
			return &ServiceError{Code: http.StatusConflict, Message: "Phone number already in use"}
		}
		updates["phone"] = input.Phone
	}

	// Sex
	if input.Sex != "" {
		updates["sex"] = input.Sex
	}

	// Icon
	if input.Icon != "" {
		updates["icon"] = input.Icon
	}

	// Email (Need Verification)
	if input.Email != "" && input.Email != user.Email.String {
		// Check uniqueness
		var count int64
		db.Model(&models.Users{}).Where("email = ?", input.Email).Count(&count)
		if count > 0 {
			return &ServiceError{Code: http.StatusConflict, Message: "Email already in use"}
		}

		// Verify Code against CURRENT email (user must prove ownership of account to change email)
		if !VerifyCode(user.Email.String, input.Code) {
			return &ServiceError{Code: http.StatusBadRequest, Message: "Invalid verification code"}
		}
		// Invalidate code
		DeleteCode(user.Email.String)

		updates["email"] = input.Email
	}

	if len(updates) == 0 {
		return nil
	}

	return db.Model(&user).Updates(updates).Error
}

// UpdatePassword 更新密码
func UpdatePassword(userID uint, input models.UpdatePasswordRequest) error {
	db := database.DB
	var user models.Users
	if err := db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Verify Old Password
	if !hashutils.CheckPasswordHash(input.OldPassword, user.Password) {
		return &ServiceError{Code: http.StatusUnauthorized, Message: "Incorrect old password"}
	}

	hashedPassword, err := hashutils.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	return db.Model(&user).Update("password", hashedPassword).Error
}
