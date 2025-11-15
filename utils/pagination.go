package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Pagination 结构体用于存储分页信息
type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Sort     string `json:"sort"`
}

// GetPagination 从 Gin Context 中解析分页查询参数
func GetPagination(c *gin.Context) *Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sort := c.DefaultQuery("sort", "created_at desc") // 默认按创建时间倒序

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 { // 限制每页最大数量
		pageSize = 100
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	}
}

// Paginate GORM 的分页 Scope
func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.PageSize
		return db.Offset(offset).Limit(p.PageSize).Order(p.Sort)
	}
}
