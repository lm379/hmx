package hashutils

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		{
			name:      "Valid password",
			password:  "password123",
			wantError: false,
		},
		{
			name:      "Empty password",
			password:  "",
			wantError: false,
		},
		{
			name:      "Long password",
			password:  strings.Repeat("a", 50),
			wantError: false,
		},
		{
			name:      "Special characters",
			password:  "!@#$%^&*()_+-=[]{}|;:,.<>?",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantError {
				t.Errorf("HashPassword() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if hash == "" {
					t.Errorf("HashPassword() returned empty hash")
				}
				if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
					t.Errorf("HashPassword() returned invalid bcrypt hash format")
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	validPassword := "test123"
	invalidPassword := "wrong123"

	hash, err := HashPassword(validPassword)
	if err != nil {
		t.Fatalf("Failed to create test hash: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password",
			password: validPassword,
			hash:     hash,
			want:     true,
		},
		{
			name:     "Incorrect password",
			password: invalidPassword,
			hash:     hash,
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Invalid hash",
			password: validPassword,
			hash:     "invalid_hash",
			want:     false,
		},
		{
			name:     "Empty hash",
			password: validPassword,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPasswordHash(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPasswordAndCheck(t *testing.T) {
	password := "mySecurePassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash() should return true for correct password")
	}

	if CheckPasswordHash("wrongPassword", hash) {
		t.Error("CheckPasswordHash() should return false for wrong password")
	}
}

func TestBcryptCost(t *testing.T) {
	// Test that the default cost is set correctly
	if BcryptCost != 12 {
		t.Errorf("Expected BcryptCost to be 12, got %d", BcryptCost)
	}

	// Test that the hash cost matches the expected cost
	password := "testPassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("bcrypt.Cost() error = %v", err)
	}

	if cost != BcryptCost {
		t.Errorf("Expected hash cost to be %d, got %d", BcryptCost, cost)
	}
}

func TestHashPasswordIdempotent(t *testing.T) {
	password := "samePassword"

	// Hash the same password multiple times
	hash1, err1 := HashPassword(password)
	if err1 != nil {
		t.Fatalf("First HashPassword() error = %v", err1)
	}

	hash2, err2 := HashPassword(password)
	if err2 != nil {
		t.Fatalf("Second HashPassword() error = %v", err2)
	}

	// Hashes should be different (due to salt)
	if hash1 == hash2 {
		t.Error("Hashing the same password should produce different hashes (due to salt)")
	}

	// But both should verify correctly
	if !CheckPasswordHash(password, hash1) {
		t.Error("First hash should verify correctly")
	}
	if !CheckPasswordHash(password, hash2) {
		t.Error("Second hash should verify correctly")
	}
}
