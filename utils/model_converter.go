package utils

import (
	"github.com/lm379/hmx/models"
)

// ToOperaResponse 将 Opera 模型转换为 OperaResponse，自动处理 URL
func ToOperaResponse(opera models.Opera) *models.OperaResponse {
	response := &models.OperaResponse{
		OperaID:     opera.OperaID,
		OperaTitle:  opera.OperaTitle,
		VideoPath:   GetFullURL(opera.VideoPath),
		Description: opera.Description,
		AiSummary:   opera.AiSummary,
		CreatedAt:   opera.CreatedAt,
		UpdatedAt:   opera.UpdatedAt,
	}

	// 处理可空字段
	if opera.ReleaseDate.Valid {
		date := opera.ReleaseDate.Time.Format("2006-01-02")
		response.ReleaseDate = &date
	}

	if opera.Duration.Valid {
		response.Duration = &opera.Duration.String
	}

	if opera.MusicPath.Valid && opera.MusicPath.String != "" {
		fullURL := GetFullURL(opera.MusicPath.String)
		response.MusicPath = &fullURL
	}

	if opera.Avatar.Valid && opera.Avatar.String != "" {
		fullURL := GetFullURL(opera.Avatar.String)
		response.Avatar = &fullURL
	}

	return response
}

// ToOperaResponseList 批量转换 Opera 列表
func ToOperaResponseList(operas []models.Opera) []*models.OperaResponse {
	result := make([]*models.OperaResponse, 0, len(operas))
	for _, opera := range operas {
		result = append(result, ToOperaResponse(opera))
	}
	return result
}

// ToUserResponse 将 Users 模型转换为 UserResponse，自动处理 URL
func ToUserResponse(user *models.Users) *models.UserResponse {
	if user == nil {
		return nil
	}

	response := &models.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Phone:     user.Phone,
		Sex:       user.Sex,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// 处理可空字段
	if user.Email.Valid {
		response.Email = &user.Email.String
	}

	if user.Icon.Valid && user.Icon.String != "" {
		fullURL := GetFullURL(user.Icon.String)
		response.Icon = &fullURL
	}

	return response
}
