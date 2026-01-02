package services

import (
	"errors"
	"net/http"

	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/hashutils"
	"github.com/lm379/hmx/pkg/pagination"
	"github.com/lm379/hmx/pkg/validator"
)

// GetUserLikes 获取用户点赞的视频列表
func GetUserLikes(userID uint, pagination *pagination.Pagination) ([]models.Opera, int64, error) {
	interactionRepo := repository.NewInteractionRepo()
	operaRepo := repository.NewOperaRepo()

	likedOperaIDs, err := interactionRepo.GetUserLikedOperaIDs(userID)
	if err != nil {
		return nil, 0, err
	}

	if len(likedOperaIDs) == 0 {
		return []models.Opera{}, 0, nil
	}

	// 使用 opera repo 获取作品列表
	operas, total, err := operaRepo.GetByIDs(likedOperaIDs, pagination)
	return operas, total, err
}

// GetUserFavorites 获取用户收藏的视频列表
func GetUserFavorites(userID uint, pagination *pagination.Pagination) ([]models.Opera, int64, error) {
	interactionRepo := repository.NewInteractionRepo()
	operaRepo := repository.NewOperaRepo()

	favOperaIDs, err := interactionRepo.GetUserFavoritedOperaIDs(userID)
	if err != nil {
		return nil, 0, err
	}

	if len(favOperaIDs) == 0 {
		return []models.Opera{}, 0, nil
	}

	// 使用 opera repo 获取作品列表
	operas, total, err := operaRepo.GetByIDs(favOperaIDs, pagination)
	return operas, total, err
}

// GetUserHistory 获取用户观看历史
func GetUserHistory(userID uint, pagination *pagination.Pagination) ([]models.Opera, int64, error) {
	interactionRepo := repository.NewInteractionRepo()
	operaRepo := repository.NewOperaRepo()

	// 获取观看历史的作品ID列表
	offset := (pagination.Page - 1) * pagination.PageSize
	operaIDs, total, err := interactionRepo.GetUserHistoryOperaIDs(userID, pagination.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	if len(operaIDs) == 0 {
		return []models.Opera{}, 0, nil
	}

	// 根据ID列表获取作品，保持顺序
	operas, err := operaRepo.GetByIDsInOrder(operaIDs)
	return operas, total, err
}

// RecordPlayHistory 记录播放历史
func RecordPlayHistory(userID *uint, operaID uint) error {
	operaRepo := repository.NewOperaRepo()
	interactionRepo := repository.NewInteractionRepo()

	// 检查作品是否存在
	exists, err := operaRepo.Exists(operaID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("opera not found")
	}

	// 记录播放历史
	return interactionRepo.RecordOrUpdatePlayHistory(userID, operaID)
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(userID uint, input models.UpdateUserProfileRequest) error {
	userRepo := repository.NewUserRepo()
	
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	updates := make(map[string]interface{})

	// Username
	if input.Username != "" && input.Username != user.Username {
		exists, err := userRepo.UsernameExists(input.Username)
		if err != nil {
			return err
		}
		if exists {
			return &ServiceError{Code: http.StatusConflict, Message: "Username already taken"}
		}
		updates["username"] = input.Username
	}

	// Phone
	if input.Phone != "" && input.Phone != user.Phone {
		if !validator.ValidatePhone(input.Phone) {
			return &ServiceError{Code: http.StatusBadRequest, Message: "Invalid phone number format"}
		}
		exists, err := userRepo.PhoneExists(input.Phone)
		if err != nil {
			return err
		}
		if exists {
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
		if !validator.ValidateEmail(input.Email) {
			return &ServiceError{Code: http.StatusBadRequest, Message: "Invalid email format"}
		}
		// Check uniqueness
		exists, err := userRepo.EmailExists(input.Email)
		if err != nil {
			return err
		}
		if exists {
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

	return userRepo.Update(userID, updates)
}

// UpdatePassword 更新密码
func UpdatePassword(userID uint, input models.UpdatePasswordRequest) error {
	userRepo := repository.NewUserRepo()
	
	user, err := userRepo.GetByID(userID)
	if err != nil {
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

	return userRepo.UpdatePassword(userID, hashedPassword)
}

// GetUserByID 根据ID获取用户
func GetUserByID(userID uint) (*models.Users, error) {
	userRepo := repository.NewUserRepo()
	return userRepo.GetByID(userID)
}

// GetAllUsers 获取所有用户 (Admin)
func GetAllUsers(pagination *pagination.Pagination) ([]models.Users, int64, error) {
	userRepo := repository.NewUserRepo()
	return userRepo.GetAll(pagination)
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(userID uint, role models.UserRole) error {
	userRepo := repository.NewUserRepo()
	return userRepo.UpdateRole(userID, string(role))
}
