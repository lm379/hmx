package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/converter"
	"github.com/lm379/hmx/pkg/pagination"
	resp "github.com/lm379/hmx/pkg/response"
	"gorm.io/gorm"
)

// HandleGetArtists (GET /api/v1/artists)
func HandleGetArtists(c *gin.Context) {
	pagination := pagination.GetPagination(c)

	artists, total, err := services.GetAllArtists(pagination)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch artists")
		return
	}

	artistResponses := converter.ToArtistResponseList(artists)

	resp.Success(c, gin.H{
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
		resp.BadRequest(c, "Invalid artist ID")
		return
	}

	artist, err := services.GetArtistByID(uint(artistID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "Artist not found")
			return
		}
		resp.InternalServerError(c, "Database error")
		return
	}

	artistResponse := converter.ToArtistResponse(artist)
	resp.Success(c, artistResponse)
}
