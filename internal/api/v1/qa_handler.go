package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleAsk (POST /api/v1/qa/ask)
// 支持游客（session_id）和已登录用户
func HandleAsk(c *gin.Context) {
	var req models.AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if req.Query == "" {
		resp.BadRequest(c, "问题不能为空")
		return
	}

	// 获取用户 ID（可选）
	var userID *uuid.UUID
	if rawID, exists := c.Get("userID"); exists {
		if uid, ok := rawID.(uint); ok && uid > 0 {
			// 将 uint 转换为 uuid 需要从数据库查找，这里直接使用数值 ID
			// 由于现有系统 users.user_id 是 uint 而不是 uuid，需要特殊处理
			// 我们使用一个确定性的 UUID namespace
			u := uintToUUID(uid)
			userID = &u
		}
	}

	// 调用 RAG 核心流程
	result, err := services.AskQuestion(req.Query, userID, req.SessionID, req.OperaID)
	if err != nil {
		resp.InternalServerError(c, "问答处理失败: "+err.Error())
		return
	}

	resp.Success(c, result)
}

// HandleGetQAHistory (GET /api/v1/qa/history)
// TryAuth：登录用户查自己历史，游客通过 session_id 查询
func HandleGetQAHistory(c *gin.Context) {
	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var userID *uuid.UUID
	sessionID := c.Query("session_id")

	if rawID, exists := c.Get("userID"); exists {
		if uid, ok := rawID.(uint); ok && uid > 0 {
			u := uintToUUID(uid)
			userID = &u
		}
	}

	if userID == nil && sessionID == "" {
		resp.BadRequest(c, "需要登录或提供 session_id")
		return
	}

	items, total, err := services.GetQAHistory(userID, sessionID, page, pageSize)
	if err != nil {
		resp.InternalServerError(c, "获取历史失败: "+err.Error())
		return
	}

	resp.Success(c, gin.H{
		"list": items,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// HandleSubmitFeedback (POST /api/v1/qa/:qa_id/feedback)
func HandleSubmitFeedback(c *gin.Context) {
	qaIDStr := c.Param("qa_id")
	qaID, err := uuid.Parse(qaIDStr)
	if err != nil {
		resp.BadRequest(c, "无效的 qa_id")
		return
	}

	var req models.FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	var userID *uuid.UUID
	if rawID, exists := c.Get("userID"); exists {
		if uid, ok := rawID.(uint); ok && uid > 0 {
			u := uintToUUID(uid)
			userID = &u
		}
	}

	if err := services.SubmitFeedback(qaID, userID, req.Rating, req.Comments); err != nil {
		resp.InternalServerError(c, "提交反馈失败: "+err.Error())
		return
	}

	resp.Success(c, gin.H{"success": true})
}

// uintToUUID 将 uint 类型的用户 ID 转换为确定性 UUID（用于 QA 历史关联）
// 使用 UUID v5 namespace + 用户 ID 字符串生成确定性 UUID
func uintToUUID(id uint) uuid.UUID {
	ns := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // DNS namespace
	return uuid.NewSHA1(ns, []byte("user:"+strconv.FormatUint(uint64(id), 10)))
}
