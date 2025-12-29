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

	result, err := services.ToggleLike(userID, operaID)
	if err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			resp.Error(c, svcErr.Code, svcErr.Message)
		} else {
			resp.InternalServerError(c, err.Error())
		}
		return
	}
	resp.Success(c, result)
}

// HandleToggleFavorite (POST /api/v1/operas/:id/favorite)
func HandleToggleFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	operaID, err := getOperaIDParam(c)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	result, err := services.ToggleFavorite(userID, operaID)
	if err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			resp.Error(c, svcErr.Code, svcErr.Message)
		} else {
			resp.InternalServerError(c, err.Error())
		}
		return
	}
	resp.Success(c, result)
}

// HandleShare (POST /api/v1/operas/:id/share)
func HandleShare(c *gin.Context) {
	var uidPtr *uint
	if v, exists := c.Get("userID"); exists {
		u := v.(uint)
		uidPtr = &u
	}

	operaID, err := getOperaIDParam(c)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	result, err := services.RecordShare(uidPtr, operaID)
	if err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			resp.Error(c, svcErr.Code, svcErr.Message)
		} else {
			resp.InternalServerError(c, err.Error())
		}
		return
	}
	// 如果是新分享，返回201；否则返回200
	if result.NewShare {
		resp.Created(c, result)
	} else {
		resp.Success(c, result)
	}
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

	result, err := services.CreateComment(userID, operaID, input)
	if err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			resp.Error(c, svcErr.Code, svcErr.Message)
		} else {
			resp.InternalServerError(c, err.Error())
		}
		return
	}
	resp.Created(c, result)
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

	result, err := services.DeleteComment(userID, uint(commentID), userRole)
	if err != nil {
		if svcErr, ok := err.(*services.ServiceError); ok {
			resp.Error(c, svcErr.Code, svcErr.Message)
		} else {
			resp.InternalServerError(c, err.Error())
		}
		return
	}
	resp.Success(c, result)
}
