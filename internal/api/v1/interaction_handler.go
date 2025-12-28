package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	resp "github.com/lm379/hmx/pkg/response"
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
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	response, status := services.ToggleLike(userID, operaID)
	if status >= 400 {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleToggleFavorite (POST /api/v1/operas/:id/favorite)
func HandleToggleFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	operaID, err := getOperaIDParam(c)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	response, status := services.ToggleFavorite(userID, operaID)
	if status >= 400 {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleCreateComment (POST /api/v1/operas/:id/comments)
func HandleCreateComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	operaID, err := getOperaIDParam(c)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	var input models.CreateCommentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	response, status := services.CreateComment(userID, operaID, input)
	if status >= 400 {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Created(c, response)
}

// HandleDeleteComment (DELETE /api/v1/comments/:id)
func HandleDeleteComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("role").(string)

	idParam := c.Param("id")
	commentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid comment ID")
		return
	}

	response, status := services.DeleteComment(userID, uint(commentID), userRole)
	if status >= 400 {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}
