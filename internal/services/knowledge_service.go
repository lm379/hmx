package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/queue"
	"github.com/lm379/hmx/pkg/ai"
	"github.com/pgvector/pgvector-go"
)

// fetchSubtitleContent 从 CDN/对象存储读取字幕文件内容
func fetchSubtitleContent(srtPath string) (string, error) {
	cdnURL := config.AppConfig.S3CustomDomain
	if cdnURL == "" {
		cdnURL = config.AppConfig.S3Endpoint + "/" + config.AppConfig.S3Bucket
	}
	if len(srtPath) > 0 && srtPath[0] == '/' {
		srtPath = srtPath[1:]
	}
	fullURL := cdnURL + "/" + srtPath

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("获取字幕文件失败: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取字幕文件 HTTP 错误: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取字幕内容失败: %w", err)
	}
	return string(body), nil
}

const (
	chunkSize    = 512 // 每块目标字符数（中文约 512 字）
	chunkOverlap = 50  // 相邻块重叠字符数
	topK         = 5   // 检索返回 Top-K 块
)

// UploadDocument 创建知识库文档记录
func UploadDocument(title, content string, sourceType models.SourceType, createdBy uint) (*models.KnowledgeDocument, error) {
	doc := &models.KnowledgeDocument{
		DocID:      uuid.New(),
		Title:      title,
		Content:    content,
		SourceType: sourceType,
		CreatedBy:  createdBy,
		IsActive:   true,
	}
	if err := database.DB.Create(doc).Error; err != nil {
		return nil, fmt.Errorf("创建知识库文档失败: %w", err)
	}
	return doc, nil
}

// EnqueueDocumentEmbedding 将文档向量化任务写入 Redis 队列，并启动 goroutine 驱动 Worker 消费
func EnqueueDocumentEmbedding(docID uuid.UUID) {
	task := &queue.Task{
		ID:    docID.String(),
		Type:  queue.TaskTypeDocumentEmbedding,
		DocID: docID.String(),
	}
	if err := queue.EnqueueTask(task); err != nil {
		log.Printf("[KnowledgeService] 向量化任务入队失败 (doc %s): %v，回退到直接处理", docID, err)
		// Redis 不可用时降级为直接 goroutine 处理
		go func() {
			if err := ProcessDocumentChunks(docID); err != nil {
				log.Printf("[KnowledgeService] 降级处理失败 (doc %s): %v", docID, err)
			}
		}()
		return
	}
	log.Printf("[KnowledgeService] 向量化任务已入队 (doc %s)", docID)
}

// ProcessDocumentChunks 将文档分块、向量化并存储（同步，适合后台任务调用）
func ProcessDocumentChunks(docID uuid.UUID) error {
	var doc models.KnowledgeDocument
	if err := database.DB.Where("doc_id = ? AND is_active = true", docID).First(&doc).Error; err != nil {
		return fmt.Errorf("找不到文档 %s: %w", docID, err)
	}

	// 标记为处理中
	database.DB.Model(&doc).Update("embedding_status", models.EmbeddingStatusProcessing)

	chunks := splitIntoChunks(doc.Content, chunkSize, chunkOverlap)
	log.Printf("[KnowledgeService] 文档 %s 分块数: %d", docID, len(chunks))

	successCount := 0
	for i, chunk := range chunks {
		// 为每块生成向量
		embedding, err := ai.GenerateEmbedding(chunk)
		if err != nil {
			log.Printf("[KnowledgeService] 块 %d 向量化失败: %v", i, err)
			continue
		}

		ke := &models.KnowledgeEmbedding{
			ChunkID:    uuid.New(),
			DocID:      docID,
			ChunkIndex: i,
			ChunkText:  chunk,
			Embedding:  pgvector.NewVector(float32Slice(embedding)),
		}
		if err := database.DB.Create(ke).Error; err != nil {
			log.Printf("[KnowledgeService] 存储块 %d 失败: %v", i, err)
			continue
		}
		successCount++
	}

	// 更新文档状态和分块数
	status := models.EmbeddingStatusCompleted
	if successCount == 0 && len(chunks) > 0 {
		status = models.EmbeddingStatusFailed
	}
	database.DB.Model(&doc).Updates(map[string]interface{}{
		"embedding_status": status,
		"chunks_count":     successCount,
	})

	log.Printf("[KnowledgeService] 文档 %s 向量化完成，状态: %s，成功块数: %d/%d", docID, status, successCount, len(chunks))
	return nil
}

// RetrieveRelevantChunks 基于余弦相似度检索相关文档块
func RetrieveRelevantChunks(query string, topKOverride int) ([]RetrievedChunk, error) {
	k := topK
	if topKOverride > 0 {
		k = topKOverride
	}

	// 生成查询向量
	queryEmbedding, err := ai.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("生成查询向量失败: %w", err)
	}

	vec := pgvector.NewVector(float32Slice(queryEmbedding))

	// 使用 pgvector cosine 距离检索（<=> 为余弦距离，越小越相似）
	type result struct {
		ChunkID   uuid.UUID
		DocID     uuid.UUID
		ChunkText string
		DocTitle  string
		Score     float64
	}

	var rows []result
	sqlQuery := `
		SELECT ke.chunk_id, ke.doc_id, ke.chunk_text, kd.title AS doc_title,
		       1 - (ke.embedding <=> ?) AS score
		FROM knowledge_embeddings ke
		JOIN knowledge_documents kd ON ke.doc_id = kd.doc_id
		WHERE kd.is_active = true
		ORDER BY ke.embedding <=> ?
		LIMIT ?`

	if err := database.DB.Raw(sqlQuery, vec, vec, k).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("向量检索失败: %w", err)
	}

	chunks := make([]RetrievedChunk, len(rows))
	for i, r := range rows {
		chunks[i] = RetrievedChunk{
			ChunkID:   r.ChunkID,
			DocID:     r.DocID,
			DocTitle:  r.DocTitle,
			ChunkText: r.ChunkText,
			Score:     r.Score,
		}
	}
	return chunks, nil
}

// RetrievedChunk 检索结果
type RetrievedChunk struct {
	ChunkID   uuid.UUID
	DocID     uuid.UUID
	DocTitle  string
	ChunkText string
	Score     float64
}

// GetDocuments 获取文档列表（分页）
func GetDocuments(page, pageSize int) ([]models.DocumentListItem, int64, error) {
	var docs []models.KnowledgeDocument
	var total int64

	db := database.DB.Model(&models.KnowledgeDocument{}).Where("is_active = true")
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&docs).Error; err != nil {
		return nil, 0, fmt.Errorf("获取文档列表失败: %w", err)
	}

	items := make([]models.DocumentListItem, len(docs))
	for i, doc := range docs {
		items[i] = models.DocumentListItem{
			DocID:           doc.DocID,
			Title:           doc.Title,
			SourceType:      doc.SourceType,
			ChunksCount:     doc.ChunksCount,
			EmbeddingStatus: doc.EmbeddingStatus,
			IsActive:        doc.IsActive,
			CreatedAt:       doc.CreatedAt,
		}
	}
	return items, total, nil
}

// GetDocumentDetail 获取单个文档详情（含内容）
func GetDocumentDetail(docID uuid.UUID) (*models.DocumentDetailResponse, error) {
	var doc models.KnowledgeDocument
	if err := database.DB.Where("doc_id = ?", docID).First(&doc).Error; err != nil {
		return nil, fmt.Errorf("文档不存在: %w", err)
	}

	// 查询分块列表（按 chunk_index 排序，不加载 embedding 向量节省带宽）
	var embeddings []models.KnowledgeEmbedding
	database.DB.Select("chunk_id, chunk_index, chunk_text").
		Where("doc_id = ?", docID).
		Order("chunk_index ASC").
		Find(&embeddings)

	chunks := make([]models.ChunkItem, len(embeddings))
	for i, e := range embeddings {
		chunks[i] = models.ChunkItem{
			ChunkID:    e.ChunkID,
			ChunkIndex: e.ChunkIndex,
			ChunkText:  e.ChunkText,
		}
	}

	return &models.DocumentDetailResponse{
		DocID:           doc.DocID,
		Title:           doc.Title,
		Content:         doc.Content,
		SourceType:      doc.SourceType,
		ChunksCount:     doc.ChunksCount,
		EmbeddingStatus: doc.EmbeddingStatus,
		IsActive:        doc.IsActive,
		CreatedBy:       doc.CreatedBy,
		CreatedAt:       doc.CreatedAt,
		UpdatedAt:       doc.UpdatedAt,
		Chunks:          chunks,
	}, nil
}

// ReactivateDocument 重新激活文档（is_active=true），并重新触发向量化
func ReactivateDocument(docID uuid.UUID) error {
	result := database.DB.Model(&models.KnowledgeDocument{}).
		Where("doc_id = ?", docID).
		Updates(map[string]interface{}{
			"is_active":        true,
			"embedding_status": models.EmbeddingStatusPending,
		})
	if result.Error != nil {
		return fmt.Errorf("激活文档失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("文档不存在")
	}
	// 入队重新向量化
	EnqueueDocumentEmbedding(docID)
	return nil
}

// GetKnowledgeStats 获取知识库统计信息
func GetKnowledgeStats() map[string]interface{} {
	var total, professional, opera, completed, pending, failed int64
	var totalChunks int64

	database.DB.Model(&models.KnowledgeDocument{}).Where("is_active = true").Count(&total)
	database.DB.Model(&models.KnowledgeDocument{}).Where("is_active = true AND source_type = ?", "professional").Count(&professional)
	database.DB.Model(&models.KnowledgeDocument{}).Where("is_active = true AND source_type = ?", "opera").Count(&opera)
	database.DB.Model(&models.KnowledgeDocument{}).Where("is_active = true AND embedding_status = ?", "completed").Count(&completed)
	database.DB.Model(&models.KnowledgeDocument{}).Where("is_active = true AND embedding_status IN ?", []string{"pending", "processing"}).Count(&pending)
	database.DB.Model(&models.KnowledgeDocument{}).Where("is_active = true AND embedding_status = ?", "failed").Count(&failed)
	database.DB.Model(&models.KnowledgeEmbedding{}).
		Joins("JOIN knowledge_documents ON knowledge_embeddings.doc_id = knowledge_documents.doc_id AND knowledge_documents.is_active = true").
		Count(&totalChunks)

	avgChunks := 0.0
	if total > 0 {
		avgChunks = float64(totalChunks) / float64(total)
	}

	return map[string]interface{}{
		"total_documents":    total,
		"professional_docs":  professional,
		"opera_docs":         opera,
		"completed_docs":     completed,
		"pending_docs":       pending,
		"failed_docs":        failed,
		"total_chunks":       totalChunks,
		"avg_chunks_per_doc": avgChunks,
	}
}

// SoftDeleteDocument 软删除文档（标记 is_active=false）
func SoftDeleteDocument(docID uuid.UUID) error {
	result := database.DB.Model(&models.KnowledgeDocument{}).
		Where("doc_id = ?", docID).
		Update("is_active", false)
	if result.Error != nil {
		return fmt.Errorf("删除文档失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("文档不存在")
	}
	return nil
}

// ImportOperaAsDocument 从 Opera 字幕内容创建知识库文档
func ImportOperaAsDocument(operaTitle, subtitleContent string, operaID uint, adminUserID uint) (*models.KnowledgeDocument, error) {
	title := fmt.Sprintf("%s - 剧情字幕", operaTitle)
	doc, err := UploadDocument(title, subtitleContent, models.SourceTypeOpera, adminUserID)
	if err != nil {
		return nil, err
	}
	// 入队向量化
	EnqueueDocumentEmbedding(doc.DocID)
	return doc, nil
}

// ImportAllOperasToKnowledge 将所有未隐藏且有字幕的作品导入知识库（后台批量任务）
func ImportAllOperasToKnowledge(adminUserID uint) {
	log.Printf("[KnowledgeService] 开始批量导入作品字幕到知识库...")

	type operaRow struct {
		OperaID    uint
		OperaTitle string
		SrtPath    string
	}

	var operas []operaRow
	if err := database.DB.Table("opera").
		Select("opera_id, opera_title, srt_path").
		Where("is_hidden = false AND srt_path IS NOT NULL AND srt_path != ''").
		Scan(&operas).Error; err != nil {
		log.Printf("[KnowledgeService] 查询作品失败: %v", err)
		return
	}

	log.Printf("[KnowledgeService] 共找到 %d 部有字幕的作品", len(operas))

	for _, opera := range operas {
		// 检查是否已导入
		var count int64
		database.DB.Model(&models.KnowledgeDocument{}).
			Where("title = ? AND is_active = true", fmt.Sprintf("%s - 剧情字幕", opera.OperaTitle)).
			Count(&count)
		if count > 0 {
			log.Printf("[KnowledgeService] 作品 %s 已导入，跳过", opera.OperaTitle)
			continue
		}

		// 从 CDN 读取字幕
		subtitleContent, err := fetchSubtitleContent(opera.SrtPath)
		if err != nil {
			log.Printf("[KnowledgeService] 读取字幕失败 (Opera %d): %v", opera.OperaID, err)
			continue
		}

		if _, err := ImportOperaAsDocument(opera.OperaTitle, subtitleContent, opera.OperaID, adminUserID); err != nil {
			log.Printf("[KnowledgeService] 导入作品 %d 失败: %v", opera.OperaID, err)
		}
	}

	log.Printf("[KnowledgeService] 批量导入完成")
}

// =============================================
// 内部工具函数
// =============================================

// splitIntoChunks 将长文本分割为固定大小的块（字符级，支持中文）
func splitIntoChunks(text string, size, overlap int) []string {
	// 将文本转换为 rune 切片（正确处理中文字符）
	runes := []rune(text)
	total := len(runes)
	if total == 0 {
		return nil
	}

	var chunks []string
	step := size - overlap
	if step <= 0 {
		step = size
	}

	for start := 0; start < total; start += step {
		end := start + size
		if end > total {
			end = total
		}
		chunk := string(runes[start:end])
		// 过滤空白块
		if strings.TrimSpace(chunk) != "" {
			chunks = append(chunks, chunk)
		}
		if end >= total {
			break
		}
	}
	return chunks
}

// float32Slice 将 []float64 转为 []float32
func float32Slice(f64 []float64) []float32 {
	f32 := make([]float32, len(f64))
	for i, v := range f64 {
		f32[i] = float32(v)
	}
	return f32
}

// 确保 utf8 包被使用（用于字符计数辅助）
var _ = utf8.RuneCountInString
