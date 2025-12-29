package jwtutils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lm379/hmx/config"
)

// Claims 自定义 JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT
func GenerateToken(userID uint, username string, role string) (string, error) {
	return generateTokenWithExpiry(userID, username, role, config.AppConfig.JWTAccessTokenExpiresIn)
}

// GenerateAccessToken 生成 Access Token
func GenerateAccessToken(userID uint, username string, role string) (string, error) {
	return generateTokenWithExpiry(userID, username, role, config.AppConfig.JWTAccessTokenExpiresIn)
}

// GenerateRefreshToken 生成 Refresh Token
func GenerateRefreshToken(userID uint, username string, role string) (string, error) {
	return generateTokenWithExpiry(userID, username, role, config.AppConfig.JWTRefreshTokenExpiresIn)
}

// generateTokenWithExpiry 生成指定过期时间的 Token
func generateTokenWithExpiry(userID uint, username string, role string, expiresIn time.Duration) (string, error) {
	cfg := config.AppConfig
	expiresAt := time.Now().Add(expiresIn)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "hmx-backend",
		},
	}

	// 使用 HS256 签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT
func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.AppConfig
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是 HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
