package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/pagination"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleGetRecommendations 获取个性化推荐（支持分页）
func HandleGetRecommendations(c *gin.Context) {
	// 从中间件获取用户ID（可选，支持游客）
	var userID uint = 0
	if uid, exists := c.Get("userID"); exists {
		userID = uid.(uint)
	}

	// 获取分页参数
	paginationParams := pagination.GetPagination(c)
	page := paginationParams.Page
	pageSize := paginationParams.PageSize

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 相似度权重和兴趣分权重（默认各占50%）
	similarityWeightStr := c.DefaultQuery("similarity_weight", "0.5")
	interestWeightStr := c.DefaultQuery("interest_weight", "0.5")

	similarityWeight, err := strconv.ParseFloat(similarityWeightStr, 64)
	if err != nil || similarityWeight < 0 || similarityWeight > 1 {
		similarityWeight = 0.5
	}

	interestWeight, err := strconv.ParseFloat(interestWeightStr, 64)
	if err != nil || interestWeight < 0 || interestWeight > 1 {
		interestWeight = 0.5
	}

	channel := c.DefaultQuery("channel", services.RecommendationChannelForYou)
	if channel != services.RecommendationChannelForYou &&
		channel != services.RecommendationChannelHot &&
		channel != services.RecommendationChannelLatest {
		channel = services.RecommendationChannelForYou
	}

	// 调用推荐服务（带分页）
	recService := services.NewRecommendationService()
	operaIDs, total, reason, err := recService.RecommendByChannel(channel, userID, page, pageSize, similarityWeight, interestWeight)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "Failed to get recommendations: "+err.Error())
		return
	}

	if len(operaIDs) == 0 {
		resp.Success(c, gin.H{
			"operas": []interface{}{},
			"total":  0,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      0,
				"total_page": 0,
			},
		})
		return
	}

	// 获取作品详情
	operaService := services.NewOperaService()
	operas, err := operaService.GetOperasByIDs(operaIDs, userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "Failed to fetch opera details")
		return
	}
	for i := range operas {
		operas[i].RecommendReason = reason
	}

	// 计算总页数
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	resp.Success(c, gin.H{
		"operas": operas,
		"total":  total,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": totalPage,
		},
	})
}

// HandleGetSimilarOperas 获取相似作品推荐
func HandleGetSimilarOperas(c *gin.Context) {
	// 获取作品ID
	operaIDStr := c.Param("id")
	operaID, err := strconv.ParseUint(operaIDStr, 10, 32)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "Invalid opera ID")
		return
	}

	// 获取推荐数量
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	// 调用推荐服务
	recService := services.NewRecommendationService()
	operaIDs, err := recService.RecommendSimilarOperas(uint(operaID), limit)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "Failed to get similar operas: "+err.Error())
		return
	}

	if len(operaIDs) == 0 {
		resp.Success(c, gin.H{
			"operas": []interface{}{},
			"total":  0,
		})
		return
	}

	// 获取用户ID（可选）
	var userID uint = 0
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	// 获取作品详情
	operaService := services.NewOperaService()
	operas, err := operaService.GetOperasByIDs(operaIDs, userID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "Failed to fetch opera details")
		return
	}

	resp.Success(c, gin.H{
		"operas": operas,
		"total":  len(operas),
	})
}

// HandleGenerateEmbedding 手动触发生成作品向量（管理员功能）
func HandleGenerateEmbedding(c *gin.Context) {
	operaIDStr := c.Param("id")
	operaID, err := strconv.ParseUint(operaIDStr, 10, 32)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "Invalid opera ID")
		return
	}

	// 获取请求体（可选force参数）
	var req struct {
		Force bool `json:"force"` // 是否强制重新生成
	}
	c.ShouldBindJSON(&req)

	// 如果不是强制生成，检查是否已有向量
	if !req.Force {
		opera, err := services.GetOperaByID(uint(operaID))
		if err == nil {
			slice := opera.Embedding.Slice()
			if len(slice) > 0 {
				resp.NoContentWithMsg(c, "Embedding already exists")
				return
			}
		}
	}

	// 将任务加入队列（异步）
	batchID, _, err := services.BatchGenerateEmbeddingsAsync([]uint{uint(operaID)}, req.Force)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "Failed to start embedding generation: "+err.Error())
		return
	}

	resp.Success(c, gin.H{
		"message":  "Embedding generation started",
		"task_id":  batchID,
		"opera_id": operaID,
	})
}

// HandleBatchGenerateEmbeddings 批量生成所有作品的向量（管理员功能）
func HandleBatchGenerateEmbeddings(c *gin.Context) {
	var req struct {
		OperaIDs []uint `json:"opera_ids"`
		Force    bool   `json:"force"` // 是否强制重新生成
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "Invalid request")
		return
	}

	if len(req.OperaIDs) == 0 {
		resp.BadRequest(c, "No opera IDs provided")
		return
	}

	batchID, total, err := services.BatchGenerateEmbeddingsAsync(req.OperaIDs, req.Force)
	if err != nil {
		resp.InternalServerError(c, "Failed to start batch task: "+err.Error())
		return
	}

	if total == 0 {
		resp.NoContentWithMsg(c, "All operas already have embeddings")
		return
	}

	resp.Success(c, gin.H{
		"message":  "Batch embedding generation started",
		"batch_id": batchID,
		"total":    total,
	})
}
