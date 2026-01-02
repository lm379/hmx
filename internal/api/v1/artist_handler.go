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

	// 使用列表响应（不包含作品列表）
	artistResponses := converter.ToArtistListResponseList(artists)

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

	// 使用详情响应（包含简化的作品列表）
	artistResponse := converter.ToArtistDetailResponse(artist)
	// 为每个作品添加统计数据
	if len(artistResponse.Operas) > 0 {
		operaIDs := make([]uint, len(artistResponse.Operas))
		for i, op := range artistResponse.Operas {
			operaIDs[i] = op.OperaID
		}
		likes, _, _, plays := services.BatchGetCounts(operaIDs)
		for i := range artistResponse.Operas {
			artistResponse.Operas[i].PlayCount = plays[artistResponse.Operas[i].OperaID]
			artistResponse.Operas[i].LikeCount = likes[artistResponse.Operas[i].OperaID]
		}
	}
	resp.Success(c, artistResponse)
}
