package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/models"
	"github.com/lm379/hmx/services"
	"github.com/lm379/hmx/utils"
	"gorm.io/gorm"
)

// HandleCreateOpera (POST /api/v1/operas)
func HandleCreateOpera(c *gin.Context) {
	var input models.CreateOperaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 从 JWT 中间件获取 userID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "Unauthorized")
		return
	}

	// 调用 service 创建 Opera
	response, status := services.CreateOpera(input, userID.(uint))
	if status != http.StatusCreated {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Created(c, response)
}

// HandleGetOperas (GET /api/v1/operas)
func HandleGetOperas(c *gin.Context) {
	pagination := utils.GetPagination(c)

	operas, total, err := services.GetAllOperas(pagination)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch operas")
		return
	}

	// 转换为响应格式，自动处理 URL
	operaResponses := utils.ToOperaResponseList(operas)

	utils.Success(c, gin.H{
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
		utils.BadRequest(c, "Invalid opera ID")
		return
	}

	opera, err := services.GetOperaByID(uint(operaID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Opera not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 转换为响应格式，自动处理 URL
	operaResponse := utils.ToOperaResponse(opera)
	utils.Success(c, operaResponse)
}
