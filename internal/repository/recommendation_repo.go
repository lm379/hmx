package repository

import (
	"github.com/lm379/hmx/database"
	"gorm.io/gorm"
)

// RecommendationRepo 推荐相关的数据库操作
type RecommendationRepo struct{}

// NewRecommendationRepo 创建推荐仓库实例
func NewRecommendationRepo() *RecommendationRepo {
	return &RecommendationRepo{}
}

// getDB 获取数据库实例
func (r *RecommendationRepo) getDB() *gorm.DB {
	return database.DB
}

// UserInterestScore 用户兴趣分视图模型
type UserInterestScore struct {
	UserID     uint    `gorm:"column:user_id"`
	OperaID    uint    `gorm:"column:opera_id"`
	TotalScore float64 `gorm:"column:total_score"`
}

func (UserInterestScore) TableName() string {
	return "user_interest_scores"
}

// GetOperaInterestScore 获取作品的兴趣分
// 公式：(播放 × 1) + (点赞 × 3) + (收藏 × 5)
func (r *RecommendationRepo) GetOperaInterestScore(operaID uint) (float64, error) {
	var totalScore float64

	err := r.getDB().Model(&UserInterestScore{}).
		Select("COALESCE(SUM(total_score), 0)").
		Where("opera_id = ?", operaID).
		Scan(&totalScore).Error

	return totalScore, err
}

// GetUserPreferenceVector 获取用户的偏好向量（基于用户历史行为）
func (r *RecommendationRepo) GetUserPreferenceVector(userID uint) ([]float64, error) {
	var result struct {
		Embedding []float64
	}

	err := r.getDB().Raw(`
		WITH user_operas AS (
			-- 收藏的作品（权重5）
			SELECT o.embedding, 5 as weight
			FROM opera o
			INNER JOIN favorites f ON o.opera_id = f.opera_id
			WHERE f.user_id = ? AND o.embedding IS NOT NULL
			
			UNION ALL
			
			-- 点赞的作品（权重3）
			SELECT o.embedding, 3 as weight
			FROM opera o
			INNER JOIN likes l ON o.opera_id = l.opera_id
			WHERE l.user_id = ? AND o.embedding IS NOT NULL
			AND o.opera_id NOT IN (SELECT opera_id FROM favorites WHERE user_id = ?)
			
			UNION ALL
			
			-- 播放的作品（权重1，限制最近20个）
			SELECT o.embedding, 1 as weight
			FROM opera o
			INNER JOIN (
				SELECT DISTINCT opera_id 
				FROM play_history 
				WHERE user_id = ? 
				ORDER BY updated_at DESC 
				LIMIT 20
			) ph ON o.opera_id = ph.opera_id
			WHERE o.embedding IS NOT NULL
			AND o.opera_id NOT IN (
				SELECT opera_id FROM favorites WHERE user_id = ?
				UNION
				SELECT opera_id FROM likes WHERE user_id = ?
			)
		)
		SELECT AVG(embedding) as embedding
		FROM user_operas
	`, userID, userID, userID, userID, userID, userID).Scan(&result).Error

	return result.Embedding, err
}

// RecommendByVector 基于向量相似度推荐作品
// 注意：涉及pgvector操作符(<=>)，必须使用Raw SQL
func (r *RecommendationRepo) RecommendByVector(userVector []float64, excludeOperaIDs []uint, similarityWeight float64, interestWeight float64, limit int) ([]uint, error) {
	type OperaScore struct {
		OperaID uint
	}

	var results []OperaScore
	var err error

	// 使用不同的查询语句避免字符串拼接SQL注入风险
	if len(excludeOperaIDs) > 0 {
		query := `
			WITH opera_scores AS (
				SELECT 
					o.opera_id,
					1 - (o.embedding <=> ?::vector) as similarity,
					COALESCE(
						(SELECT SUM(total_score)
						FROM user_interest_scores
						WHERE opera_id = o.opera_id), 0
					) as interest_score
				FROM opera o
				WHERE o.embedding IS NOT NULL
				AND o.is_hidden = false
				AND o.opera_id NOT IN ?
			)
			SELECT 
				opera_id
			FROM opera_scores
			ORDER BY (similarity * ? + LEAST(interest_score / 1000.0, 1.0) * ?) DESC
			LIMIT ?
		`
		err = r.getDB().Raw(query, userVector, excludeOperaIDs, similarityWeight, interestWeight, limit).Scan(&results).Error
	} else {
		query := `
			WITH opera_scores AS (
				SELECT 
					o.opera_id,
					1 - (o.embedding <=> ?::vector) as similarity,
					COALESCE(
						(SELECT SUM(total_score)
						FROM user_interest_scores
						WHERE opera_id = o.opera_id), 0
					) as interest_score
				FROM opera o
				WHERE o.embedding IS NOT NULL
				AND o.is_hidden = false
			)
			SELECT 
				opera_id
			FROM opera_scores
			ORDER BY (similarity * ? + LEAST(interest_score / 1000.0, 1.0) * ?) DESC
			LIMIT ?
		`
		err = r.getDB().Raw(query, userVector, similarityWeight, interestWeight, limit).Scan(&results).Error
	}

	if err != nil {
		return nil, err
	}

	operaIDs := make([]uint, len(results))
	for i, r := range results {
		operaIDs[i] = r.OperaID
	}

	return operaIDs, nil
}

// RecommendByInterestScore 基于兴趣分推荐作品
func (r *RecommendationRepo) RecommendByInterestScore(excludeOperaIDs []uint, limit int) ([]uint, error) {
	type OperaWithScore struct {
		OperaID uint    `gorm:"column:opera_id"`
		Score   float64 `gorm:"column:score"`
	}

	var results []OperaWithScore

	query := r.getDB().Table("opera o").
		Select("o.opera_id, COALESCE(scores.score, 0) as score").
		Joins("LEFT JOIN (SELECT opera_id, SUM(total_score) as score FROM user_interest_scores GROUP BY opera_id) scores ON o.opera_id = scores.opera_id").
		Where("o.is_hidden = ?", false)

	if len(excludeOperaIDs) > 0 {
		query = query.Where("o.opera_id NOT IN ?", excludeOperaIDs)
	}

	err := query.Order("score DESC").
		Limit(limit).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	operaIDs := make([]uint, len(results))
	for i, r := range results {
		operaIDs[i] = r.OperaID
	}

	return operaIDs, nil
}

// GetSimilarOperas 获取相似作品
// 注意：涉及pgvector操作符(<=>)，必须使用Raw SQL
func (r *RecommendationRepo) GetSimilarOperas(operaID uint, limit int) ([]uint, error) {
	var operaIDs []uint

	err := r.getDB().Raw(`
		SELECT o2.opera_id
		FROM opera o1
		CROSS JOIN opera o2
		WHERE o1.opera_id = ?
		AND o2.opera_id != ?
		AND o1.embedding IS NOT NULL
		AND o2.embedding IS NOT NULL
		AND o2.is_hidden = false
		ORDER BY o2.embedding <=> o1.embedding
		LIMIT ?
	`, operaID, operaID, limit).Pluck("opera_id", &operaIDs).Error

	return operaIDs, err
}

// GetUserInteractedOperaIDs 获取用户已交互的作品ID列表
func (r *RecommendationRepo) GetUserInteractedOperaIDs(userID uint) ([]uint, error) {
	var operaIDs []uint

	// GORM不支持UNION，使用Raw SQL
	err := r.getDB().Raw(`
		SELECT DISTINCT opera_id FROM (
			SELECT opera_id FROM favorites WHERE user_id = ?
			UNION
			SELECT opera_id FROM likes WHERE user_id = ?
			UNION
			SELECT DISTINCT opera_id FROM play_history WHERE user_id = ?
		) AS interacted
	`, userID, userID, userID).Pluck("opera_id", &operaIDs).Error

	return operaIDs, err
}
