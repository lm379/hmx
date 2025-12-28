package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/jwtutils"
)

// AuthMiddleware 验证 JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// 格式: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format is 'Bearer <token>'"})
			return
		}

		tokenString := parts[1]
		claims, err := jwtutils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// 将用户信息存储到 Gin Context 中，以便后续 Handler 使用
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// TryAuthMiddleware 尝试验证 JWT Token，如果失败也继续执行（用于支持游客+用户的接口）
func TryAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims, err := jwtutils.ValidateToken(tokenString)
			if err == nil {
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
			}
		}
		c.Next()
	}
}

// AdminMiddleware 检查用户是否具有管理员权限
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Auth context not found"})
			return
		}

		// 检查角色是否为 'Administrator'
		if role.(string) != string(models.Administrator) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: Requires administrator privileges"})
			return
		}

		c.Next()
	}
}
