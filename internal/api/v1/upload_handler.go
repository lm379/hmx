package v1

import (
	"strconv"

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

	// 验证 ContentType
	if input.UploadType == "avatars" && len(input.ContentType) < 6 || input.UploadType == "avatars" && input.ContentType[:6] != "image/" {
		resp.BadRequest(c, "Invalid content_type for avatar. Must be an image.")
		return
	}
	if input.UploadType == "videos" && len(input.ContentType) < 6 || input.UploadType == "videos" && input.ContentType[:6] != "video/" {
		resp.BadRequest(c, "Invalid content_type for video. Must be a video.")
		return
	}

	// 获取当前用户ID
	var userID uint
	if id, exists := c.Get("userID"); exists {
		userID = id.(uint)
	}

	// 构造存储路径 (硬编码)
	var basePath string
	switch input.UploadType {
	case "avatars":
		if userID == 0 {
			resp.Unauthorized(c, "Login required for avatar upload")
			return
		}
		basePath = "user/avatar/" + strconv.FormatUint(uint64(userID), 10)
	case "videos":
		basePath = "tmp"
	default:
		basePath = input.UploadType
	}

	// 从 service 获取预签名 URL (pass basePath as the uploadType/folder)
	uploadURL, objectKey, err := services.GeneratePresignedUploadURL(c, basePath, input.Filename, input.ContentType)
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
