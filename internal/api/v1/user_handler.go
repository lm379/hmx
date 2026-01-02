package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/converter"
	"github.com/lm379/hmx/pkg/pagination"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleGetMe (GET /api/v1/users/me)
func HandleGetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		resp.Unauthorized(c, "Unauthorized")
		return
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		resp.NotFound(c, "User not found")
		return
	}

	// 转换为响应格式，自动处理 URL
	userResponse := converter.ToUserResponse(user)
	resp.Success(c, userResponse)
}

// HandleGetUserLikes (GET /api/v1/users/me/likes)
func HandleGetUserLikes(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pagination := pagination.GetPagination(c)

	operas, total, err := services.GetUserLikes(userID, pagination)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch liked operas")
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

// HandleGetUserFavorites (GET /api/v1/users/me/favorites)
func HandleGetUserFavorites(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pagination := pagination.GetPagination(c)

	operas, total, err := services.GetUserFavorites(userID, pagination)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch favorited operas")
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

// HandleGetUserHistory (GET /api/v1/users/me/history)
func HandleGetUserHistory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	pagination := pagination.GetPagination(c)

	operas, total, err := services.GetUserHistory(userID, pagination)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch watch history")
		return
	}

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

// HandleRecordHistory (POST /api/v1/operas/:id/history)
func HandleRecordHistory(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	var userIDPtr *uint
	if userID, exists := c.Get("userID"); exists {
		id := userID.(uint)
		userIDPtr = &id
	}

	if err := services.RecordPlayHistory(userIDPtr, uint(operaID)); err != nil {
		resp.InternalServerError(c, "Failed to record history")
		return
	}

	resp.Success(c, gin.H{"message": "History recorded"})
}

// HandleUpdateUserProfile (PUT /api/v1/users/me)
func HandleUpdateUserProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var input models.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	if err := services.UpdateUserProfile(userID, input); err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			resp.Error(c, svcErr.Code, svcErr.Message)
		} else {
			resp.InternalServerError(c, err.Error())
		}
		return
	}

	resp.Success(c, gin.H{"message": "Profile updated successfully"})
}

// HandleUpdatePassword (POST /api/v1/users/me/password)
func HandleUpdatePassword(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var input models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	if err := services.UpdatePassword(userID, input); err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			resp.Error(c, svcErr.Code, svcErr.Message)
		} else {
			resp.InternalServerError(c, err.Error())
		}
		return
	}

	resp.Success(c, gin.H{"message": "Password updated successfully"})
}

// HandleSendCodeToCurrentUser (POST /api/v1/users/me/email-code)
func HandleSendCodeToCurrentUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	// Fetch current user email
	user, err := services.GetUserByID(userID)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch user info")
		return
	}

	if !user.Email.Valid || user.Email.String == "" {
		resp.BadRequest(c, "User has no email bound")
		return
	}

	response, status := services.SendVerificationCode(user.Email.String)
	if status != 200 {
		resp.Error(c, status, response["error"].(string))
		return
	}

	resp.Success(c, response)
}
