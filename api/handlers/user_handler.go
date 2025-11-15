package handlers

import (
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
