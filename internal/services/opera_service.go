package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/pagination"
)

var operaRepo = repository.NewOperaRepo()

// Helper to process new artist names
func processNewArtistNames(newNames []string) ([]uint, error) {
	var ids []uint
	for _, name := range newNames {
		if name == "" {
			continue
		}
		// Check if exists
		artist, err := artistRepo.FindByName(name)
		if err != nil {
			// Create new
			newArtist := models.Artist{
				Name: name,
			}
			if err := artistRepo.Create(&newArtist); err != nil {
				return nil, err
			}
			ids = append(ids, newArtist.ArtistID)
		} else {
			ids = append(ids, artist.ArtistID)
		}
	}
	return ids, nil
}

// CreateOpera 在数据库中创建新的 Opera 记录
func CreateOpera(input models.CreateOperaRequest, userID uint) (gin.H, int) {
	// 合并已有ID和新创建的ID
	allArtistIDs := input.ArtistIDs
	if len(input.NewArtistNames) > 0 {
		newIDs, err := processNewArtistNames(input.NewArtistNames)
		if err != nil {
			return gin.H{"error": "Failed to process new artist names"}, http.StatusInternalServerError
		}
		allArtistIDs = append(allArtistIDs, newIDs...)
	}

	// 查找艺术家
	var artists []*models.Artist
	if len(allArtistIDs) > 0 {
		var err error
		artists, err = operaRepo.FindArtistsByIDs(allArtistIDs)
		if err != nil {
			return gin.H{"error": "Invalid artist IDs provided"}, http.StatusBadRequest
		}
	}

	// 创建 Opera 实例
	newOpera := models.Opera{
		OperaTitle:  input.Title,
		Description: input.Description,
		VideoPath:   input.VideoPath,
		Avatar:      sql.NullString{String: input.AvatarPath, Valid: input.AvatarPath != ""},
		Artists:     artists,
	}

	// 创建 Opera 记录
	if err := operaRepo.Create(&newOpera); err != nil {
		return gin.H{"error": "Failed to create opera record: " + err.Error()}, http.StatusInternalServerError
	}

	return gin.H{"message": "Opera created successfully", "opera_id": newOpera.OperaID}, http.StatusCreated
}

// UpdateOpera 更新 Opera 记录
func UpdateOpera(operaID uint, input models.UpdateOperaRequest) (gin.H, int) {
	// 检查作品是否存在
	opera, err := operaRepo.GetByID(operaID)
	if err != nil {
		return gin.H{"error": "Opera not found"}, http.StatusNotFound
	}

	updates := make(map[string]interface{})
	if input.Title != "" {
		updates["opera_title"] = input.Title
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.VideoPath != "" {
		updates["video_path"] = input.VideoPath
	}
	if input.AvatarPath != "" {
		updates["avatar"] = input.AvatarPath
	}
	if input.IsHidden != nil {
		updates["is_hidden"] = *input.IsHidden
	}

	// 更新基本字段
	if len(updates) > 0 {
		if err := operaRepo.Update(operaID, updates); err != nil {
			return gin.H{"error": "Failed to update opera"}, http.StatusInternalServerError
		}
	}

	// 更新艺术家关联
	if input.ArtistIDs != nil || input.NewArtistNames != nil {
		allArtistIDs := input.ArtistIDs
		if len(input.NewArtistNames) > 0 {
			newIDs, err := processNewArtistNames(input.NewArtistNames)
			if err != nil {
				return gin.H{"error": "Failed to process new artist names"}, http.StatusInternalServerError
			}
			allArtistIDs = append(allArtistIDs, newIDs...)
		}

		var artists []*models.Artist
		if len(allArtistIDs) > 0 {
			artists, err = operaRepo.FindArtistsByIDs(allArtistIDs)
			if err != nil {
				return gin.H{"error": "Invalid artist IDs"}, http.StatusBadRequest
			}
		} else {
			// 如果两个列表都为空（但 input 不为 nil），说明是要清空关联
			// 这里逻辑有点微妙，如果 ArtistIDs 是空数组，通常意味着清空
			artists = []*models.Artist{}
		}
		if err := operaRepo.UpdateArtists(opera, artists); err != nil {
			return gin.H{"error": "Failed to update artists"}, http.StatusInternalServerError
		}
	}

	return gin.H{"message": "Opera updated successfully"}, http.StatusOK
}

// GetOperas 获取作品列表 (带分页和过滤)
func GetOperas(pagination *pagination.Pagination, includeHidden bool) ([]models.Opera, int64, error) {
	return operaRepo.GetAll(pagination, includeHidden)
}

// GetOperaByID 获取单个作品
func GetOperaByID(operaID uint) (models.Opera, error) {
	opera, err := operaRepo.GetByID(operaID)
	if err != nil {
		return models.Opera{}, err
	}
	return *opera, nil
}

// DeleteOpera 删除作品 (管理员)
func DeleteOpera(operaID uint) (gin.H, int) {
	// 检查作品是否存在
	exists, err := operaRepo.Exists(operaID)
	if err != nil {
		return gin.H{"error": "Failed to check opera"}, http.StatusInternalServerError
	}
	if !exists {
		return gin.H{"error": "Opera not found"}, http.StatusNotFound
	}

	if err := operaRepo.Delete(operaID); err != nil {
		return gin.H{"error": "Failed to delete opera"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Opera deleted successfully"}, http.StatusOK
}

// HandleTranscodeCallback 处理腾讯云转码完成回调
func HandleTranscodeCallback(callback models.TencentCloudCallbackRequest) (gin.H, int) {
	// 检查事件类型
	if callback.EventName != "TaskFinish" {
		return gin.H{"message": "Event ignored", "event": callback.EventName}, http.StatusOK
	}

	// 检查是否有任务详情
	if len(callback.JobsDetail) == 0 {
		return gin.H{"error": "No job details in callback"}, http.StatusBadRequest
	}

	// 获取第一个任务详情
	job := callback.JobsDetail[0]

	// 检查任务状态
	if job.State != "Success" {
		return gin.H{"error": "Job not successful", "state": job.State, "message": job.Message}, http.StatusBadRequest
	}

	// 获取原始文件路径（上传时在 tmp 目录）
	originalPath := job.Input.Object

	// 从路径中提取视频标题（文件名，不含扩展名）
	videoTitle := extractTitleFromPath(originalPath)
	if videoTitle == "" {
		return gin.H{"error": "Failed to extract title from path: " + originalPath}, http.StatusBadRequest
	}

	// 查找对应的 Opera 记录（通过 opera_title）
	opera, err := operaRepo.FindByTitle(videoTitle)
	if err != nil {
		return gin.H{"error": "Opera not found with title: " + videoTitle}, http.StatusNotFound
	}

	// 准备更新字段
	updates := make(map[string]interface{})

	// 根据任务类型处理不同的输出
	switch job.Tag {
	case "Transcode":
		// 转码任务 - 更新视频路径和时长
		if job.Operation.MediaResult != nil && len(job.Operation.MediaResult.OutputFile.ObjectName) > 0 {
			// 更新视频路径（HLS m3u8）
			hlsPath := job.Operation.MediaResult.OutputFile.ObjectName[0]
			updates["video_path"] = hlsPath

			// 提取并转换时长
			if job.Operation.MediaInfo != nil && job.Operation.MediaInfo.Format.Duration != "" {
				duration, err := parseDurationToTime(job.Operation.MediaInfo.Format.Duration)
				if err == nil {
					updates["duration"] = duration
				}
			}
		}

	case "SmartCover":
		// 智能封面任务 - 更新封面路径
		if job.Operation.MediaResult != nil && len(job.Operation.MediaResult.OutputFile.ObjectName) > 0 {
			coverPath := job.Operation.MediaResult.OutputFile.ObjectName[0]
			updates["avatar"] = coverPath
		}

	case "SpeechRecognition":
		// 语音识别任务 - 更新字幕路径
		if job.Operation.SpeechRecognitionResult != nil && job.Operation.SpeechRecognitionResult.ObjectName != "" {
			srtPath := job.Operation.SpeechRecognitionResult.ObjectName
			updates["srt_path"] = srtPath
		}
	}

	// 执行更新
	if len(updates) > 0 {
		if err := operaRepo.Update(opera.OperaID, updates); err != nil {
			return gin.H{"error": "Failed to update opera after callback"}, http.StatusInternalServerError
		}
	}

	return gin.H{
		"message":  "Callback processed successfully",
		"opera_id": opera.OperaID,
		"job_type": job.Tag,
		"updates":  updates,
	}, http.StatusOK
}

// extractTitleFromPath 从文件路径中提取标题（文件名，不含扩展名）
// 例如: "tmp/xxx.mp4" -> "xxx"
func extractTitleFromPath(path string) string {
	// 获取文件名（去掉目录部分）
	filename := filepath.Base(path)
	// 去掉扩展名
	ext := filepath.Ext(filename)
	title := strings.TrimSuffix(filename, ext)
	return title
}

// parseDurationToTime 将秒数字符串转换为 HH:MM:SS 格式
func parseDurationToTime(durationStr string) (string, error) {
	var seconds float64
	_, err := fmt.Sscanf(durationStr, "%f", &seconds)
	if err != nil {
		return "", err
	}

	// 转换为时:分:秒
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs), nil
}
