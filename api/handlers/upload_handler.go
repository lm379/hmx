package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/models"
	"github.com/lm379/hmx/services"
	"github.com/lm379/hmx/utils"
)

// HandleRequestUploadURL (POST /api/v1/uploads/presign)
func HandleRequestUploadURL(c *gin.Context) {
	var input models.PresignRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 验证 uploadType
	if input.UploadType != "videos" && input.UploadType != "avatars" {
		utils.BadRequest(c, "Invalid upload_type. Must be 'videos' or 'avatars'")
		return
	}

	// 从 service 获取预签名 URL
	uploadURL, objectKey, err := services.GeneratePresignedUploadURL(c, input.UploadType, input.Filename, input.ContentType)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate presigned URL")
		return
	}

	// 返回给前端
	utils.Success(c, models.PresignResponse{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
	})
}
