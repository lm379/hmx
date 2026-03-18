package urlutils

import (
	"testing"

	"github.com/lm379/hmx/config"
)

func TestGetFullURL(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	tests := []struct {
		name           string
		path           string
		s3CustomDomain string
		s3Endpoint     string
		s3Bucket       string
		expected       string
	}{
		{
			name:           "Empty path",
			path:           "",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "",
		},
		{
			name:           "Already full URL with https",
			path:           "https://cdn.example.com/video.mp4",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://cdn.example.com/video.mp4",
		},
		{
			name:           "Already full URL with http",
			path:           "http://cdn.example.com/video.mp4",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "http://cdn.example.com/video.mp4",
		},
		{
			name:           "Relative path with custom domain",
			path:           "videos/test.mp4",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://cdn.example.com/videos/test.mp4",
		},
		{
			name:           "Relative path with trailing slash in custom domain",
			path:           "videos/test.mp4",
			s3CustomDomain: "https://cdn.example.com/",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://cdn.example.com/videos/test.mp4",
		},
		{
			name:           "Relative path with leading slash",
			path:           "/videos/test.mp4",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://cdn.example.com/videos/test.mp4",
		},
		{
			name:           "Relative path without custom domain",
			path:           "videos/test.mp4",
			s3CustomDomain: "",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://s3.amazonaws.com/my-bucket/videos/test.mp4",
		},
		{
			name:           "Relative path without custom domain, with trailing slash endpoint",
			path:           "videos/test.mp4",
			s3CustomDomain: "",
			s3Endpoint:     "https://s3.amazonaws.com/",
			s3Bucket:       "my-bucket",
			expected:       "https://s3.amazonaws.com/my-bucket/videos/test.mp4",
		},
		{
			name:           "Relative path without custom domain, with leading slash",
			path:           "/videos/test.mp4",
			s3CustomDomain: "",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://s3.amazonaws.com/my-bucket/videos/test.mp4",
		},
		{
			name:           "Root path with custom domain",
			path:           "/",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://cdn.example.com/",
		},
		{
			name:           "Simple filename with custom domain",
			path:           "test.mp4",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://cdn.example.com/test.mp4",
		},
		{
			name:           "Path with multiple slashes",
			path:           "//videos//test.mp4",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://cdn.example.com//videos//test.mp4",
		},
		{
			name:           "Path with subdirectories",
			path:           "videos/2024/01/test.mp4",
			s3CustomDomain: "https://cdn.example.com",
			s3Endpoint:     "https://s3.amazonaws.com",
			s3Bucket:       "my-bucket",
			expected:       "https://cdn.example.com/videos/2024/01/test.mp4",
		},
		{
			name:           "Empty custom domain uses endpoint",
			path:           "avatars/user123.jpg",
			s3CustomDomain: "",
			s3Endpoint:     "https://s3.ap-southeast-1.amazonaws.com",
			s3Bucket:       "hmx-bucket",
			expected:       "https://s3.ap-southeast-1.amazonaws.com/hmx-bucket/avatars/user123.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config for this test
			config.AppConfig = config.Config{
				S3CustomDomain: tt.s3CustomDomain,
				S3Endpoint:     tt.s3Endpoint,
				S3Bucket:       tt.s3Bucket,
			}

			got := GetFullURL(tt.path)
			if got != tt.expected {
				t.Errorf("GetFullURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetFullURLWithConfig(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	// Test priority: custom domain > endpoint
	config.AppConfig = config.Config{
		S3CustomDomain: "https://cdn.example.com",
		S3Endpoint:     "https://s3.amazonaws.com",
		S3Bucket:       "my-bucket",
	}

	result := GetFullURL("videos/test.mp4")
	expected := "https://cdn.example.com/videos/test.mp4"
	if result != expected {
		t.Errorf("Expected custom domain to be prioritized, got %v, want %v", result, expected)
	}

	// Test with empty custom domain
	config.AppConfig.S3CustomDomain = ""
	result = GetFullURL("videos/test.mp4")
	expected = "https://s3.amazonaws.com/my-bucket/videos/test.mp4"
	if result != expected {
		t.Errorf("Expected endpoint to be used when custom domain is empty, got %v, want %v", result, expected)
	}
}

func TestGetFullURLEmptyString(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	config.AppConfig = config.Config{
		S3CustomDomain: "https://cdn.example.com",
		S3Endpoint:     "https://s3.amazonaws.com",
		S3Bucket:       "my-bucket",
	}

	result := GetFullURL("")
	if result != "" {
		t.Errorf("GetFullURL() with empty string should return empty, got %v", result)
	}
}

func TestGetFullURLConsistentBehavior(t *testing.T) {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	config.AppConfig = config.Config{
		S3CustomDomain: "https://cdn.example.com",
		S3Endpoint:     "https://s3.amazonaws.com",
		S3Bucket:       "my-bucket",
	}

	// Test that same path produces same URL
	path := "videos/test.mp4"
	result1 := GetFullURL(path)
	result2 := GetFullURL(path)

	if result1 != result2 {
		t.Errorf("GetFullURL() should produce consistent results, got %v and %v", result1, result2)
	}
}
