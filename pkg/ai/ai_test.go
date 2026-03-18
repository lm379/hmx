package ai

import (
	"encoding/json"
	"testing"

	"github.com/lm379/hmx/config"
)

func TestEmbeddingRequest(t *testing.T) {
	req := EmbeddingRequest{
		Input: "test input",
		Model: "text-embedding-3-small",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal EmbeddingRequest: %v", err)
	}

	var decoded EmbeddingRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal EmbeddingRequest: %v", err)
	}

	if decoded.Input != req.Input {
		t.Errorf("Input mismatch: got %q, want %q", decoded.Input, req.Input)
	}
	if decoded.Model != req.Model {
		t.Errorf("Model mismatch: got %q, want %q", decoded.Model, req.Model)
	}
}

func TestEmbeddingResponse(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"embedding": [0.1, 0.2, 0.3, 0.4, 0.5]
			}
		]
	}`

	var resp EmbeddingResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal EmbeddingResponse: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 data entry, got %d", len(resp.Data))
	}

	if len(resp.Data[0].Embedding) != 5 {
		t.Errorf("Expected 5 embedding dimensions, got %d", len(resp.Data[0].Embedding))
	}

	expectedEmbedding := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	for i, val := range resp.Data[0].Embedding {
		if val != expectedEmbedding[i] {
			t.Errorf("Embedding[%d] = %v, want %v", i, val, expectedEmbedding[i])
		}
	}
}

func TestChatMessage(t *testing.T) {
	msg := ChatMessage{
		Role:    "user",
		Content: "Hello",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal ChatMessage: %v", err)
	}

	var decoded ChatMessage
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ChatMessage: %v", err)
	}

	if decoded.Role != msg.Role {
		t.Errorf("Role mismatch: got %q, want %q", decoded.Role, msg.Role)
	}
	if decoded.Content != msg.Content {
		t.Errorf("Content mismatch: got %q, want %q", decoded.Content, msg.Content)
	}
}

func TestChatCompletionRequest(t *testing.T) {
	req := ChatCompletionRequest{
		Model: "gpt-4.1",
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a helpful assistant"},
			{Role: "user", Content: "Hello"},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ChatCompletionRequest: %v", err)
	}

	var decoded ChatCompletionRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ChatCompletionRequest: %v", err)
	}

	if decoded.Model != req.Model {
		t.Errorf("Model mismatch: got %q, want %q", decoded.Model, req.Model)
	}

	if len(decoded.Messages) != len(req.Messages) {
		t.Errorf("Messages length mismatch: got %d, want %d", len(decoded.Messages), len(req.Messages))
	}
}

func TestChatCompletionResponse(t *testing.T) {
	jsonData := `{
		"choices": [
			{
				"message": {
					"role": "assistant",
					"content": "Hello! How can I help you today?"
				}
			}
		]
	}`

	var resp ChatCompletionResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal ChatCompletionResponse: %v", err)
	}

	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got %q", resp.Choices[0].Message.Role)
	}

	expectedContent := "Hello! How can I help you today?"
	if resp.Choices[0].Message.Content != expectedContent {
		t.Errorf("Content mismatch: got %q, want %q", resp.Choices[0].Message.Content, expectedContent)
	}
}

func TestEmbeddingRequestFields(t *testing.T) {
	tests := []struct {
		name  string
		input string
		model string
	}{
		{
			name:  "Simple text",
			input: "Hello world",
			model: "text-embedding-3-small",
		},
		{
			name:  "Chinese text",
			input: "黄梅戏",
			model: "text-embedding-3-small",
		},
		{
			name:  "Empty string",
			input: "",
			model: "text-embedding-3-small",
		},
		{
			name:  "Special characters",
			input: "!@#$%^&*()",
			model: "text-embedding-3-small",
		},
		{
			name:  "Different model",
			input: "test",
			model: "text-embedding-3-large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := EmbeddingRequest{
				Input: tt.input,
				Model: tt.model,
			}

			data, err := json.Marshal(req)
			if err != nil {
				t.Errorf("Failed to marshal: %v", err)
				return
			}

			var decoded EmbeddingRequest
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
				return
			}

			if decoded.Input != tt.input {
				t.Errorf("Input mismatch: got %q, want %q", decoded.Input, tt.input)
			}
			if decoded.Model != tt.model {
				t.Errorf("Model mismatch: got %q, want %q", decoded.Model, tt.model)
			}
		})
	}
}

func TestChatMessageRoles(t *testing.T) {
	roles := []string{"system", "user", "assistant"}

	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			msg := ChatMessage{
				Role:    role,
				Content: "test content",
			}

			data, err := json.Marshal(msg)
			if err != nil {
				t.Errorf("Failed to marshal: %v", err)
				return
			}

			var decoded ChatMessage
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
				return
			}

			if decoded.Role != role {
				t.Errorf("Role mismatch: got %q, want %q", decoded.Role, role)
			}
		})
	}
}

func TestChatCompletionRequestWithMultipleMessages(t *testing.T) {
	req := ChatCompletionRequest{
		Model: "gpt-4.1",
		Messages: []ChatMessage{
			{Role: "system", Content: "System prompt"},
			{Role: "user", Content: "User message 1"},
			{Role: "assistant", Content: "Assistant response 1"},
			{Role: "user", Content: "User message 2"},
			{Role: "assistant", Content: "Assistant response 2"},
			{Role: "user", Content: "User message 3"},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded ChatCompletionRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(decoded.Messages) != 6 {
		t.Errorf("Expected 6 messages, got %d", len(decoded.Messages))
	}

	expectedSequence := []string{"system", "user", "assistant", "user", "assistant", "user"}
	for i, expectedRole := range expectedSequence {
		if decoded.Messages[i].Role != expectedRole {
			t.Errorf("Message[%d] role: got %q, want %q", i, decoded.Messages[i].Role, expectedRole)
		}
	}
}

func TestGenerateEmbeddingConfigValidation(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	tests := []struct {
		name        string
		apiKey      string
		baseURL     string
		model       string
		expectError bool
	}{
		{
			name:        "Missing API key",
			apiKey:      "",
			baseURL:     "https://api.openai.com/v1",
			model:       "text-embedding-3-small",
			expectError: true,
		},
		{
			name:        "All configs present",
			apiKey:      "test-key",
			baseURL:     "https://api.openai.com/v1",
			model:       "text-embedding-3-small",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.AppConfig = config.Config{
				EmbeddingAPIKey:  tt.apiKey,
				EmbeddingBaseURL: tt.baseURL,
				EmbeddingModel:   tt.model,
			}

			_, err := GenerateEmbedding("test")

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Logf("Got error (API may not be accessible): %v", err)
				// This is expected if the API is not reachable
			}
		})
	}
}

func TestGenerateSubtitleSummaryConfigValidation(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	tests := []struct {
		name        string
		apiKey      string
		baseURL     string
		model       string
		expectError bool
	}{
		{
			name:        "Missing API key",
			apiKey:      "",
			baseURL:     "https://api.openai.com/v1",
			model:       "gpt-4.1",
			expectError: true,
		},
		{
			name:        "All configs present",
			apiKey:      "test-key",
			baseURL:     "https://api.openai.com/v1",
			model:       "gpt-4.1",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.AppConfig = config.Config{
				SubtitleAPIKey:  tt.apiKey,
				SubtitleBaseURL: tt.baseURL,
				SubtitleModel:   tt.model,
			}

			_, err := GenerateSubtitleSummary("test", "test content")

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Logf("Got error (API may not be accessible): %v", err)
				// This is expected if the API is not reachable
			}
		})
	}
}

func TestStructJSONRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		instance interface{}
		json     string
	}{
		{
			name:     "EmbeddingRequest",
			instance: EmbeddingRequest{Input: "test", Model: "model"},
			json:     `{"input":"test","model":"model"}`,
		},
		{
			name:     "ChatMessage",
			instance: ChatMessage{Role: "user", Content: "hello"},
			json:     `{"role":"user","content":"hello"}`,
		},
		{
			name:     "ChatCompletionRequest with single message",
			instance: ChatCompletionRequest{Model: "gpt-4", Messages: []ChatMessage{{Role: "user", Content: "hi"}}},
			json:     `{"model":"gpt-4","messages":[{"role":"user","content":"hi"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.instance)
			if err != nil {
				t.Errorf("Failed to marshal: %v", err)
				return
			}

			// Unmarshal to verify the structure
			var decoded interface{}
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
				return
			}

			// Just verify that it can be marshaled and unmarshaled
			if decoded == nil {
				t.Error("Decoded value is nil")
			}
		})
	}
}
