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
		resp.InternalServerError(c, "Database error")
		return
	}

	// 转换为响应格式，自动处理 URL
	operaResponse := converter.ToOperaResponse(opera)
	resp.Success(c, operaResponse)
}
