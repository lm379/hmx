package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/models"
	"github.com/lm379/hmx/services"
	"github.com/lm379/hmx/utils"
)

// getOperaIDParam 从 URL param 中解析 opera ID
func getOperaIDParam(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(operaID), nil
}

// HandleToggleLike (POST /api/v1/operas/:id/like)
func HandleToggleLike(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	operaID, err := getOperaIDParam(c)
	if err != nil {
		utils.BadRequest(c, "Invalid opera ID")
		return
	}

	response, status := services.ToggleLike(userID, operaID)
	if status >= 400 {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Success(c, response)
}

// HandleToggleFavorite (POST /api/v1/operas/:id/favorite)
func HandleToggleFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	operaID, err := getOperaIDParam(c)
	if err != nil {
		utils.BadRequest(c, "Invalid opera ID")
		return
	}

	response, status := services.ToggleFavorite(userID, operaID)
	if status >= 400 {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Success(c, response)
}

// HandleCreateComment (POST /api/v1/operas/:id/comments)
func HandleCreateComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	operaID, err := getOperaIDParam(c)
	if err != nil {
		utils.BadRequest(c, "Invalid opera ID")
		return
	}

	var input models.CreateCommentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, status := services.CreateComment(userID, operaID, input)
	if status >= 400 {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Created(c, response)
}

// HandleDeleteComment (DELETE /api/v1/comments/:id)
func HandleDeleteComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("role").(string)

	idParam := c.Param("id")
	commentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid comment ID")
		return
	}

	response, status := services.DeleteComment(userID, uint(commentID), userRole)
	if status >= 400 {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Success(c, response)
}
