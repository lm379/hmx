package v1

import (
	"errors"

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
	validTypes := map[string]bool{
		"user_avatar":     true,
		"artist_avatar":   true,
		"opera_cover":     true,
		"news_cover":      true,
		"education_pdf":   true,
		"education_cover": true,
		"video_upload":    true,
	}
	if !validTypes[input.UploadType] {
		resp.BadRequest(c, "Invalid upload_type. Must be 'user_avatar', 'artist_avatar', 'opera_cover', 'news_cover', 'education_pdf', 'education_cover', or 'video_upload'")
		return
	}

	// 验证 ContentType
	if input.UploadType == "education_pdf" {
		// 教育资源只允许上传 PDF，避免前台阅读器加载非预期文件。
		if input.ContentType != "application/pdf" {
			resp.BadRequest(c, "Invalid content_type for education_pdf. Must be application/pdf.")
			return
		}
	} else if input.UploadType == "video_upload" {
		// 视频类型必须是video/*
		if len(input.ContentType) < 6 || input.ContentType[:6] != "video/" {
			resp.BadRequest(c, "Invalid content_type for video. Must be a video.")
			return
		}
	} else {
		// 其他类型必须是图片
		if len(input.ContentType) < 6 || input.ContentType[:6] != "image/" {
			resp.BadRequest(c, "Invalid content_type. Must be an image.")
			return
		}
	}

	// 获取当前用户ID和角色（从JWT token）
	var userID uint
	var userRole models.UserRole

	if id, exists := c.Get("userID"); exists {
		userID = id.(uint)
	} else {
		resp.Unauthorized(c, "Login required")
		return
	}

	if role, exists := c.Get("role"); exists {
		userRole = models.UserRole(role.(string))
	} else {
		resp.Unauthorized(c, "Role not found in token")
		return
	}

	// 调用S3 service生成预签名URL
	uploadURL, objectKey, err := services.GeneratePresignedUploadURL(
		c,
		input.UploadType,
		userID,
		userRole,
		input.TargetID,
		input.Filename,
		input.ContentType,
	)
	if err != nil {
		// 使用 errors.Is 判断错误类型
		if errors.Is(err, services.ErrPermissionDenied) {
			resp.Forbidden(c, err.Error())
			return
		}
		if errors.Is(err, services.ErrInvalidUploadType) {
			resp.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, services.ErrOperaNotFound) {
			resp.NotFound(c, err.Error())
			return
		}
		resp.InternalServerError(c, "Failed to generate presigned URL: "+err.Error())
		return
	}

	// 返回给前端
	resp.Success(c, models.PresignResponse{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
	})
}

// HandleTranscodeCallback (POST /api/v1/callback/transcode)
// 处理腾讯云数据万象任务完成回调
func HandleTranscodeCallback(c *gin.Context) {
	var callback models.TencentCloudCallbackRequest
	if err := c.ShouldBindJSON(&callback); err != nil {
		resp.BadRequest(c, "Invalid callback data: "+err.Error())
		return
	}

	// 调用服务层处理
	response, status := services.HandleTranscodeCallback(callback)
	c.JSON(status, response)
}
