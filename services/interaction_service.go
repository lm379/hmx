package services

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/models"
	"gorm.io/gorm"
)

// checkOperaExists 检查作品是否存在
func checkOperaExists(db *gorm.DB, operaID uint) error {
	var opera models.Opera
	if err := db.First(&opera, operaID).Error; err != nil {
		return err
	}
	return nil
}

// ToggleLike 切换点赞状态
func ToggleLike(userID, operaID uint) (gin.H, int) {
	db := database.DB

	// 检查作品是否存在
	if err := checkOperaExists(db, operaID); err != nil {
		return gin.H{"error": "Opera not found"}, http.StatusNotFound
	}

	like := models.Like{UserID: userID, OperaID: operaID}

	// 尝试查找点赞
	err := db.Where("user_id = ? AND opera_id = ?", userID, operaID).First(&like).Error
	if err == nil {
		// 找到了 -> 取消点赞
		if err := db.Delete(&like).Error; err != nil {
			return gin.H{"error": "Failed to unlike"}, http.StatusInternalServerError
		}
		return gin.H{"status": "unliked"}, http.StatusOK
	} else if err == gorm.ErrRecordNotFound {
		// 没找到 -> 创建点赞
		if err := db.Create(&like).Error; err != nil {
			return gin.H{"error": "Failed to like"}, http.StatusInternalServerError
		}
		return gin.H{"status": "liked"}, http.StatusCreated
	}

	return gin.H{"error": "Database error"}, http.StatusInternalServerError
}

// ToggleFavorite 切换收藏状态
func ToggleFavorite(userID, operaID uint) (gin.H, int) {
	db := database.DB

	// 检查作品
	if err := checkOperaExists(db, operaID); err != nil {
		return gin.H{"error": "Opera not found"}, http.StatusNotFound
	}

	fav := models.Favorite{UserID: userID, OperaID: operaID}

	// 切换逻辑
	err := db.Where("user_id = ? AND opera_id = ?", userID, operaID).First(&fav).Error
	if err == nil {
		if err := db.Delete(&fav).Error; err != nil {
			return gin.H{"error": "Failed to unfavorite"}, http.StatusInternalServerError
		}
		return gin.H{"status": "unfavorited"}, http.StatusOK
	} else if err == gorm.ErrRecordNotFound {
		if err := db.Create(&fav).Error; err != nil {
			return gin.H{"error": "Failed to favorite"}, http.StatusInternalServerError
		}
		return gin.H{"status": "favorited"}, http.StatusCreated
	}

	return gin.H{"error": "Database error"}, http.StatusInternalServerError
}

// CreateComment 创建评论
func CreateComment(userID, operaID uint, input models.CreateCommentRequest) (gin.H, int) {
	db := database.DB

	// 检查作品
	if err := checkOperaExists(db, operaID); err != nil {
		return gin.H{"error": "Opera not found"}, http.StatusNotFound
	}

	// 检查父评论是否存在
	if input.ParentCommentID != nil {
		var parentComment models.Comment
		// 确保父评论也属于同一个 Opera
		err := db.First(&parentComment, "comment_id = ? AND opera_id = ?", *input.ParentCommentID, operaID).Error
		if err != nil {
			return gin.H{"error": "Parent comment not found on this opera"}, http.StatusBadRequest
		}
	}

	// 创建评论
	comment := models.Comment{
		UserID:          sql.NullInt64{Int64: int64(userID), Valid: true},
		OperaID:         operaID,
		ParentCommentID: input.ParentCommentID,
		CommentText:     input.CommentText,
	}

	if err := db.Create(&comment).Error; err != nil {
		return gin.H{"error": "Failed to create comment"}, http.StatusInternalServerError
	}

	// (在真实项目中，您可能希望将创建的评论对象返回)
	return gin.H{"message": "Comment created successfully", "comment_id": comment.CommentID}, http.StatusCreated
}

// DeleteComment 删除评论
func DeleteComment(userID, commentID uint, userRole string) (gin.H, int) {
	db := database.DB
	var comment models.Comment

	if err := db.First(&comment, commentID).Error; err != nil {
		return gin.H{"error": "Comment not found"}, http.StatusNotFound
	}

	// 权限检查：必须是评论所有者 或 管理员
	if comment.UserID.Int64 != int64(userID) && userRole != string(models.Administrator) {
		return gin.H{"error": "Forbidden: You cannot delete this comment"}, http.StatusForbidden
	}

	if err := db.Delete(&comment).Error; err != nil {
		return gin.H{"error": "Failed to delete comment"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Comment deleted successfully"}, http.StatusOK
}
