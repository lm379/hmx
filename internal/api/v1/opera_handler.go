package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/converter"
	"github.com/lm379/hmx/pkg/pagination"
	resp "github.com/lm379/hmx/pkg/response"
	"gorm.io/gorm"
)

// HandleCreateOpera (POST /api/v1/operas)
func HandleCreateOpera(c *gin.Context) {
	var input models.CreateOperaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	// 从 JWT 中间件获取 userID
	userID, exists := c.Get("userID")
	if !exists {
		resp.Unauthorized(c, "Unauthorized")
		return
	}

	// 调用 service 创建 Opera
	response, status := services.CreateOpera(input, userID.(uint))
	if status != http.StatusCreated {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Created(c, response)
}

// HandleGetOperas (GET /api/v1/operas)
func HandleGetOperas(c *gin.Context) {
	pagination := pagination.GetPagination(c)

	operas, total, err := services.GetOperas(pagination, false)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch operas")
		return
	}

	// 转换为列表响应格式（简化版）
	operaResponses := converter.ToOperaListResponseList(operas)

	// 批量获取播放和点赞计数
	if len(operaResponses) > 0 {
		operaIDs := make([]uint, len(operaResponses))
		for i, r := range operaResponses {
			operaIDs[i] = r.OperaID
		}

		likes, _, _, plays := services.BatchGetCounts(operaIDs)
		for _, r := range operaResponses {
			r.LikeCount = likes[r.OperaID]
			r.PlayCount = plays[r.OperaID]
		}
	}

	resp.Success(c, gin.H{
		"list": operaResponses,
		"pagination": gin.H{
			"total":     total,
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
		},
	})
}

// HandleGetOperaByID (GET /api/v1/operas/:id)
func HandleGetOperaByID(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	opera, err := services.GetOperaByID(uint(operaID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "Opera not found")
			return
		}
		resp.InternalServerError(c, "System error")
		return
	}

	// 转换为详情响应格式（完整版）
	operaResponse := converter.ToOperaDetailResponse(opera)
	// 计数与用户状态
	lc, fc, sc, pc := services.GetCounts(operaResponse.OperaID)
	operaResponse.LikeCount = lc
	operaResponse.FavoriteCount = fc
	operaResponse.ShareCount = sc
	operaResponse.PlayCount = pc
	if v, exists := c.Get("userID"); exists {
		uid := v.(uint)
		operaResponse.Liked = services.IsLiked(uid, operaResponse.OperaID)
		operaResponse.Favorited = services.IsFavorited(uid, operaResponse.OperaID)
	}
	resp.Success(c, operaResponse)
}

// HandleRequestOperaSummary (POST /api/v1/operas/:id/request-summary)
// 用户端请求生成AI摘要（需要登录，只能生成不存在摘要的视频）
func HandleRequestOperaSummary(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	// 检查用户是否登录
	_, exists := c.Get("userID")
	if !exists {
		resp.Unauthorized(c, "Please login first")
		return
	}

	// 获取作品
	opera, err := services.GetOperaByID(uint(operaID))
	if err != nil {
		resp.NotFound(c, "Opera not found")
		return
	}

	// 检查是否有字幕文件
	if !opera.SrtPath.Valid || opera.SrtPath.String == "" {
		resp.BadRequest(c, "Opera has no subtitle file")
		return
	}

	// 检查是否已有摘要，如果有则返回204
	if opera.AiSummary != "" {
		resp.NoContentWithMsg(c, "Summary already exists")
		return
	}

	// 将任务加入队列（异步，force=false）
	batchID, _, err := services.BatchGenerateOperaSummariesAsync([]uint{uint(operaID)}, false)
	if err != nil {
		resp.InternalServerError(c, "Failed to start summary generation: "+err.Error())
		return
	}

	resp.Success(c, gin.H{
		"message":  "Summary generation started",
		"task_id":  batchID,
		"opera_id": operaID,
	})
}
