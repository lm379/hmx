package services

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/pagination"
)

var artistRepo = repository.NewArtistRepo()

// GetAllArtists 获取所有艺术家 (带分页)
func GetAllArtists(pagination *pagination.Pagination) ([]models.Artist, int64, error) {
	return artistRepo.GetAll(pagination)
}

// GetArtistByID 获取单个艺术家及其作品
func GetArtistByID(artistID uint) (models.Artist, error) {
	artist, err := artistRepo.GetByID(artistID)
	if err != nil {
		return models.Artist{}, err
	}
	return *artist, nil
}

// CreateArtist 创建艺术家
func CreateArtist(input models.CreateArtistRequest) (gin.H, int) {
	artist := models.Artist{
		Name: input.Name,
		Bio:  input.Bio,
		ArtistAvatar: sql.NullString{
			String: input.Avatar,
			Valid:  input.Avatar != "",
		},
	}

	if err := artistRepo.Create(&artist); err != nil {
		return gin.H{"error": "Failed to create artist: " + err.Error()}, http.StatusInternalServerError
	}

	return gin.H{"message": "Artist created successfully", "artist_id": artist.ArtistID}, http.StatusCreated
}

// UpdateArtist 更新艺术家
func UpdateArtist(id uint, input models.UpdateArtistRequest) (gin.H, int) {
	// 检查艺术家是否存在
	exists, err := artistRepo.Exists(id)
	if err != nil {
		return gin.H{"error": "Failed to check artist"}, http.StatusInternalServerError
	}
	if !exists {
		return gin.H{"error": "Artist not found"}, http.StatusNotFound
	}

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Bio != "" {
		updates["bio"] = input.Bio
	}
	if input.Avatar != "" {
		updates["artist_avatar"] = input.Avatar
	}

	if len(updates) > 0 {
		if err := artistRepo.Update(id, updates); err != nil {
			return gin.H{"error": "Failed to update artist"}, http.StatusInternalServerError
		}
	}

	return gin.H{"message": "Artist updated successfully"}, http.StatusOK
}

// DeleteArtist 删除艺术家
func DeleteArtist(id uint) (gin.H, int) {
	if err := artistRepo.Delete(id); err != nil {
		return gin.H{"error": "Failed to delete artist"}, http.StatusInternalServerError
	}
	return gin.H{"message": "Artist deleted successfully"}, http.StatusOK
}
