package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
)

// GenerateAnswer 调用 LLM 基于检索到的上下文生成回答
func GenerateAnswer(query string, chunks []RetrievedChunk, operaContext string) (string, error) {
	apiKey := config.AppConfig.SubtitleAPIKey
	if apiKey == "" {
		return "", fmt.Errorf("LLM API key 未配置")
	}

	baseURL := config.AppConfig.SubtitleBaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := config.AppConfig.SubtitleModel
	if model == "" {
		model = "gpt-4.1"
	}

	// 构建上下文文本
	var contextParts []string
	for i, chunk := range chunks {
		contextParts = append(contextParts,
			fmt.Sprintf("[参考资料 %d] 来源：%s\n%s", i+1, chunk.DocTitle, chunk.ChunkText))
	}
	contextText := strings.Join(contextParts, "\n\n---\n\n")

	// 构建提示词
	userPrompt := fmt.Sprintf("请根据以下参考资料回答问题。如果参考资料中没有相关信息，请如实说明。\n\n参考资料：\n%s", contextText)
	if operaContext != "" {
		userPrompt += fmt.Sprintf("\n\n当前作品上下文：%s", operaContext)
	}
	userPrompt += fmt.Sprintf("\n\n用户问题：%s", query)

	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "你是一个专业的黄梅戏知识助手，能够根据提供的参考资料准确回答关于黄梅戏的各种问题。请用简洁、专业的中文回答。如果无法从参考资料中找到答案，请诚实说明，不要编造信息。",
			},
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	url := baseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API 返回状态码 %d: %s", resp.StatusCode, string(body))
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("解析 LLM 响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("LLM 返回空响应")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// CreateQAHistory 创建问答历史记录
func CreateQAHistory(userID *uuid.UUID, sessionID string, question, answer string, relatedDocIDs []uuid.UUID, relatedOperaIDs []uint) (*models.QAHistory, error) {
	qa := &models.QAHistory{
		QAID:            uuid.New(),
		UserID:          userID,
		SessionID:       sessionID,
		Question:        question,
		Answer:          answer,
		RelatedDocIDs:   models.UUIDArray(relatedDocIDs),
		RelatedOperaIDs: models.Uint32Array(relatedOperaIDs),
	}
	if err := database.DB.Create(qa).Error; err != nil {
		return nil, fmt.Errorf("保存问答历史失败: %w", err)
	}
	return qa, nil
}

// GetQAHistory 获取用户或 session 的问答历史（分页）
func GetQAHistory(userID *uuid.UUID, sessionID string, page, pageSize int) ([]models.QAHistoryItem, int64, error) {
	var records []models.QAHistory
	var total int64

	db := database.DB.Model(&models.QAHistory{})
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	} else if sessionID != "" {
		db = db.Where("session_id = ?", sessionID)
	} else {
		return nil, 0, fmt.Errorf("需要提供 user_id 或 session_id")
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("获取问答历史失败: %w", err)
	}

	items := make([]models.QAHistoryItem, len(records))
	for i, r := range records {
		// 将 UUIDArray 转换为 []string
		docIDs := make([]string, len(r.RelatedDocIDs))
		for j, id := range r.RelatedDocIDs {
			docIDs[j] = id.String()
		}
		items[i] = models.QAHistoryItem{
			QAID:            r.QAID,
			Question:        r.Question,
			Answer:          r.Answer,
			RelatedDocIDs:   docIDs,
			RelatedOperaIDs: []uint(r.RelatedOperaIDs),
			CreatedAt:       r.CreatedAt,
		}
	}
	return items, total, nil
}

// SubmitFeedback 提交问答反馈
func SubmitFeedback(qaID uuid.UUID, userID *uuid.UUID, rating int, comments string) error {
	// 验证 QA 记录是否存在
	var qa models.QAHistory
	if err := database.DB.Where("qa_id = ?", qaID).First(&qa).Error; err != nil {
		return fmt.Errorf("问答记录不存在: %w", err)
	}

	feedback := &models.QAFeedback{
		FeedbackID: uuid.New(),
		QAID:       qaID,
		UserID:     userID,
		Rating:     rating,
		Comments:   comments,
	}
	if err := database.DB.Create(feedback).Error; err != nil {
		return fmt.Errorf("保存反馈失败: %w", err)
	}
	return nil
}

// AskQuestion RAG 核心流程：检索 → 生成 → 存储
func AskQuestion(query string, userID *uuid.UUID, sessionID string, operaID *uint) (*models.AskResponse, error) {
	log.Printf("[RAG] 开始处理问题: %q (userID=%v, sessionID=%s)", query, userID, sessionID)

	// 1. 向量检索相关文档块
	chunks, err := RetrieveRelevantChunks(query, topK)
	if err != nil {
		return nil, fmt.Errorf("检索失败: %w", err)
	}
	log.Printf("[RAG] 检索到 %d 个相关块", len(chunks))

	// 2. 获取 Opera 上下文（可选）
	operaContext := ""
	var relatedOperaIDs []uint
	if operaID != nil {
		operaRepo, err := getOperaTitle(*operaID)
		if err == nil {
			operaContext = operaRepo
			relatedOperaIDs = append(relatedOperaIDs, *operaID)
		}
	}

	// 3. 调用 LLM 生成回答
	answer, err := GenerateAnswer(query, chunks, operaContext)
	if err != nil {
		return nil, fmt.Errorf("生成回答失败: %w", err)
	}

	// 4. 收集来源信息
	relatedDocIDs := make([]uuid.UUID, 0, len(chunks))
	seenDocs := make(map[uuid.UUID]bool)
	sources := make([]models.ChunkSource, len(chunks))
	for i, chunk := range chunks {
		sources[i] = models.ChunkSource{
			ChunkID:   chunk.ChunkID,
			DocID:     chunk.DocID,
			DocTitle:  chunk.DocTitle,
			ChunkText: chunk.ChunkText,
			Score:     chunk.Score,
		}
		if !seenDocs[chunk.DocID] {
			relatedDocIDs = append(relatedDocIDs, chunk.DocID)
			seenDocs[chunk.DocID] = true
		}
	}

	// 5. 保存历史记录
	qa, err := CreateQAHistory(userID, sessionID, query, answer, relatedDocIDs, relatedOperaIDs)
	if err != nil {
		log.Printf("[RAG] 保存历史失败（不影响返回）: %v", err)
	}

	response := &models.AskResponse{
		Answer:        answer,
		Sources:       sources,
		RelatedOperas: relatedOperaIDs,
	}
	if qa != nil {
		response.QAID = qa.QAID
	}

	log.Printf("[RAG] 问答完成，QA ID: %s", response.QAID)
	return response, nil
}

// getOperaTitle 获取作品标题（供 RAG 上下文使用）
func getOperaTitle(operaID uint) (string, error) {
	var opera struct {
		OperaTitle string
	}
	if err := database.DB.Table("opera").Select("opera_title").Where("opera_id = ?", operaID).Scan(&opera).Error; err != nil {
		return "", err
	}
	return opera.OperaTitle, nil
}
