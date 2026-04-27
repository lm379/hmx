package services

import (
	"fmt"

	"github.com/lm379/hmx/internal/repository"
)

const (
	RecommendationChannelForYou = "for_you"
	RecommendationChannelHot    = "hot"
	RecommendationChannelLatest = "latest"
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
	var total int64

	// 计算真实的可用总数 = 全局可见总数 - 排除的ID数
	totalVisible, err := s.operaRepo.CountVisible()
	if err != nil {
		// 如果获取失败，回退到当前切片长度
		total = int64(len(recommendedIDs))
	} else {
		total = totalVisible - int64(len(excludeIDs))
		if total < 0 {
			total = 0
		}
	}

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
			// 注意：total不需要更新，因为popularIDs本来就在totalVisible里
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

// RecommendByChannel 为首页不同频道推荐作品。
func (s *RecommendationService) RecommendByChannel(channel string, userID uint, page int, limit int, similarityWeight float64, interestWeight float64) ([]uint, int64, string, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	switch channel {
	case RecommendationChannelHot:
		operaIDs, err := s.recommendationRepo.RecommendByInterestScore(nil, limit*page)
		if err != nil {
			return nil, 0, "", err
		}
		total, err := s.operaRepo.CountVisible()
		if err != nil {
			total = int64(len(operaIDs))
		}
		return paginateIDs(operaIDs, page, limit), total, "互动热度较高", nil
	case RecommendationChannelLatest:
		offset := (page - 1) * limit
		operaIDs, total, err := s.operaRepo.GetLatestOperas(limit, offset)
		return operaIDs, total, "最新收录", err
	default:
		operaIDs, total, err := s.RecommendForUser(userID, page, limit, similarityWeight, interestWeight)
		if err != nil {
			return nil, 0, "", err
		}
		if userID == 0 {
			return operaIDs, total, "全站热门内容", nil
		}
		return operaIDs, total, "为你推荐", nil
	}
}

func paginateIDs(ids []uint, page int, limit int) []uint {
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit
	if startIdx >= len(ids) {
		return []uint{}
	}
	if endIdx > len(ids) {
		endIdx = len(ids)
	}
	return ids[startIdx:endIdx]
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
