package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/services"
	"github.com/lm379/hmx/utils"
	"gorm.io/gorm"
)

// HandleGetArtists (GET /api/v1/artists)
func HandleGetArtists(c *gin.Context) {
	pagination := utils.GetPagination(c)

	artists, total, err := services.GetAllArtists(pagination)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch artists")
		return
	}

	artistResponses := utils.ToArtistResponseList(artists)

	utils.Success(c, gin.H{
		"list": artistResponses,
		"pagination": gin.H{
			"total":     total,
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
		},
	})
}

// HandleGetArtistByID (GET /api/v1/artists/:id)
func HandleGetArtistByID(c *gin.Context) {
	idParam := c.Param("id")
	artistID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid artist ID")
		return
	}

	artist, err := services.GetArtistByID(uint(artistID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Artist not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	artistResponse := utils.ToArtistResponse(artist)
	utils.Success(c, artistResponse)
}
