package services

import (
	"database/sql"
	"net/http"

	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/urlutils"
	"gorm.io/gorm"
)

var interactionRepo = repository.NewInteractionRepo()

// InteractionResponse 交互操作的响应结构
type InteractionResponse struct {
	Status        string `json:"status,omitempty"`
	Liked         bool   `json:"liked"`
	Favorited     bool   `json:"favorited"`
	LikeCount     int64  `json:"like_count"`
	FavoriteCount int64  `json:"favorite_count"`
	ShareCount    int64  `json:"share_count"`
	PlayCount     int64  `json:"play_count"`
}

// ToggleLike 切换点赞状态
func ToggleLike(userID, operaID uint) (*InteractionResponse, error) {
	// 检查作品是否存在
	if err := interactionRepo.CheckOperaExists(operaID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &ServiceError{Code: http.StatusNotFound, Message: "Opera not found"}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to check opera"}
	}

	var status string
	var liked bool
	// 尝试查找点赞
	_, err := interactionRepo.GetLike(userID, operaID)
	if err == nil {
		// 找到了 -> 取消点赞
		if err := interactionRepo.DeleteLike(userID, operaID); err != nil {
			return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to unlike"}
		}
		status = "unliked"
		liked = false
	} else if err == gorm.ErrRecordNotFound {
		// 没找到 -> 创建点赞
		if err := interactionRepo.CreateLike(userID, operaID); err != nil {
			return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to like"}
		}
		status = "liked"
		liked = true
	} else {
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Database error"}
	}

	lc, fc, sc, pc := GetCounts(operaID)
	return &InteractionResponse{
		Status:        status,
		Liked:         liked,
		Favorited:     interactionRepo.IsFavorited(userID, operaID),
		LikeCount:     lc,
		FavoriteCount: fc,
		ShareCount:    sc,
		PlayCount:     pc,
	}, nil
}

// ToggleFavorite 切换收藏状态
func ToggleFavorite(userID, operaID uint) (*InteractionResponse, error) {
	// 检查作品
	if err := interactionRepo.CheckOperaExists(operaID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &ServiceError{Code: http.StatusNotFound, Message: "Opera not found"}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to check opera"}
	}

	var status string
	var favorited bool
	// 切换逻辑
	_, err := interactionRepo.GetFavorite(userID, operaID)
	if err == nil {
		if err := interactionRepo.DeleteFavorite(userID, operaID); err != nil {
			return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to unfavorite"}
		}
		status = "unfavorited"
		favorited = false
	} else if err == gorm.ErrRecordNotFound {
		if err := interactionRepo.CreateFavorite(userID, operaID); err != nil {
			return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to favorite"}
		}
		status = "favorited"
		favorited = true
	} else {
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Database error"}
	}

	lc, fc, sc, pc := GetCounts(operaID)
	return &InteractionResponse{
		Status:        status,
		Liked:         interactionRepo.IsLiked(userID, operaID),
		Favorited:     favorited,
		LikeCount:     lc,
		FavoriteCount: fc,
		ShareCount:    sc,
		PlayCount:     pc,
	}, nil
}

// GetCounts 返回点赞、收藏、分享、播放计数
func GetCounts(operaID uint) (likeCount int64, favoriteCount int64, shareCount int64, playCount int64) {
	return interactionRepo.CountLikes(operaID),
		interactionRepo.CountFavorites(operaID),
		interactionRepo.CountShares(operaID),
		interactionRepo.CountPlays(operaID)
}

// BatchGetCounts 批量获取多个作品的计数
func BatchGetCounts(operaIDs []uint) (likes, favorites, shares, plays map[uint]int64) {
	// 使用单次查询获取所有计数
	return interactionRepo.BatchCountAll(operaIDs)
}

// IsLiked 检查用户是否已点赞（导出供handler使用）
func IsLiked(userID, operaID uint) bool {
	return interactionRepo.IsLiked(userID, operaID)
}

// IsFavorited 检查用户是否已收藏（导出供handler使用）
func IsFavorited(userID, operaID uint) bool {
	return interactionRepo.IsFavorited(userID, operaID)
}

// ShareResponse 分享操作的响应结构
type ShareResponse struct {
	Message       string `json:"message"`
	LikeCount     int64  `json:"like_count"`
	FavoriteCount int64  `json:"favorite_count"`
	ShareCount    int64  `json:"share_count"`
	PlayCount     int64  `json:"play_count"`
	NewShare      bool   `json:"new_share"` // 是否是新分享
}

// RecordShare 记录一次分享行为（游客或用户）
// 同一用户对同一作品只记录一次分享
func RecordShare(userID *uint, operaID uint) (*ShareResponse, error) {
	// 检查作品是否存在
	if err := interactionRepo.CheckOperaExists(operaID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &ServiceError{Code: http.StatusNotFound, Message: "Opera not found"}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to check opera"}
	}

	// 使用 FirstOrCreate 避免重复记录，同一用户对同一作品只算一次
	isNew, err := interactionRepo.FirstOrCreateShare(userID, operaID)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to record share"}
	}

	lc, fc, sc, pc := GetCounts(operaID)
	return &ShareResponse{
		Message:       "Share recorded",
		LikeCount:     lc,
		FavoriteCount: fc,
		ShareCount:    sc,
		PlayCount:     pc,
		NewShare:      isNew,
	}, nil
}

// CommentResponse 评论操作的响应结构
type CommentResponse struct {
	Message   string `json:"message"`
	CommentID uint   `json:"comment_id"`
}

// CreateComment 创建评论
func CreateComment(userID, operaID uint, input models.CreateCommentRequest) (*CommentResponse, error) {
	// 检查作品
	if err := interactionRepo.CheckOperaExists(operaID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &ServiceError{Code: http.StatusNotFound, Message: "Opera not found"}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to check opera"}
	}

	// 检查父评论是否存在
	if input.ParentCommentID != nil {
		parentComment, err := interactionRepo.GetComment(*input.ParentCommentID)
		// 确保父评论也属于同一个 Opera
		if err != nil || parentComment.OperaID != operaID {
			return nil, &ServiceError{Code: http.StatusBadRequest, Message: "Parent comment not found on this opera"}
		}
	}

	// 创建评论
	comment := &models.Comment{
		UserID:          sql.NullInt64{Int64: int64(userID), Valid: true},
		OperaID:         operaID,
		ParentCommentID: input.ParentCommentID,
		CommentText:     input.CommentText,
	}

	if err := interactionRepo.CreateComment(comment); err != nil {
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to create comment"}
	}

	return &CommentResponse{
		Message:   "Comment created successfully",
		CommentID: comment.CommentID,
	}, nil
}

// ServiceError 服务层错误类型
type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// DeleteCommentResponse 删除评论的响应结构
type DeleteCommentResponse struct {
	Message string `json:"message"`
}

// DeleteComment 删除评论
func DeleteComment(userID, commentID uint, userRole string) (*DeleteCommentResponse, error) {
	comment, err := interactionRepo.GetComment(commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &ServiceError{Code: http.StatusNotFound, Message: "Comment not found"}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to get comment"}
	}

	// 权限检查：必须是评论所有者 或 管理员
	if comment.UserID.Int64 != int64(userID) && userRole != string(models.Administrator) {
		return nil, &ServiceError{Code: http.StatusForbidden, Message: "Forbidden: You cannot delete this comment"}
	}

	if err := interactionRepo.DeleteComment(commentID); err != nil {
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to delete comment"}
	}

	return &DeleteCommentResponse{Message: "Comment deleted successfully"}, nil
}

// CommentLikeResponse 评论点赞响应
type CommentLikeResponse struct {
	Liked     bool  `json:"liked"`
	LikeCount int64 `json:"like_count"`
}

// ToggleCommentLike 切换评论点赞状态
func ToggleCommentLike(userID, commentID uint) (*CommentLikeResponse, error) {
	// Check if comment exists
	if _, err := interactionRepo.GetComment(commentID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &ServiceError{Code: http.StatusNotFound, Message: "Comment not found"}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to check comment"}
	}

	var liked bool
	// Try to get existing like
	_, err := interactionRepo.GetCommentLike(userID, commentID)
	if err == nil {
		// Found -> Unlike
		if err := interactionRepo.DeleteCommentLike(userID, commentID); err != nil {
			return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to unlike comment"}
		}
		liked = false
	} else if err == gorm.ErrRecordNotFound {
		// Not found -> Like
		if err := interactionRepo.CreateCommentLike(userID, commentID); err != nil {
			return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to like comment"}
		}
		liked = true
	} else {
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Database error"}
	}

	count := interactionRepo.CountCommentLikes(commentID)
	return &CommentLikeResponse{
		Liked:     liked,
		LikeCount: count,
	}, nil
}

// GetComments 获取评论列表
func GetComments(operaID uint, userID *uint) ([]models.CommentDTO, error) {
	comments, err := interactionRepo.GetCommentsByOperaID(operaID, userID)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusInternalServerError, Message: "Failed to fetch comments"}
	}

	// 处理头像 URL
	for i := range comments {
		if comments[i].UserIcon != nil && *comments[i].UserIcon != "" {
			fullURL := urlutils.GetFullURL(*comments[i].UserIcon)
			comments[i].UserIcon = &fullURL
		}
	}

	return comments, nil
}
