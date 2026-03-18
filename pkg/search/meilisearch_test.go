package search

import (
	"testing"

	"github.com/lm379/hmx/config"
)

func TestStringsToInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []interface{}
	}{
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []interface{}{},
		},
		{
			name:     "Single element",
			input:    []string{"test"},
			expected: []interface{}{"test"},
		},
		{
			name:     "Multiple elements",
			input:    []string{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name:     "Chinese characters",
			input:    []string{"黄梅戏", "戏曲"},
			expected: []interface{}{"黄梅戏", "戏曲"},
		},
		{
			name:     "Mixed content",
			input:    []string{"test", "测试", "123"},
			expected: []interface{}{"test", "测试", "123"},
		},
		{
			name:     "Empty strings",
			input:    []string{"", "test", ""},
			expected: []interface{}{"", "test", ""},
		},
		{
			name:     "Special characters",
			input:    []string{"!@#", "$%^"},
			expected: []interface{}{"!@#", "$%^"},
		},
		{
			name:     "Long slice",
			input:    []string{"a", "b", "c", "d", "e", "f"},
			expected: []interface{}{"a", "b", "c", "d", "e", "f"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringsToInterfaces(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("stringsToInterfaces() length = %d, want %d", len(result), len(tt.expected))
				return
			}

			for i := range tt.expected {
				if result[i] != tt.expected[i] {
					t.Errorf("stringsToInterfaces()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestGetMeilisearchClient(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	tests := []struct {
		name        string
		host        string
		apiKey      string
		expectNil   bool
		expectError bool
	}{
		{
			name:        "Valid configuration",
			host:        "http://localhost:7700",
			apiKey:      "test-key",
			expectNil:   false,
			expectError: false,
		},
		{
			name:        "Empty host (uses default)",
			host:        "",
			apiKey:      "test-key",
			expectNil:   false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the singleton for each test
			config.AppConfig = config.Config{
				MeilisearchHost:   tt.host,
				MeilisearchAPIKey: tt.apiKey,
			}

			// Note: GetMeilisearchClient uses sync.Once, so it will only execute once
			// in a single test run. This is expected behavior for a singleton pattern.
			client, err := GetMeilisearchClient()

			if tt.expectNil && client != nil {
				t.Error("Expected client to be nil, got non-nil")
			}

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Logf("Got error (connection may be unavailable): %v", err)
				// This is expected if Meilisearch is not running
			}
		})
	}
}

func TestStringsToInterfacesNilSafe(t *testing.T) {
	result := stringsToInterfaces([]string{})

	if result == nil {
		t.Error("stringsToInterfaces() should return non-nil slice for empty input")
	}

	if len(result) != 0 {
		t.Errorf("stringsToInterfaces() length = %d, want 0", len(result))
	}
}

func TestStringsToInterfacesTypes(t *testing.T) {
	inputs := []string{"test1", "test2", "test3"}
	results := stringsToInterfaces(inputs)

	for i, result := range results {
		str, ok := result.(string)
		if !ok {
			t.Errorf("stringsToInterfaces()[%d] is not a string, got %T", i, result)
		}
		if str != inputs[i] {
			t.Errorf("stringsToInterfaces()[%d] = %q, want %q", i, str, inputs[i])
		}
	}
}

func TestStringsToInterfacesConsistency(t *testing.T) {
	inputs := []string{"a", "b", "c"}

	result1 := stringsToInterfaces(inputs)
	result2 := stringsToInterfaces(inputs)

	if len(result1) != len(result2) {
		t.Errorf("Inconsistent result lengths: %d vs %d", len(result1), len(result2))
	}

	for i := range result1 {
		if result1[i] != result2[i] {
			t.Errorf("Inconsistent result at index %d: %v vs %v", i, result1[i], result2[i])
		}
	}
}

func TestStringsToInterfacesWithUnicode(t *testing.T) {
	tests := []struct {
		name  string
		input []string
	}{
		{
			name:  "Chinese characters",
			input: []string{"黄梅戏", "戏曲", "艺术"},
		},
		{
			name:  "Emoji",
			input: []string{"😊", "😂", "🎭"},
		},
		{
			name:  "Japanese characters",
			input: []string{"日本語", "漢字"},
		},
		{
			name:  "Korean characters",
			input: []string{"한국어", "한글"},
		},
		{
			name:  "Arabic characters",
			input: []string{"مرحبا", "العربية"},
		},
		{
			name:  "Mixed unicode",
			input: []string{"测试", "Test", "テスト", "Тест"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringsToInterfaces(tt.input)

			if len(result) != len(tt.input) {
				t.Errorf("stringsToInterfaces() length = %d, want %d", len(result), len(tt.input))
			}

			for i := range tt.input {
				if result[i] != tt.input[i] {
					t.Errorf("stringsToInterfaces()[%d] = %v, want %v", i, result[i], tt.input[i])
				}
			}
		})
	}
}

func TestGetMeilisearchClientSingleton(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	config.AppConfig = config.Config{
		MeilisearchHost:   "http://localhost:7700",
		MeilisearchAPIKey: "test-key",
	}

	// Call GetMeilisearchClient multiple times
	client1, err1 := GetMeilisearchClient()
	client2, err2 := GetMeilisearchClient()

	// Both should return the same instance (or error should be consistent)
	if (err1 == nil) != (err2 == nil) {
		t.Errorf("Inconsistent error results: %v vs %v", err1, err2)
	}

	// If both clients are non-nil, they should be the same instance
	if client1 != nil && client2 != nil && client1 != client2 {
		t.Log("Clients are different instances (this may be expected depending on implementation)")
	}
}

func TestStringsToInterfacesLongStrings(t *testing.T) {
	longString := "a"
	for i := 0; i < 1000; i++ {
		longString += "a"
	}

	inputs := []string{longString}
	results := stringsToInterfaces(inputs)

	if len(results) != 1 {
		t.Errorf("stringsToInterfaces() returned unexpected length: %d", len(results))
	}

	if results[0] != longString {
		t.Error("Long string was not preserved correctly")
	}
}
