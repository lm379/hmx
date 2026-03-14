package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/converter"
	"github.com/lm379/hmx/pkg/pagination"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleDeleteOpera (DELETE /api/v1/admin/operas/:id)
func HandleDeleteOpera(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	response, status := services.DeleteOpera(uint(operaID))
	if status >= 400 {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleAdminGetOperas (GET /api/v1/admin/operas)
func HandleAdminGetOperas(c *gin.Context) {
	pagination := pagination.GetPagination(c)

	// 获取所有视频，包括隐藏的
	operas, total, err := services.GetOperas(pagination, true)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch operas")
		return
	}

	operaResponses := converter.ToOperaResponseList(operas)

	// Count stats
	if len(operaResponses) > 0 {
		operaIDs := make([]uint, len(operaResponses))
		for i, r := range operaResponses {
			operaIDs[i] = r.OperaID
		}

		likes, favorites, shares, plays := services.BatchGetCounts(operaIDs)
		for _, r := range operaResponses {
			r.LikeCount = likes[r.OperaID]
			r.FavoriteCount = favorites[r.OperaID]
			r.ShareCount = shares[r.OperaID]
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

// HandleAdminUpdateOpera (PUT /api/v1/admin/operas/:id)
func HandleAdminUpdateOpera(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	var input models.UpdateOperaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	response, status := services.UpdateOpera(uint(operaID), input)
	if status != http.StatusOK {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleAdminCreateArtist (POST /api/v1/admin/artists)
func HandleAdminCreateArtist(c *gin.Context) {
	var input models.CreateArtistRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	response, status := services.CreateArtist(input)
	if status != http.StatusCreated {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Created(c, response)
}

// HandleAdminUpdateArtist (PUT /api/v1/admin/artists/:id)
func HandleAdminUpdateArtist(c *gin.Context) {
	idParam := c.Param("id")
	artistID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid artist ID")
		return
	}

	var input models.UpdateArtistRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	response, status := services.UpdateArtist(uint(artistID), input)
	if status != http.StatusOK {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleAdminDeleteArtist (DELETE /api/v1/admin/artists/:id)
func HandleAdminDeleteArtist(c *gin.Context) {
	idParam := c.Param("id")
	artistID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid artist ID")
		return
	}

	response, status := services.DeleteArtist(uint(artistID))
	if status != http.StatusOK {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleAdminGetUsers (GET /api/v1/admin/users)
func HandleAdminGetUsers(c *gin.Context) {
	pagination := pagination.GetPagination(c)

	users, total, err := services.GetAllUsers(pagination)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch users")
		return
	}

	userResponses := make([]*models.UserResponse, len(users))
	for i, u := range users {
		// Create a copy of the loop variable to avoid pointing to the same address
		user := u
		userResponses[i] = converter.ToUserResponse(&user)
	}

	resp.Success(c, gin.H{
		"list": userResponses,
		"pagination": gin.H{
			"total":     total,
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
		},
	})
}

// HandleAdminUpdateUserRole (PUT /api/v1/admin/users/:id/role)
func HandleAdminUpdateUserRole(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid user ID")
		return
	}

	var input models.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	if err := services.UpdateUserRole(uint(userID), input.Role); err != nil {
		resp.InternalServerError(c, "Failed to update role")
		return
	}
	resp.Success(c, gin.H{"message": "Role updated successfully"})
}

// HandleAdminGetDashboardStats (GET /api/v1/admin/stats)
func HandleAdminGetDashboardStats(c *gin.Context) {
	stats, err := services.GetDashboardStats()
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch dashboard stats")
		return
	}
	resp.Success(c, stats)
}

// HandleGenerateOperaSummary (POST /api/v1/admin/operas/:id/generate-summary)
func HandleGenerateOperaSummary(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	// 获取请求体（可选force参数）
	var req struct {
		Force bool `json:"force"` // 是否强制重新生成
	}
	c.ShouldBindJSON(&req)

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

	// 如果不是强制生成且已有摘要,返回204
	if !req.Force && opera.AiSummary != "" {
		resp.NoContentWithMsg(c, "Summary already exists")
		return
	}

	// 将任务加入队列（异步）
	batchID, _, err := services.BatchGenerateOperaSummariesAsync([]uint{uint(operaID)}, req.Force)
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

// HandleBatchGenerateOperaSummaries (POST /api/v1/admin/operas/batch-summary)
func HandleBatchGenerateOperaSummaries(c *gin.Context) {
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

	batchID, total, err := services.BatchGenerateOperaSummariesAsync(req.OperaIDs, req.Force)
	if err != nil {
		resp.InternalServerError(c, "Failed to start batch task: "+err.Error())
		return
	}

	if total == 0 {
		resp.NoContentWithMsg(c, "All operas already have summaries")
		return
	}

	resp.Success(c, gin.H{
		"message":  "Batch summary generation started",
		"batch_id": batchID,
		"total":    total,
	})
}

// =============================================
// 知识库管理端点（Admin Knowledge）
// =============================================

// HandleAdminUploadKnowledge (POST /api/v1/admin/knowledge/upload)
// multipart/form-data: title, source_type, file (.txt/.md)
func HandleAdminUploadKnowledge(c *gin.Context) {
	var req models.UploadDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取操作者 ID
	rawID, exists := c.Get("userID")
	if !exists {
		resp.Unauthorized(c, "未认证")
		return
	}
	adminUserID := rawID.(uint)

	// 创建文档记录
	doc, err := services.UploadDocument(req.Title, req.Content, req.SourceType, adminUserID)
	if err != nil {
		resp.InternalServerError(c, "上传文档失败: "+err.Error())
		return
	}

	// 入队向量化
	services.EnqueueDocumentEmbedding(doc.DocID)

	resp.Success(c, gin.H{
		"doc_id": doc.DocID,
		"status": "processing",
	})
}

// HandleAdminGetKnowledgeDocuments (GET /api/v1/admin/knowledge/documents)
func HandleAdminGetKnowledgeDocuments(c *gin.Context) {
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

	items, total, err := services.GetDocuments(page, pageSize)
	if err != nil {
		resp.InternalServerError(c, "获取文档列表失败: "+err.Error())
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

// HandleAdminDeleteKnowledgeDocument (DELETE /api/v1/admin/knowledge/:doc_id)
func HandleAdminDeleteKnowledgeDocument(c *gin.Context) {
	docIDStr := c.Param("doc_id")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		resp.BadRequest(c, "无效的 doc_id")
		return
	}

	if err := services.SoftDeleteDocument(docID); err != nil {
		resp.InternalServerError(c, "删除文档失败: "+err.Error())
		return
	}

	resp.Success(c, gin.H{"success": true})
}

// HandleAdminImportOperas (POST /api/v1/admin/knowledge/import-operas)
// 一次性触发：将所有已显示作品的字幕导入知识库
func HandleAdminImportOperas(c *gin.Context) {
	rawID, exists := c.Get("userID")
	if !exists {
		resp.Unauthorized(c, "未认证")
		return
	}
	adminUserID := rawID.(uint)

	// 异步执行
	go services.ImportAllOperasToKnowledge(adminUserID)

	resp.Success(c, gin.H{
		"message": "作品字幕导入任务已启动，正在后台处理",
	})
}

// HandleAdminGetKnowledgeDocumentDetail (GET /api/v1/admin/knowledge/documents/:doc_id)
// 获取单个知识库文档详情（含全文内容）
func HandleAdminGetKnowledgeDocumentDetail(c *gin.Context) {
	docIDStr := c.Param("doc_id")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		resp.BadRequest(c, "无效的 doc_id")
		return
	}

	doc, err := services.GetDocumentDetail(docID)
	if err != nil {
		resp.NotFound(c, "文档不存在")
		return
	}

	resp.Success(c, doc)
}

// HandleAdminReactivateKnowledgeDocument (PUT /api/v1/admin/knowledge/documents/:doc_id/reactivate)
// 重新激活已软删除的文档，并重新触发向量化
func HandleAdminReactivateKnowledgeDocument(c *gin.Context) {
	docIDStr := c.Param("doc_id")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		resp.BadRequest(c, "无效的 doc_id")
		return
	}

	if err := services.ReactivateDocument(docID); err != nil {
		resp.InternalServerError(c, "激活文档失败: "+err.Error())
		return
	}

	resp.Success(c, gin.H{"success": true, "message": "文档已激活，正在重新向量化"})
}

// HandleAdminGetKnowledgeStats (GET /api/v1/admin/knowledge/stats)
// 获取知识库统计信息
func HandleAdminGetKnowledgeStats(c *gin.Context) {
	stats := services.GetKnowledgeStats()
	resp.Success(c, stats)
}
