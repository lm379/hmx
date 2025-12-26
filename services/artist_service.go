package services

import (
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/models"
	"github.com/lm379/hmx/utils"
)

// GetAllArtists 获取所有艺术家 (带分页)
func GetAllArtists(pagination *utils.Pagination) ([]models.Artist, int64, error) {
	var artists []models.Artist
	var total int64
	db := database.DB.Model(&models.Artist{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Scopes(pagination.Paginate()).Find(&artists).Error
	return artists, total, err
}

// GetArtistByID 获取单个艺术家及其作品
func GetArtistByID(artistID uint) (models.Artist, error) {
	var artist models.Artist
	db := database.DB

	// Preload Operas associated with the artist
	// Also preload Artists inside Operas if we want to show other artists in that opera card?
	// For now just Operas is enough.
	err := db.Preload("Operas").Where("artist_id = ?", artistID).First(&artist).Error
	return artist, err
}
