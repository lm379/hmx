package services

import (
	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/repository"
)

var statsRepo = repository.NewStatsRepo()

// GetDashboardStats 获取仪表盘统计数据
func GetDashboardStats() (gin.H, error) {
	stats, err := statsRepo.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	return gin.H{
		"total_users":     stats.TotalUsers,
		"total_operas":    stats.TotalOperas,
		"total_artists":   stats.TotalArtists,
		"total_likes":     stats.TotalLikes,
		"total_plays":     stats.TotalPlays,
		"total_favorites": stats.TotalFavorites,
		"total_shares":    stats.TotalShares,
	}, nil
}
