package services

import (
	"fmt"

	"github.com/lm379/hmx/internal/repository"
)

// RecommendationService 推荐服务
type RecommendationService struct {
	operaRepo          *repository.OperaRepo
	interactionRepo    *repository.InteractionRepo
	recommendationRepo *repository.RecommendationRepo
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService() *RecommendationService {
	return &RecommendationService{
		operaRepo:          repository.NewOperaRepo(),
		interactionRepo:    repository.NewInteractionRepo(),
		recommendationRepo: repository.NewRecommendationRepo(),
	}
}

// RecommendForUser 为用户推荐作品（混合推荐算法，支持分页）
// userID: 用户ID，0 表示游客
// limit: 每页数量
// page: 页码（从1开始）
// 如果推荐结果少于limit，自动用热门视频补充至limit条（去重）
func (s *RecommendationService) RecommendForUser(userID uint, page int, limit int, similarityWeight float64, interestWeight float64) ([]uint, int64, error) {
	if page < 1 {
		page = 1
	}

	// 获取用户已交互的作品ID（用于排除）
	excludeIDs := []uint{}
	if userID != 0 {
		interactedIDs, err := s.recommendationRepo.GetUserInteractedOperaIDs(userID)
		if err == nil {
			excludeIDs = interactedIDs
		}
	}

	var recommendedIDs []uint
	var err error

	// 游客或新用户直接使用热门推荐
	if userID == 0 {
		recommendedIDs, err = s.recommendationRepo.RecommendByInterestScore(nil, limit*page)
		if err != nil {
			return nil, 0, err
		}
	} else {
		// 尝试获取用户偏好向量
		userVector, vectorErr := s.recommendationRepo.GetUserPreferenceVector(userID)
		if vectorErr != nil || len(userVector) == 0 {
			// 降级为纯兴趣分推荐
			recommendedIDs, err = s.recommendationRepo.RecommendByInterestScore(excludeIDs, limit*page)
			if err != nil {
				return nil, 0, err
			}
		} else {
			// 使用混合推荐算法
			recommendedIDs, err = s.recommendationRepo.RecommendByVector(userVector, excludeIDs, similarityWeight, interestWeight, limit*page)
			if err != nil {
				return nil, 0, err
			}
		}
	}

	// 总数（用于前端显示）
	total := int64(len(recommendedIDs))

	// 如果推荐结果少于请求的总数，用热门视频补充
	minRequired := limit * page
	if len(recommendedIDs) < minRequired {
		// 计算需要补充多少条
		needed := minRequired - len(recommendedIDs)

		// 获取热门视频（排除已推荐的）
		excludeForPopular := append(excludeIDs, recommendedIDs...)
		popularIDs, err := s.operaRepo.GetPopularOperas(excludeForPopular, needed)
		if err == nil && len(popularIDs) > 0 {
			recommendedIDs = append(recommendedIDs, popularIDs...)
			total = int64(len(recommendedIDs))
		}
	}

	// 分页处理
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit

	if startIdx >= len(recommendedIDs) {
		return []uint{}, total, nil
	}

	if endIdx > len(recommendedIDs) {
		endIdx = len(recommendedIDs)
	}

	pageResults := recommendedIDs[startIdx:endIdx]

	return pageResults, total, nil
}

// RecommendSimilarOperas 推荐与指定作品相似的其他作品
func (s *RecommendationService) RecommendSimilarOperas(operaID uint, limit int) ([]uint, error) {
	operaIDs, err := s.recommendationRepo.GetSimilarOperas(operaID, limit)
	if err != nil {
		return nil, err
	}

	if len(operaIDs) == 0 {
		return nil, fmt.Errorf("no similar operas found")
	}

	return operaIDs, nil
}
