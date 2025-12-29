package repository

import (
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"gorm.io/gorm"
)

// InteractionRepo 交互相关的数据库操作
type InteractionRepo struct{}

// NewInteractionRepo 创建交互仓库实例
func NewInteractionRepo() *InteractionRepo {
	return &InteractionRepo{}
}

// getDB 获取数据库实例
func (r *InteractionRepo) getDB() *gorm.DB {
	return database.DB
}

// CheckOperaExists 检查作品是否存在
func (r *InteractionRepo) CheckOperaExists(operaID uint) error {
	var opera models.Opera
	return r.getDB().First(&opera, operaID).Error
}

// GetLike 获取点赞记录
func (r *InteractionRepo) GetLike(userID, operaID uint) (*models.Like, error) {
	var like models.Like
	err := r.getDB().Where("user_id = ? AND opera_id = ?", userID, operaID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// CreateLike 创建点赞记录
func (r *InteractionRepo) CreateLike(userID, operaID uint) error {
	like := models.Like{UserID: userID, OperaID: operaID}
	return r.getDB().Create(&like).Error
}

// DeleteLike 删除点赞记录
func (r *InteractionRepo) DeleteLike(userID, operaID uint) error {
	return r.getDB().Where("user_id = ? AND opera_id = ?", userID, operaID).Delete(&models.Like{}).Error
}

// CountLikes 统计作品的点赞数
func (r *InteractionRepo) CountLikes(operaID uint) int64 {
	var count int64
	r.getDB().Model(&models.Like{}).Where("opera_id = ?", operaID).Count(&count)
	return count
}

// BatchCountLikes 批量统计多个作品的点赞数
func (r *InteractionRepo) BatchCountLikes(operaIDs []uint) map[uint]int64 {
	if len(operaIDs) == 0 {
		return make(map[uint]int64)
	}
	type Result struct {
		OperaID uint
		Count   int64
	}
	var results []Result
	r.getDB().Model(&models.Like{}).
		Select("opera_id, COUNT(*) as count").
		Where("opera_id IN ?", operaIDs).
		Group("opera_id").
		Scan(&results)
	
	countMap := make(map[uint]int64)
	for _, result := range results {
		countMap[result.OperaID] = result.Count
	}
	return countMap
}

// IsLiked 检查用户是否已点赞
func (r *InteractionRepo) IsLiked(userID, operaID uint) bool {
	var count int64
	r.getDB().Model(&models.Like{}).Where("user_id = ? AND opera_id = ?", userID, operaID).Count(&count)
	return count > 0
}

// GetFavorite 获取收藏记录
func (r *InteractionRepo) GetFavorite(userID, operaID uint) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.getDB().Where("user_id = ? AND opera_id = ?", userID, operaID).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

// CreateFavorite 创建收藏记录
func (r *InteractionRepo) CreateFavorite(userID, operaID uint) error {
	favorite := models.Favorite{UserID: userID, OperaID: operaID}
	return r.getDB().Create(&favorite).Error
}

// DeleteFavorite 删除收藏记录
func (r *InteractionRepo) DeleteFavorite(userID, operaID uint) error {
	return r.getDB().Where("user_id = ? AND opera_id = ?", userID, operaID).Delete(&models.Favorite{}).Error
}

// CountFavorites 统计作品的收藏数
func (r *InteractionRepo) CountFavorites(operaID uint) int64 {
	var count int64
	r.getDB().Model(&models.Favorite{}).Where("opera_id = ?", operaID).Count(&count)
	return count
}

// BatchCountFavorites 批量统计多个作品的收藏数
func (r *InteractionRepo) BatchCountFavorites(operaIDs []uint) map[uint]int64 {
	if len(operaIDs) == 0 {
		return make(map[uint]int64)
	}
	type Result struct {
		OperaID uint
		Count   int64
	}
	var results []Result
	r.getDB().Model(&models.Favorite{}).
		Select("opera_id, COUNT(*) as count").
		Where("opera_id IN ?", operaIDs).
		Group("opera_id").
		Scan(&results)
	
	countMap := make(map[uint]int64)
	for _, result := range results {
		countMap[result.OperaID] = result.Count
	}
	return countMap
}

// IsFavorited 检查用户是否已收藏
func (r *InteractionRepo) IsFavorited(userID, operaID uint) bool {
	var count int64
	r.getDB().Model(&models.Favorite{}).Where("user_id = ? AND opera_id = ?", userID, operaID).Count(&count)
	return count > 0
}

// FirstOrCreateShare 查找或创建分享记录（同一用户对同一作品只记录一次）
func (r *InteractionRepo) FirstOrCreateShare(userID *uint, operaID uint) (bool, error) {
	share := models.Share{OperaID: operaID, UserID: userID}
	result := r.getDB().Where(models.Share{UserID: userID, OperaID: operaID}).FirstOrCreate(&share)
	if result.Error != nil {
		return false, result.Error
	}
	// 如果RowsAffected > 0，说明是新创建的
	return result.RowsAffected > 0, nil
}

// CountShares 统计作品的分享数（去重后）
func (r *InteractionRepo) CountShares(operaID uint) int64 {
	var count int64
	r.getDB().Model(&models.Share{}).Where("opera_id = ?", operaID).Count(&count)
	return count
}

// BatchCountShares 批量统计多个作品的分享数
func (r *InteractionRepo) BatchCountShares(operaIDs []uint) map[uint]int64 {
	if len(operaIDs) == 0 {
		return make(map[uint]int64)
	}
	type Result struct {
		OperaID uint
		Count   int64
	}
	var results []Result
	r.getDB().Model(&models.Share{}).
		Select("opera_id, COUNT(*) as count").
		Where("opera_id IN ?", operaIDs).
		Group("opera_id").
		Scan(&results)
	
	countMap := make(map[uint]int64)
	for _, result := range results {
		countMap[result.OperaID] = result.Count
	}
	return countMap
}

// CountPlays 统计作品的播放数
func (r *InteractionRepo) CountPlays(operaID uint) int64 {
	var count int64
	r.getDB().Model(&models.PlayHistory{}).
		Select("COALESCE(SUM(count), 0)").
		Where("opera_id = ?", operaID).
		Scan(&count)
	return count
}

// BatchCountAll 一次查询获取所有计数（点赞、收藏、分享、播放）
func (r *InteractionRepo) BatchCountAll(operaIDs []uint) (likes, favorites, shares, plays map[uint]int64) {
	if len(operaIDs) == 0 {
		return make(map[uint]int64), make(map[uint]int64), make(map[uint]int64), make(map[uint]int64)
	}

	type Result struct {
		OperaID uint
		Likes   int64
		Favs    int64
		Shares  int64
		Plays   int64
	}

	var results []Result
	
	// 使用子查询一次性获取所有计数
	r.getDB().Raw(`
		SELECT 
			o.opera_id,
			COALESCE(l.count, 0) as likes,
			COALESCE(f.count, 0) as favs,
			COALESCE(s.count, 0) as shares,
			COALESCE(p.count, 0) as plays
		FROM unnest(?::int[]) AS o(opera_id)
		LEFT JOIN (
			SELECT opera_id, COUNT(*) as count 
			FROM likes 
			WHERE opera_id = ANY(?)
			GROUP BY opera_id
		) l ON o.opera_id = l.opera_id
		LEFT JOIN (
			SELECT opera_id, COUNT(*) as count 
			FROM favorites 
			WHERE opera_id = ANY(?)
			GROUP BY opera_id
		) f ON o.opera_id = f.opera_id
		LEFT JOIN (
			SELECT opera_id, COUNT(*) as count 
			FROM shares 
			WHERE opera_id = ANY(?)
			GROUP BY opera_id
		) s ON o.opera_id = s.opera_id
		LEFT JOIN (
			SELECT opera_id, SUM(count) as count 
			FROM play_history 
			WHERE opera_id = ANY(?)
			GROUP BY opera_id
		) p ON o.opera_id = p.opera_id
	`, operaIDs, operaIDs, operaIDs, operaIDs, operaIDs).Scan(&results)

	likes = make(map[uint]int64)
	favorites = make(map[uint]int64)
	shares = make(map[uint]int64)
	plays = make(map[uint]int64)

	for _, result := range results {
		likes[result.OperaID] = result.Likes
		favorites[result.OperaID] = result.Favs
		shares[result.OperaID] = result.Shares
		plays[result.OperaID] = result.Plays
	}

	return
}

// CreateComment 创建评论
func (r *InteractionRepo) CreateComment(comment *models.Comment) error {
	return r.getDB().Create(comment).Error
}

// GetComment 获取评论
func (r *InteractionRepo) GetComment(commentID uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.getDB().First(&comment, commentID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteComment 删除评论
func (r *InteractionRepo) DeleteComment(commentID uint) error {
	return r.getDB().Delete(&models.Comment{}, commentID).Error
}
