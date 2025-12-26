package services

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/models"
	"github.com/lm379/hmx/utils"
)

// CreateOpera 在数据库中创建新的 Opera 记录
func CreateOpera(input models.CreateOperaRequest, userID uint) (gin.H, int) {
	db := database.DB

	// 查找艺术家
	var artists []*models.Artist
	if len(input.ArtistIDs) > 0 {
		if err := db.Find(&artists, input.ArtistIDs).Error; err != nil {
			return gin.H{"error": "Invalid artist IDs provided"}, http.StatusBadRequest
		}
	}

	// 创建 Opera 实例
	newOpera := models.Opera{
		OperaTitle:  input.Title,
		Description: input.Description,
		VideoPath:   input.VideoPath,
		Avatar:      sql.NullString{String: input.AvatarPath, Valid: input.AvatarPath != ""},
		Artists:     artists, // 关联艺术家
	}

	// 创建 Opera 记录
	err := db.Create(&newOpera).Error
	if err != nil {
		return gin.H{"error": "Failed to create opera record: " + err.Error()}, http.StatusInternalServerError
	}

	return gin.H{"message": "Opera created successfully", "opera_id": newOpera.OperaID}, http.StatusCreated
}

// GetAllOperas 获取所有作品 (带分页)
func GetAllOperas(pagination *utils.Pagination) ([]models.Opera, int64, error) {
	var operas []models.Opera
	var total int64
	db := database.DB.Model(&models.Opera{})

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用分页查询并加载艺术家数据
	err := db.Preload("Artists").Scopes(pagination.Paginate()).Find(&operas).Error

	return operas, total, err
}

// GetOperaByID 获取单个作品
func GetOperaByID(operaID uint) (models.Opera, error) {
	var opera models.Opera
	db := database.DB

	err := db.Preload("Artists").Where("opera_id = ?", operaID).
		First(&opera).Error

	return opera, err
}

// DeleteOpera 删除作品 (管理员)
func DeleteOpera(operaID uint) (gin.H, int) {
	var opera models.Opera
	db := database.DB

	if err := db.First(&opera, operaID).Error; err != nil {
		return gin.H{"error": "Opera not found"}, http.StatusNotFound
	}

	if err := db.Delete(&opera).Error; err != nil {
		return gin.H{"error": "Failed to delete opera"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Opera deleted successfully"}, http.StatusOK
}
