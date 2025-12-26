package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/models"
	"github.com/lm379/hmx/services"
	"github.com/lm379/hmx/utils"
)

// HandleGetMe (GET /api/v1/users/me)
func HandleGetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "Unauthorized")
		return
	}

	var user models.Users
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	// 转换为响应格式，自动处理 URL
	userResponse := utils.ToUserResponse(&user)
	utils.Success(c, userResponse)
}

// HandleGetUserLikes (GET /api/v1/users/me/likes)
func HandleGetUserLikes(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pagination := utils.GetPagination(c)

	operas, total, err := services.GetUserLikes(userID, pagination)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch liked operas")
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

// HandleGetUserFavorites (GET /api/v1/users/me/favorites)
func HandleGetUserFavorites(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pagination := utils.GetPagination(c)

	operas, total, err := services.GetUserFavorites(userID, pagination)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch favorited operas")
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

// HandleGetUserHistory (GET /api/v1/users/me/history)
func HandleGetUserHistory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pagination := utils.GetPagination(c)

	operas, total, err := services.GetUserHistory(userID, pagination)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch watch history")
		return
	}

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

// HandleRecordHistory (POST /api/v1/operas/:id/history)
func HandleRecordHistory(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid opera ID")
		return
	}

	var userIDPtr *uint
	if userID, exists := c.Get("userID"); exists {
		id := userID.(uint)
		userIDPtr = &id
	}

	if err := services.RecordPlayHistory(userIDPtr, uint(operaID)); err != nil {
		utils.InternalServerError(c, "Failed to record history")
		return
	}

	utils.Success(c, gin.H{"message": "History recorded"})
}
