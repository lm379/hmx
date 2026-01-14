package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/services"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleSearchOperas 搜索 Opera
// @Summary 搜索作品
// @Description 在作品索引中搜索
// @Tags Search
// @Accept json
// @Produce json
// @Param query query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param include_hidden query bool false "是否包含隐藏内容"
// @Param sort_by query string false "排序字段，如: created_at:desc"
// @Success 200 {object} services.SearchOperasResponse
// @Router /api/v1/search/operas [get]
func HandleSearchOperas(c *gin.Context) {
	var req services.SearchOperasRequest

	// 解析查询参数
	req.Query = c.Query("query")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	req.Page = page

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	req.PageSize = pageSize

	req.IncludeHidden = c.Query("include_hidden") == "true"
	req.SortBy = c.Query("sort_by")

	// 调用搜索服务
	searchService := services.NewSearchService()
	response, err := searchService.SearchOperas(req)
	if err != nil {
		resp.InternalServerError(c, "Search failed: "+err.Error())
		return
	}

	resp.Success(c, response)
}

// HandleSearchArtists 搜索 Artist
// @Summary 搜索艺术家
// @Description 在艺术家索引中搜索
// @Tags Search
// @Accept json
// @Produce json
// @Param query query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段，如: created_at:desc"
// @Success 200 {object} services.SearchArtistsResponse
// @Router /api/v1/search/artists [get]
func HandleSearchArtists(c *gin.Context) {
	var req services.SearchArtistsRequest

	// 解析查询参数
	req.Query = c.Query("query")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	req.Page = page

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	req.PageSize = pageSize

	req.SortBy = c.Query("sort_by")

	// 调用搜索服务
	searchService := services.NewSearchService()
	response, err := searchService.SearchArtists(req)
	if err != nil {
		resp.InternalServerError(c, "Search failed: "+err.Error())
		return
	}

	resp.Success(c, response)
}

// HandleReindexOperas 重新索引所有 Opera
// @Summary 重新索引所有作品
// @Description 批量重新索引所有作品到 Meilisearch (需要管理员权限)
// @Tags Search
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/search/reindex/operas [post]
func HandleReindexOperas(c *gin.Context) {
	searchService := services.NewSearchService()
	err := searchService.ReindexAllOperas()
	if err != nil {
		resp.InternalServerError(c, "Reindex failed: "+err.Error())
		return
	}

	resp.Success(c, gin.H{
		"message": "Successfully reindexed all operas",
	})
}

// HandleReindexArtists 重新索引所有 Artist
// @Summary 重新索引所有艺术家
// @Description 批量重新索引所有艺术家到 Meilisearch (需要管理员权限)
// @Tags Search
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/search/reindex/artists [post]
func HandleReindexArtists(c *gin.Context) {
	searchService := services.NewSearchService()
	err := searchService.ReindexAllArtists()
	if err != nil {
		resp.InternalServerError(c, "Reindex failed: "+err.Error())
		return
	}

	resp.Success(c, gin.H{
		"message": "Successfully reindexed all artists",
	})
}

// HandleGetSearchStats 获取搜索统计
// @Summary 获取搜索索引统计信息
// @Description 获取 Meilisearch 索引统计信息 (需要管理员权限)
// @Tags Search
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/search/stats [get]
func HandleGetSearchStats(c *gin.Context) {
	searchService := services.NewSearchService()
	stats, err := searchService.GetSearchStats()
	if err != nil {
		resp.InternalServerError(c, "Failed to get search stats: "+err.Error())
		return
	}

	resp.Success(c, stats)
}
