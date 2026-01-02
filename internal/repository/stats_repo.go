package repository

import (
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"gorm.io/gorm"
)

// StatsRepo 统计相关的数据库操作
type StatsRepo struct{}

// NewStatsRepo 创建统计仓库实例
func NewStatsRepo() *StatsRepo {
	return &StatsRepo{}
}

// getDB 获取数据库实例
func (r *StatsRepo) getDB() *gorm.DB {
	return database.DB
}

// DashboardStats 仪表盘统计数据结构
type DashboardStats struct {
	TotalUsers     int64
	TotalOperas    int64
	TotalArtists   int64
	TotalLikes     int64
	TotalPlays     int64
	TotalFavorites int64
	TotalShares    int64
}

// GetDashboardStats 获取仪表盘统计数据
func (r *StatsRepo) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}
	db := r.getDB()

	// 统计用户总数
	if err := db.Model(&models.Users{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// 统计作品总数
	if err := db.Model(&models.Opera{}).Count(&stats.TotalOperas).Error; err != nil {
		return nil, err
	}

	// 统计艺术家总数
	if err := db.Model(&models.Artist{}).Count(&stats.TotalArtists).Error; err != nil {
		return nil, err
	}

	// 统计点赞总数
	if err := db.Model(&models.Like{}).Count(&stats.TotalLikes).Error; err != nil {
		return nil, err
	}

	// 统计播放总数 (累计播放次数)
	if err := db.Model(&models.PlayHistory{}).Select("COALESCE(sum(count), 0)").Scan(&stats.TotalPlays).Error; err != nil {
		return nil, err
	}

	// 统计收藏总数
	if err := db.Model(&models.Favorite{}).Count(&stats.TotalFavorites).Error; err != nil {
		return nil, err
	}

	// 统计分享总数
	if err := db.Model(&models.Share{}).Count(&stats.TotalShares).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
