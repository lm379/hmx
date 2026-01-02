package services

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/pagination"
)

var operaRepo = repository.NewOperaRepo()

// CreateOpera 在数据库中创建新的 Opera 记录
func CreateOpera(input models.CreateOperaRequest, userID uint) (gin.H, int) {
	// 查找艺术家
	var artists []*models.Artist
	if len(input.ArtistIDs) > 0 {
		var err error
		artists, err = operaRepo.FindArtistsByIDs(input.ArtistIDs)
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
	if input.ArtistIDs != nil {
		var artists []*models.Artist
		if len(input.ArtistIDs) > 0 {
			artists, err = operaRepo.FindArtistsByIDs(input.ArtistIDs)
			if err != nil {
				return gin.H{"error": "Invalid artist IDs"}, http.StatusBadRequest
			}
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
