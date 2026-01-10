package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lm379/hmx/config"
)

// OpenAI Embedding API 请求和响应结构
type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
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
	} `json:"choices"`
}

// GenerateEmbedding 调用Embedding API生成向量
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
	// 使用聊天模型进行总结
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
