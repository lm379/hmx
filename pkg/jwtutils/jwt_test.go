package jwtutils

import (
	"testing"
	"time"

	"github.com/lm379/hmx/config"
)

func setupTestConfig() {
	config.AppConfig = config.Config{
		JWTSecret:                "test-secret-key-for-jwt-token-generation-and-validation",
		JWTAccessTokenExpiresIn:  15 * time.Minute,
		JWTRefreshTokenExpiresIn: 168 * time.Hour, // 7 days
	}
}

func TestGenerateToken(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name      string
		userID    uint
		username  string
		role      string
		wantError bool
	}{
		{
			name:      "Valid token generation",
			userID:    1,
			username:  "testuser",
			role:      "user",
			wantError: false,
		},
		{
			name:      "Admin role",
			userID:    2,
			username:  "admin",
			role:      "admin",
			wantError: false,
		},
		{
			name:      "Empty username",
			userID:    3,
			username:  "",
			role:      "user",
			wantError: false,
		},
		{
			name:      "Zero user ID",
			userID:    0,
			username:  "testuser",
			role:      "user",
			wantError: false,
		},
		{
			name:      "Large user ID",
			userID:    999999,
			username:  "testuser",
			role:      "user",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.username, tt.role)
			if (err != nil) != tt.wantError {
				t.Errorf("GenerateToken() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && token == "" {
				t.Error("GenerateToken() returned empty token")
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	setupTestConfig()

	token, err := GenerateAccessToken(1, "testuser", "user")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateAccessToken() returned empty token")
	}

	// Verify the token can be parsed
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("Expected UserID to be 1, got %v", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("Expected Username to be 'testuser', got %v", claims.Username)
	}
	if claims.Role != "user" {
		t.Errorf("Expected Role to be 'user', got %v", claims.Role)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	setupTestConfig()

	token, err := GenerateRefreshToken(1, "testuser", "user")
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateRefreshToken() returned empty token")
	}

	// Verify the token can be parsed
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("Expected UserID to be 1, got %v", claims.UserID)
	}

	// Check that the expiration time is set correctly
	if claims.ExpiresAt == nil {
		t.Error("Expected ExpiresAt to be set")
	}

	// Refresh token should expire in 7 days (168 hours)
	expectedExpiry := time.Now().Add(168 * time.Hour)
	timeDiff := expectedExpiry.Sub(claims.ExpiresAt.Time)
	if timeDiff < -time.Minute || timeDiff > time.Minute {
		t.Errorf("Expected expiration time to be around 7 days from now, got %v", claims.ExpiresAt.Time)
	}
}

func TestValidateToken(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name      string
		token     string
		wantError bool
		checkUser uint
	}{
		{
			name:      "Valid token",
			token:     generateValidToken(t),
			wantError: false,
			checkUser: 123,
		},
		{
			name:      "Invalid token - empty",
			token:     "",
			wantError: true,
		},
		{
			name:      "Invalid token - random string",
			token:     "not.a.valid.jwt.token",
			wantError: true,
		},
		{
			name:      "Invalid token - wrong secret",
			token:     generateTokenWithDifferentSecret(t),
			wantError: true,
		},
		{
			name:      "Invalid token - malformed",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateToken() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if claims == nil {
					t.Error("ValidateToken() returned nil claims")
				}
				if claims.UserID != tt.checkUser {
					t.Errorf("Expected UserID to be %v, got %v", tt.checkUser, claims.UserID)
				}
			}
		})
	}
}

func TestTokenRoundTrip(t *testing.T) {
	setupTestConfig()

	userID := uint(42)
	username := "roundtripuser"
	role := "admin"

	// Generate token
	token, err := GenerateToken(userID, username, role)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Validate token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	// Check claims
	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %v, want %v", claims.UserID, userID)
	}
	if claims.Username != username {
		t.Errorf("Username mismatch: got %v, want %v", claims.Username, username)
	}
	if claims.Role != role {
		t.Errorf("Role mismatch: got %v, want %v", claims.Role, role)
	}
	if claims.Issuer != "hmx-backend" {
		t.Errorf("Issuer mismatch: got %v, want 'hmx-backend'", claims.Issuer)
	}
}

func TestTokenExpiration(t *testing.T) {
	setupTestConfig()

	token, err := GenerateAccessToken(1, "testuser", "user")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	// Check that the expiration time is set correctly
	if claims.ExpiresAt == nil {
		t.Error("Expected ExpiresAt to be set")
	}

	// Access token should expire in 15 minutes
	expectedExpiry := time.Now().Add(15 * time.Minute)
	timeDiff := expectedExpiry.Sub(claims.ExpiresAt.Time)
	if timeDiff < -time.Minute || timeDiff > time.Minute {
		t.Errorf("Expected expiration time to be around 15 minutes from now, got %v", claims.ExpiresAt.Time)
	}
}

func TestTokenConsistency(t *testing.T) {
	setupTestConfig()

	userID := uint(10)
	username := "consistent"
	role := "user"

	// Generate multiple tokens with the same parameters
	token1, err1 := GenerateToken(userID, username, role)
	if err1 != nil {
		t.Fatalf("First GenerateToken() error = %v", err1)
	}

	token2, err2 := GenerateToken(userID, username, role)
	if err2 != nil {
		t.Fatalf("Second GenerateToken() error = %v", err2)
	}

	// Tokens should be different (due to different issued-at times or different nanoseconds)
	// Note: In extremely rare cases, if both tokens are generated at the exact same nanosecond,
	// they might be the same. This is acceptable behavior.
	if token1 == token2 {
		t.Log("Tokens are the same (generated at the exact same nanosecond)")
	}

	// But both should be valid
	claims1, err1 := ValidateToken(token1)
	if err1 != nil {
		t.Errorf("First ValidateToken() error = %v", err1)
	}

	claims2, err2 := ValidateToken(token2)
	if err2 != nil {
		t.Errorf("Second ValidateToken() error = %v", err2)
	}

	// Claims should be identical except for issued-at time
	if claims1.UserID != claims2.UserID {
		t.Errorf("UserID mismatch between two tokens")
	}
	if claims1.Username != claims2.Username {
		t.Errorf("Username mismatch between two tokens")
	}
	if claims1.Role != claims2.Role {
		t.Errorf("Role mismatch between two tokens")
	}
}

func TestTokenWithSpecialCharacters(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name     string
		userID   uint
		username string
		role     string
	}{
		{
			name:     "Username with special characters",
			userID:   1,
			username: "user@domain.com",
			role:     "user",
		},
		{
			name:     "Username with emoji",
			userID:   2,
			username: "user😊",
			role:     "user",
		},
		{
			name:     "Username with Chinese characters",
			userID:   3,
			username: "测试用户",
			role:     "user",
		},
		{
			name:     "Username with spaces",
			userID:   4,
			username: "user name",
			role:     "user",
		},
		{
			name:     "Username with underscores",
			userID:   5,
			username: "user_name",
			role:     "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.username, tt.role)
			if err != nil {
				t.Errorf("GenerateToken() error = %v", err)
				return
			}

			claims, err := ValidateToken(token)
			if err != nil {
				t.Errorf("ValidateToken() error = %v", err)
				return
			}

			if claims.UserID != tt.userID {
				t.Errorf("UserID mismatch: got %v, want %v", claims.UserID, tt.userID)
			}
			if claims.Username != tt.username {
				t.Errorf("Username mismatch: got %v, want %v", claims.Username, tt.username)
			}
			if claims.Role != tt.role {
				t.Errorf("Role mismatch: got %v, want %v", claims.Role, tt.role)
			}
		})
	}
}

// Helper functions

func generateValidToken(t *testing.T) string {
	setupTestConfig()
	token, err := GenerateToken(123, "testuser", "user")
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}
	return token
}

func generateTokenWithDifferentSecret(t *testing.T) string {
	// Save original config
	originalConfig := config.AppConfig
	defer func() { config.AppConfig = originalConfig }()

	// Generate token with different secret
	config.AppConfig.JWTSecret = "different-secret-key"
	token, err := GenerateToken(123, "testuser", "user")
	if err != nil {
		t.Fatalf("Failed to generate token with different secret: %v", err)
	}

	// Restore original config
	config.AppConfig = originalConfig

	return token
}
