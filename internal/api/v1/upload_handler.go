package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleRequestUploadURL (POST /api/v1/uploads/presign)
func HandleRequestUploadURL(c *gin.Context) {
	var input models.PresignRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	// 验证 uploadType
	if input.UploadType != "videos" && input.UploadType != "avatars" {
		resp.BadRequest(c, "Invalid upload_type. Must be 'videos' or 'avatars'")
		return
	}

	// 从 service 获取预签名 URL
	uploadURL, objectKey, err := services.GeneratePresignedUploadURL(c, input.UploadType, input.Filename, input.ContentType)
	if err != nil {
		resp.InternalServerError(c, "Failed to generate presigned URL")
		return
	}

	// 返回给前端
	resp.Success(c, models.PresignResponse{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
	})
}
