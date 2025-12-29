package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/converter"
	"github.com/lm379/hmx/pkg/pagination"
	resp "github.com/lm379/hmx/pkg/response"
	"gorm.io/gorm"
)

// HandleCreateOpera (POST /api/v1/operas)
func HandleCreateOpera(c *gin.Context) {
	var input models.CreateOperaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	// 从 JWT 中间件获取 userID
	userID, exists := c.Get("userID")
	if !exists {
		resp.Unauthorized(c, "Unauthorized")
		return
	}

	// 调用 service 创建 Opera
	response, status := services.CreateOpera(input, userID.(uint))
	if status != http.StatusCreated {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Created(c, response)
}

// HandleGetOperas (GET /api/v1/operas)
func HandleGetOperas(c *gin.Context) {
	pagination := pagination.GetPagination(c)

	operas, total, err := services.GetAllOperas(pagination)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch operas")
		return
	}

	// 转换为响应格式，自动处理 URL
	operaResponses := converter.ToOperaResponseList(operas)

	// 批量获取计数，避免N+1查询问题
	if len(operaResponses) > 0 {
		operaIDs := make([]uint, len(operaResponses))
		for i, r := range operaResponses {
			operaIDs[i] = r.OperaID
		}

		likes, favorites, shares, plays := services.BatchGetCounts(operaIDs)
		for _, r := range operaResponses {
			r.LikeCount = likes[r.OperaID]
			r.FavoriteCount = favorites[r.OperaID]
			r.ShareCount = shares[r.OperaID]
			r.PlayCount = plays[r.OperaID]
		}
	}

	resp.Success(c, gin.H{
		"list": operaResponses,
		"pagination": gin.H{
			"total":     total,
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
		},
	})
}

// HandleGetOperaByID (GET /api/v1/operas/:id)
func HandleGetOperaByID(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	opera, err := services.GetOperaByID(uint(operaID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "Opera not found")
			return
		}
		resp.InternalServerError(c, "System error")
		return
	}

	// 转换为响应格式，自动处理 URL
	operaResponse := converter.ToOperaResponse(opera)
	// 计数与用户状态
	lc, fc, sc, pc := services.GetCounts(operaResponse.OperaID)
	operaResponse.LikeCount = lc
	operaResponse.FavoriteCount = fc
	operaResponse.ShareCount = sc
	operaResponse.PlayCount = pc
	if v, exists := c.Get("userID"); exists {
		uid := v.(uint)
		operaResponse.Liked = services.IsLiked(uid, operaResponse.OperaID)
		operaResponse.Favorited = services.IsFavorited(uid, operaResponse.OperaID)
	}
	resp.Success(c, operaResponse)
}
