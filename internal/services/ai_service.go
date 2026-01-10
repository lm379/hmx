package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
)

// OpenAI Embedding API 请求和响应结构
type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Chat Completion API 请求和响应结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
		Index   int         `json:"index"`
	} `json:"choices"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// GenerateEmbedding 使用配置的Embedding服务生成文本向量
func GenerateEmbedding(text string) ([]float64, error) {
	apiKey := config.AppConfig.EmbeddingAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("embedding API key is not configured")
	}

	baseURL := config.AppConfig.EmbeddingBaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := config.AppConfig.EmbeddingModel
	if model == "" {
		model = "text-embedding-3-small"
	}

	reqBody := EmbeddingRequest{
		Input: text,
		Model: model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := baseURL + "/embeddings"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var embResp EmbeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(embResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}

	return embResp.Data[0].Embedding, nil
}

// GenerateSubtitleSummary 使用LLM对字幕内容生成摘要
func GenerateSubtitleSummary(operaTitle string, subtitleContent string) (string, error) {
	apiKey := config.AppConfig.SubtitleAPIKey
	if apiKey == "" {
		return "", fmt.Errorf("subtitle API key is not configured")
	}

	baseURL := config.AppConfig.SubtitleBaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := config.AppConfig.SubtitleModel
	// 使用聊天模型进行总结（默认使用gpt-4.1或配置的模型）
	if model == "" {
		model = "gpt-4.1"
	}

	prompt := fmt.Sprintf("这是黄梅戏%s的视频字幕，请根据该视频的字幕内容，生成一段简洁的视频摘要（200字以内）：\n\n%s", operaTitle, subtitleContent)

	reqBody := ChatCompletionRequest{
		Model: model,
		Messages: []ChatMessage{
			{Role: "system", Content: "你是一个专业的黄梅戏视频内容评鉴专家，擅长根据字幕提炼视频核心内容并生成简洁摘要。"},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := baseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// GenerateOperaEmbedding 为作品生成向量（异步）
func GenerateOperaEmbedding(operaID uint) {
	log.Printf("[AI Embedding] 开始为 Opera ID %d 生成向量...", operaID)

	operaRepo := repository.NewOperaRepo()
	opera, err := operaRepo.GetByID(operaID)
	if err != nil {
		log.Printf("[AI Embedding] 错误: 找不到 Opera ID %d: %v", operaID, err)
		return
	}

	// 构建文本：标题 + 描述 + AI摘要
	text := opera.OperaTitle
	if opera.Description != "" {
		text += " " + opera.Description
	}
	if opera.AiSummary != "" {
		text += " " + opera.AiSummary
	}

	// 调用 Embedding API 生成向量
	embedding, err := GenerateEmbedding(text)
	if err != nil {
		log.Printf("[AI Embedding] 错误: 生成向量失败 (Opera ID %d): %v", operaID, err)
		return
	}

	// 直接存储向量到 pgvector 字段
	if err := operaRepo.UpdateEmbedding(opera.OperaID, embedding); err != nil {
		log.Printf("[AI Embedding] 错误: 更新向量失败 (Opera ID %d): %v", operaID, err)
		return
	}

	log.Printf("[AI Embedding] 成功: Opera ID %d 向量已生成并保存 (维度: %d)", operaID, len(embedding))
}

// GetVideoSubtitleSummary 获取视频的AI字幕摘要
// 逻辑：
// 1. 如果数据库中已有AI摘要，直接返回
// 2. 如果没有AI摘要，先返回"正在生成中"
// 3. 检查是否存在AI字幕文件，如果存在则异步生成摘要并存入数据库
func GetVideoSubtitleSummary(operaID uint) (string, error) {
	operaRepo := repository.NewOperaRepo()
	opera, err := operaRepo.GetByID(operaID)
	if err != nil {
		return "", fmt.Errorf("找不到作品: %w", err)
	}

	// 如果已有AI摘要，直接返回
	if opera.AiSummary != "" {
		return opera.AiSummary, nil
	}

	// 没有AI摘要，检查是否有字幕文件
	if opera.SrtPath.Valid && opera.SrtPath.String != "" {
		// 有字幕文件，异步生成摘要
		go generateSummaryFromSubtitle(operaID, opera.SrtPath.String)
	}
	// 无论是否有字幕，都先返回"正在生成中"
	// 如果没有字幕，用户下次请求时依然会得到"正在生成中"

	return "正在生成中", nil
}

// generateSummaryFromSubtitle 从字幕文件生成视频摘要（异步）
func generateSummaryFromSubtitle(operaID uint, srtPath string) {
	log.Printf("[AI Subtitle Summary] 开始为 Opera ID %d 生成字幕摘要...", operaID)

	// 获取作品信息（需要标题）
	operaRepo := repository.NewOperaRepo()
	opera, err := operaRepo.GetByID(operaID)
	if err != nil {
		log.Printf("[AI Subtitle Summary] 错误: 获取作品信息失败 (Opera ID %d): %v", operaID, err)
		return
	}

	// 从CDN/对象存储读取字幕文件内容
	subtitleContent, err := readSubtitleFromCDN(srtPath)
	if err != nil {
		log.Printf("[AI Subtitle Summary] 错误: 读取字幕失败 (Opera ID %d): %v", operaID, err)
		return
	}

	// 使用LLM生成摘要
	summary, err := GenerateSubtitleSummary(opera.OperaTitle, subtitleContent)
	if err != nil {
		log.Printf("[AI Subtitle Summary] 错误: 生成摘要失败 (Opera ID %d): %v", operaID, err)
		return
	}

	// 保存到数据库
	if err := operaRepo.UpdateSummary(operaID, summary); err != nil {
		log.Printf("[AI Subtitle Summary] 错误: 更新摘要失败 (Opera ID %d): %v", operaID, err)
		return
	}

	log.Printf("[AI Subtitle Summary] 成功: Opera ID %d 摘要已生成", operaID)

	// 生成摘要成功后，自动触发向量生成
	go GenerateOperaEmbedding(operaID)
}

// GenerateOperaSummarySync 同步生成作品的AI摘要（用于管理后台）
func GenerateOperaSummarySync(operaID uint, srtPath string) (string, error) {
	log.Printf("[AI Subtitle Summary Sync] 开始为 Opera ID %d 生成字幕摘要...", operaID)

	// 获取作品信息（需要标题）
	operaRepo := repository.NewOperaRepo()
	opera, err := operaRepo.GetByID(operaID)
	if err != nil {
		return "", fmt.Errorf("获取作品信息失败: %w", err)
	}

	// 从CDN/对象存储读取字幕文件内容
	subtitleContent, err := readSubtitleFromCDN(srtPath)
	if err != nil {
		return "", fmt.Errorf("读取字幕失败: %w", err)
	}

	// 使用LLM生成摘要
	summary, err := GenerateSubtitleSummary(opera.OperaTitle, subtitleContent)
	if err != nil {
		return "", fmt.Errorf("生成摘要失败: %w", err)
	}

	// 保存到数据库
	if err := operaRepo.UpdateSummary(operaID, summary); err != nil {
		return "", fmt.Errorf("更新摘要失败: %w", err)
	}

	log.Printf("[AI Subtitle Summary Sync] 成功: Opera ID %d 摘要已生成", operaID)
	return summary, nil
}

// readSubtitleFromCDN 从CDN/对象存储读取字幕文件内容
func readSubtitleFromCDN(srtPath string) (string, error) {
	// 构建完整的CDN URL
	cdnURL := config.AppConfig.S3CustomDomain
	if cdnURL == "" {
		// 如果没有自定义域名，使用S3 endpoint
		cdnURL = config.AppConfig.S3Endpoint + "/" + config.AppConfig.S3Bucket
	}

	// 移除srtPath开头的斜杠（如果有）
	if len(srtPath) > 0 && srtPath[0] == '/' {
		srtPath = srtPath[1:]
	}

	fullURL := cdnURL + "/" + srtPath

	// 发起HTTP GET请求读取字幕文件
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch subtitle file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch subtitle, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read subtitle content: %w", err)
	}

	return string(body), nil
}

// BatchGenerateOperaSummariesAsync 批量生成作品AI摘要（异步执行）
func BatchGenerateOperaSummariesAsync(operaIDs []uint) int {
	if len(operaIDs) == 0 {
		return 0
	}

	operaRepo := repository.NewOperaRepo()

	// 获取指定的作品
	operas := make([]models.Opera, 0)
	for _, id := range operaIDs {
		opera, err := operaRepo.GetByID(id)
		if err != nil {
			continue
		}
		// 过滤没有字幕的作品
		if opera.SrtPath.Valid && opera.SrtPath.String != "" {
			operas = append(operas, *opera)
		}
	}

	total := len(operas)
	if total == 0 {
		return 0
	}

	log.Printf("[Batch Summary] 开始批量生成摘要，共 %d 个作品", total)

	// 异步执行批量生成
	go func() {
		var success, failed int
		for _, opera := range operas {
			summary, err := GenerateOperaSummarySync(opera.OperaID, opera.SrtPath.String)
			if err != nil {
				failed++
				log.Printf("[Batch Summary] 失败: Opera ID %d (%s): %v", opera.OperaID, opera.OperaTitle, err)
				continue
			}

			// 成功生成摘要，不自动触发向量生成
			_ = summary
			success++
			log.Printf("[Batch Summary] 成功: Opera ID %d (%s)", opera.OperaID, opera.OperaTitle)
		}
		log.Printf("[Batch Summary] 批量生成完成: 成功 %d, 失败 %d", success, failed)
	}()

	return total
}
