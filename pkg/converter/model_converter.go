package converter

import (
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/urlutils"
)

// ToOperaResponse 将 Opera 模型转换为 OperaResponse，自动处理 URL
func ToOperaResponse(opera models.Opera) *models.OperaResponse {
	response := &models.OperaResponse{
		OperaID:     opera.OperaID,
		OperaTitle:  opera.OperaTitle,
		VideoPath:   urlutils.GetFullURL(opera.VideoPath),
		Description: opera.Description,
		AiSummary:   opera.AiSummary,
		IsHidden:    opera.IsHidden,
		CreatedAt:   opera.CreatedAt,
		UpdatedAt:   opera.UpdatedAt,
	}

	// 处理 Artists 字段
	if opera.Artists != nil {
		artists := make([]models.SimpleArtist, 0, len(opera.Artists))
		for _, artist := range opera.Artists {
			if artist != nil {
				simpleArtist := models.SimpleArtist{
					ArtistID: artist.ArtistID,
					Name:     artist.Name,
				}
				if artist.ArtistAvatar.Valid && artist.ArtistAvatar.String != "" {
					simpleArtist.Avatar = urlutils.GetFullURL(artist.ArtistAvatar.String)
				}
				artists = append(artists, simpleArtist)
			}
		}
		response.Artists = artists
	} else {
		response.Artists = []models.SimpleArtist{}
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
		fullURL := urlutils.GetFullURL(opera.MusicPath.String)
		response.MusicPath = &fullURL
	}

	if opera.SrtPath.Valid && opera.SrtPath.String != "" {
		fullURL := urlutils.GetFullURL(opera.SrtPath.String)
		response.SrtPath = &fullURL
	}

	if opera.Avatar.Valid && opera.Avatar.String != "" {
		fullURL := urlutils.GetFullURL(opera.Avatar.String)
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
		fullURL := urlutils.GetFullURL(artist.ArtistAvatar.String)
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

// ToOperaListResponse 将 Opera 转换为列表响应（简化版）
func ToOperaListResponse(opera models.Opera) *models.OperaListResponse {
	response := &models.OperaListResponse{
		OperaID:    opera.OperaID,
		OperaTitle: opera.OperaTitle,
		IsHidden:   opera.IsHidden,
		CreatedAt:  opera.CreatedAt,
	}

	// 处理 Artists 字段
	if opera.Artists != nil {
		artists := make([]models.SimpleArtist, 0, len(opera.Artists))
		for _, artist := range opera.Artists {
			if artist != nil {
				simpleArtist := models.SimpleArtist{
					ArtistID: artist.ArtistID,
					Name:     artist.Name,
				}
				if artist.ArtistAvatar.Valid && artist.ArtistAvatar.String != "" {
					simpleArtist.Avatar = urlutils.GetFullURL(artist.ArtistAvatar.String)
				}
				artists = append(artists, simpleArtist)
			}
		}
		response.Artists = artists
	} else {
		response.Artists = []models.SimpleArtist{}
	}

	if opera.Avatar.Valid && opera.Avatar.String != "" {
		fullURL := urlutils.GetFullURL(opera.Avatar.String)
		response.Avatar = &fullURL
	}

	if opera.Duration.Valid {
		response.Duration = &opera.Duration.String
	}

	return response
}

// ToOperaListResponseList 批量转换为列表响应
func ToOperaListResponseList(operas []models.Opera) []*models.OperaListResponse {
	result := make([]*models.OperaListResponse, 0, len(operas))
	for _, opera := range operas {
		result = append(result, ToOperaListResponse(opera))
	}
	return result
}

// ToOperaDetailResponse 将 Opera 转换为详情响应（完整版）
func ToOperaDetailResponse(opera models.Opera) *models.OperaDetailResponse {
	response := &models.OperaDetailResponse{
		OperaID:     opera.OperaID,
		OperaTitle:  opera.OperaTitle,
		VideoPath:   urlutils.GetFullURL(opera.VideoPath),
		Description: opera.Description,
		AiSummary:   opera.AiSummary,
		IsHidden:    opera.IsHidden,
		CreatedAt:   opera.CreatedAt,
		UpdatedAt:   opera.UpdatedAt,
	}

	// 处理 Artists 字段
	if opera.Artists != nil {
		artists := make([]models.SimpleArtist, 0, len(opera.Artists))
		for _, artist := range opera.Artists {
			if artist != nil {
				simpleArtist := models.SimpleArtist{
					ArtistID: artist.ArtistID,
					Name:     artist.Name,
				}
				if artist.ArtistAvatar.Valid && artist.ArtistAvatar.String != "" {
					simpleArtist.Avatar = urlutils.GetFullURL(artist.ArtistAvatar.String)
				}
				artists = append(artists, simpleArtist)
			}
		}
		response.Artists = artists
	} else {
		response.Artists = []models.SimpleArtist{}
	}

	if opera.ReleaseDate.Valid {
		date := opera.ReleaseDate.Time.Format("2006-01-02")
		response.ReleaseDate = &date
	}

	if opera.Duration.Valid {
		response.Duration = &opera.Duration.String
	}

	if opera.MusicPath.Valid && opera.MusicPath.String != "" {
		fullURL := urlutils.GetFullURL(opera.MusicPath.String)
		response.MusicPath = &fullURL
	}

	if opera.SrtPath.Valid && opera.SrtPath.String != "" {
		fullURL := urlutils.GetFullURL(opera.SrtPath.String)
		response.SrtPath = &fullURL
	}

	if opera.Avatar.Valid && opera.Avatar.String != "" {
		fullURL := urlutils.GetFullURL(opera.Avatar.String)
		response.Avatar = &fullURL
	}

	return response
}

// ToSimpleOpera 将 Opera 转换为简化版（用于艺术家详情，不包含艺术家信息避免冗余）
func ToSimpleOpera(opera models.Opera) models.SimpleOpera {
	simple := models.SimpleOpera{
		OperaID:    opera.OperaID,
		OperaTitle: opera.OperaTitle,
		CreatedAt:  opera.CreatedAt,
	}

	if opera.Avatar.Valid && opera.Avatar.String != "" {
		fullURL := urlutils.GetFullURL(opera.Avatar.String)
		simple.Avatar = &fullURL
	}

	if opera.Duration.Valid {
		simple.Duration = &opera.Duration.String
	}

	return simple
}

// ToSimpleOperaList 批量转换为简化版
func ToSimpleOperaList(operas []models.Opera) []models.SimpleOpera {
	result := make([]models.SimpleOpera, 0, len(operas))
	for _, opera := range operas {
		result = append(result, ToSimpleOpera(opera))
	}
	return result
}

// ToArtistListResponse 将 Artist 转换为列表响应（简化版）
func ToArtistListResponse(artist models.Artist) *models.ArtistListResponse {
	response := &models.ArtistListResponse{
		ArtistID:  artist.ArtistID,
		Name:      artist.Name,
		Bio:       artist.Bio,
		CreatedAt: artist.CreatedAt,
	}

	if artist.ArtistAvatar.Valid && artist.ArtistAvatar.String != "" {
		fullURL := urlutils.GetFullURL(artist.ArtistAvatar.String)
		response.Avatar = &fullURL
	}

	return response
}

// ToArtistListResponseList 批量转换为列表响应
func ToArtistListResponseList(artists []models.Artist) []*models.ArtistListResponse {
	result := make([]*models.ArtistListResponse, 0, len(artists))
	for _, artist := range artists {
		result = append(result, ToArtistListResponse(artist))
	}
	return result
}

// ToArtistDetailResponse 将 Artist 转换为详情响应（完整版）
func ToArtistDetailResponse(artist models.Artist) *models.ArtistDetailResponse {
	response := &models.ArtistDetailResponse{
		ArtistID:  artist.ArtistID,
		Name:      artist.Name,
		Bio:       artist.Bio,
		CreatedAt: artist.CreatedAt,
		UpdatedAt: artist.UpdatedAt,
	}

	if artist.ArtistAvatar.Valid && artist.ArtistAvatar.String != "" {
		fullURL := urlutils.GetFullURL(artist.ArtistAvatar.String)
		response.Avatar = &fullURL
	}

	if artist.Operas != nil {
		operas := make([]models.Opera, 0, len(artist.Operas))
		for _, op := range artist.Operas {
			if op != nil {
				operas = append(operas, *op)
			}
		}
		response.Operas = ToSimpleOperaList(operas)
	} else {
		response.Operas = []models.SimpleOpera{}
	}

	return response
}

// ToUserResponse 将 Users 模型转换为 UserResponse，自动处理 URL
func ToUserResponse(user *models.Users) *models.UserResponse {
	if user == nil {
		return nil
	}

	response := &models.UserResponse{
		UserID:      user.UserID,
		Username:    user.Username,
		Phone:       user.Phone,
		Sex:         user.Sex,
		Role:        user.Role,
		LastLoginAt: user.LastLoginAt,
		LastLoginIP: user.LastIp,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	// 处理可空字段
	if user.Email.Valid {
		response.Email = &user.Email.String
	}

	if user.Icon.Valid && user.Icon.String != "" {
		fullURL := urlutils.GetFullURL(user.Icon.String)
		response.Icon = &fullURL
	}

	return response
}
