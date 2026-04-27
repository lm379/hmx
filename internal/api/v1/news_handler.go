package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/pagination"
	resp "github.com/lm379/hmx/pkg/response"
	"gorm.io/gorm"
)

// HandleGetNews 返回前台可见的已发布新闻列表。
func HandleGetNews(c *gin.Context) {
	p := pagination.GetPagination(c)
	newsService := services.NewNewsService()
	list, total, err := newsService.GetPublishedNews(p)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch news")
		return
	}

	resp.Success(c, gin.H{
		"list": list,
		"pagination": gin.H{
			"total":     total,
			"page":      p.Page,
			"page_size": p.PageSize,
		},
	})
}

// HandleGetNewsByID 返回单条已发布新闻详情，草稿不会从前台接口暴露。
func HandleGetNewsByID(c *gin.Context) {
	newsID, ok := parseNewsID(c)
	if !ok {
		return
	}

	newsService := services.NewNewsService()
	detail, err := newsService.GetPublishedNewsDetail(newsID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "News not found")
			return
		}
		resp.InternalServerError(c, "Failed to fetch news")
		return
	}
	resp.Success(c, detail)
}

// HandleAdminGetNews 返回后台新闻列表，包含草稿和已发布内容。
func HandleAdminGetNews(c *gin.Context) {
	p := pagination.GetPagination(c)
	newsService := services.NewNewsService()
	list, total, err := newsService.GetAllNews(p)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch news")
		return
	}
	resp.Success(c, gin.H{
		"list": list,
		"pagination": gin.H{
			"total":     total,
			"page":      p.Page,
			"page_size": p.PageSize,
		},
	})
}

// HandleAdminGetNewsByID 返回后台编辑所需的完整新闻详情。
func HandleAdminGetNewsByID(c *gin.Context) {
	newsID, ok := parseNewsID(c)
	if !ok {
		return
	}

	newsService := services.NewNewsService()
	detail, err := newsService.GetAdminNewsDetail(newsID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "News not found")
			return
		}
		resp.InternalServerError(c, "Failed to fetch news")
		return
	}
	resp.Success(c, detail)
}

// HandleAdminCreateNews 创建新闻；若直接发布，服务层会自动补发布时间。
func HandleAdminCreateNews(c *gin.Context) {
	var input models.CreateNewsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	newsService := services.NewNewsService()
	detail, err := newsService.CreateNews(input)
	if err != nil {
		resp.InternalServerError(c, "Failed to create news")
		return
	}
	resp.Created(c, detail)
}

// HandleAdminUpdateNews 更新新闻完整表单，支持发布/下架切换。
func HandleAdminUpdateNews(c *gin.Context) {
	newsID, ok := parseNewsID(c)
	if !ok {
		return
	}

	var input models.UpdateNewsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	newsService := services.NewNewsService()
	detail, err := newsService.UpdateNews(newsID, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "News not found")
			return
		}
		resp.InternalServerError(c, "Failed to update news")
		return
	}
	resp.Success(c, detail)
}

// HandleAdminDeleteNews 删除新闻记录。
func HandleAdminDeleteNews(c *gin.Context) {
	newsID, ok := parseNewsID(c)
	if !ok {
		return
	}

	newsService := services.NewNewsService()
	if err := newsService.DeleteNews(newsID); err != nil {
		resp.InternalServerError(c, "Failed to delete news")
		return
	}
	resp.Success(c, gin.H{"message": "News deleted"})
}

// parseNewsID 统一解析路由中的新闻 ID，并直接写出 400 响应。
func parseNewsID(c *gin.Context) (uint, bool) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid news ID")
		return 0, false
	}
	return uint(id), true
}
