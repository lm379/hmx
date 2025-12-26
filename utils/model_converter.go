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

	// 处理 Artists 字段
	if opera.Artists != nil {
		artists := make([]models.Artist, 0, len(opera.Artists))
		for _, artist := range opera.Artists {
			if artist != nil {
				artists = append(artists, *artist)
			}
		}
		response.Artists = artists
	} else {
		response.Artists = []models.Artist{}
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

	if opera.SrtPath.Valid && opera.SrtPath.String != "" {
		fullURL := GetFullURL(opera.SrtPath.String)
		response.SrtPath = &fullURL
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

// ToArtistResponse 将 Artist 模型转换为 ArtistResponse
func ToArtistResponse(artist models.Artist) *models.ArtistResponse {
	response := &models.ArtistResponse{
		ArtistID:  artist.ArtistID,
		Name:      artist.Name,
		Bio:       artist.Bio,
		CreatedAt: artist.CreatedAt,
		UpdatedAt: artist.UpdatedAt,
	}

	if artist.ArtistAvatar.Valid && artist.ArtistAvatar.String != "" {
		fullURL := GetFullURL(artist.ArtistAvatar.String)
		response.Avatar = &fullURL
	}

	if artist.Operas != nil {
		operas := make([]models.Opera, 0, len(artist.Operas))
		for _, op := range artist.Operas {
			if op != nil {
				operas = append(operas, *op)
			}
		}
		// 使用 ToOperaResponseList 来处理 URL
		response.Operas = ToOperaResponseList(operas)
	} else {
		response.Operas = []*models.OperaResponse{}
	}

	return response
}

// ToArtistResponseList 批量转换 Artist 列表
func ToArtistResponseList(artists []models.Artist) []*models.ArtistResponse {
	result := make([]*models.ArtistResponse, 0, len(artists))
	for _, artist := range artists {
		result = append(result, ToArtistResponse(artist))
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