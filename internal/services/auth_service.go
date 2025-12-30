package services

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/emailutils"
	"github.com/lm379/hmx/pkg/hashutils"
	"github.com/lm379/hmx/pkg/jwtutils"
	"github.com/lm379/hmx/pkg/validator"
	"gorm.io/gorm"
)

type EmailTemplateData struct {
	VerificationCode string
	ExpiryTime       int
}

// VerifyCode 验证验证码
func VerifyCode(email, code string) bool {
	redisKey := "verify_code:" + email
	storedCode, err := database.RDBVerifyCode.Get(database.Ctx, redisKey).Result()
	if err != nil {
		return false
	}
	return storedCode == code
}

// DeleteCode 删除验证码
func DeleteCode(email string) {
	redisKey := "verify_code:" + email
	database.RDBVerifyCode.Del(database.Ctx, redisKey)
}

// SendVerificationCode 生成并发送验证码
func SendVerificationCode(email string) (gin.H, int) {
	if !validator.ValidateEmail(email) {
		return gin.H{"error": "Invalid email format"}, http.StatusBadRequest
	}

	// 生成 6 位随机数
	max := big.NewInt(1000000)
	randomNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return gin.H{"error": "Failed to generate verification code"}, http.StatusInternalServerError
	}
	code := fmt.Sprintf("%06v", randomNum)

	// 存储到 Redis，设置过期时间
	redisKey := "verify_code:" + email
	expireTime := time.Duration(config.AppConfig.SMTPCodeExpires) * time.Minute
	err = database.RDBVerifyCode.Set(database.Ctx, redisKey, code, expireTime).Err()
	if err != nil {
		return gin.H{"error": "Failed to store verification code"}, http.StatusInternalServerError
	}

	templateStr := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>邮箱验证码</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            font-size: 14px;
            line-height: 1.6;
            color: #333;
        }
        .container {
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            background-color: #415A94;
            color: #fff;
            padding: 20px 40px;
            font-size: 21px;
            text-align: center;
        }
        .content {
            padding: 30px 40px;
        }
        .verification-code {
            text-align: center;
            margin: 25px 0;
            font-size: 32px;
            font-weight: bold;
            color: #415A94;
            letter-spacing: 5px;
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            border: 2px dashed #dee2e6;
        }
        .footer {
            padding: 20px 40px;
            font-size: 12px;
            color: #999;
            line-height: 20px;
            background: #f7f7f7;
            text-align: center;
            border-top: 1px solid #e0e0e0;
        }
        .tips {
            font-size: 12px;
            color: #666;
            margin-top: 10px;
        }
        .highlight {
            color: #e74c3c;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            安全验证
        </div>
        <div class="content">
            <p>尊敬的用户，您好！</p>
            <p>您正在进行敏感操作，请在验证页面输入以下验证码以完成操作：</p>
            
            <div class="verification-code">{{.VerificationCode}}</div>
            
            <p class="tips">此验证码 <span class="highlight">{{.ExpiryTime}}</span> 分钟内有效。请勿将验证码泄露给他人。</p>
            <p>如非本人操作，请忽略此邮件。</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿直接回复。</p>
        </div>
    </div>
</body>
</html>`
	data := EmailTemplateData{
		VerificationCode: code,
		ExpiryTime:       int(expireTime.Minutes()),
	}
	subject := "您的验证码"
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return gin.H{"error": "Failed to parse email template"}, http.StatusInternalServerError
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return gin.H{"error": "Failed to execute email template"}, http.StatusInternalServerError
	}

	err = emailutils.SendEmail(email, subject, body.String())
	if err != nil {
		return gin.H{"error": "Failed to send verification email"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Verification code sent (check logs)."}, http.StatusOK
}

// RegisterUser 注册新用户
func RegisterUser(input models.RegisterInput) (gin.H, int) {
	if !validator.ValidateEmail(input.Email) {
		return gin.H{"error": "Invalid email format"}, http.StatusBadRequest
	}
	if !validator.ValidatePhone(input.Phone) {
		return gin.H{"error": "Invalid phone number format"}, http.StatusBadRequest
	}

	// 验证 Redis 中的验证码
	redisKey := "verify_code:" + input.Email
	code, err := database.RDBVerifyCode.Get(database.Ctx, redisKey).Result()

	if err != nil {
		return gin.H{"error": "Verification code expired or invalid"}, http.StatusBadRequest
	}
	if code != input.Code {
		return gin.H{"error": "Verification code incorrect"}, http.StatusBadRequest
	}

	// 检查用户是否已存在 (Email, Phone, Username)
	var existingUser models.Users
	if err := database.DB.Where("email = ? OR phone = ? OR username = ?", input.Email, input.Phone, input.Username).First(&existingUser).Error; err == nil {
		return gin.H{"error": "Email, Phone or Username already registered"}, http.StatusConflict
	}

	// 哈希密码
	hashedPassword, err := hashutils.HashPassword(input.Password)
	if err != nil {
		return gin.H{"error": "Failed to hash password"}, http.StatusInternalServerError
	}

	// 创建用户
	newUser := models.Users{
		Username: input.Username,
		Phone:    input.Phone,
		Email:    sql.NullString{String: input.Email, Valid: true},
		Password: hashedPassword,
		Role:     models.User,  // 默认为普通用户
		Sex:      models.Other, // 默认
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return gin.H{"error": "Failed to create user"}, http.StatusInternalServerError
	}

	// 注册成功，删除验证码
	database.RDBVerifyCode.Del(database.Ctx, redisKey)

	return gin.H{"message": "Registration successful"}, http.StatusCreated
}

// LoginUser 登录用户并返回 JWT
func LoginUser(input models.LoginInput, clientIP string) (gin.H, int) {
	var user models.Users

	// 查找用户 (支持 Email 或 Username)
	if err := database.DB.Where("email = ? OR username = ?", input.Account, input.Account).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gin.H{"error": "Invalid credentials"}, http.StatusUnauthorized
		}
		return gin.H{"error": "Database error"}, http.StatusInternalServerError
	}

	// 检查密码
	if !hashutils.CheckPasswordHash(input.Password, user.Password) {
		return gin.H{"error": "Invalid credentials"}, http.StatusUnauthorized
	}

	// 更新最后登录时间和 IP
	now := time.Now()
	database.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_ip":       clientIP,
	})

	// 生成 Access Token 和 Refresh Token
	accessToken, err := jwtutils.GenerateAccessToken(user.UserID, user.Username, string(user.Role))
	if err != nil {
		return gin.H{"error": "Failed to generate access token"}, http.StatusInternalServerError
	}

	refreshToken, err := jwtutils.GenerateRefreshToken(user.UserID, user.Username, string(user.Role))
	if err != nil {
		return gin.H{"error": "Failed to generate refresh token"}, http.StatusInternalServerError
	}

	// 将 Refresh Token 存储到 Redis
	refreshTokenKey := fmt.Sprintf("refresh_token:%d", user.UserID)
	err = database.RDBToken.Set(database.Ctx, refreshTokenKey, refreshToken, config.AppConfig.JWTRefreshTokenExpiresIn).Err()
	if err != nil {
		return gin.H{"error": "Failed to store refresh token"}, http.StatusInternalServerError
	}

	return gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, http.StatusOK
}

// ForgetPassword 忘记密码处理
func ForgetPassword(input models.ForgetPasswordInput) (gin.H, int) {
	// 验证 Redis 中的验证码
	redisKey := "verify_code:" + input.Email
	code, err := database.RDBVerifyCode.Get(database.Ctx, redisKey).Result()

	if err != nil {
		return gin.H{"error": "Verification code expired or invalid"}, http.StatusBadRequest
	}
	if code != input.Code {
		return gin.H{"error": "Verification code incorrect"}, http.StatusBadRequest
	}

	// 查找用户
	var user models.Users
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gin.H{"error": "User not found"}, http.StatusNotFound
		}
		return gin.H{"error": "Database error"}, http.StatusInternalServerError
	}

	// 哈希新密码
	hashedPassword, err := hashutils.HashPassword(input.NewPassword)
	if err != nil {
		return gin.H{"error": "Failed to hash password"}, http.StatusInternalServerError
	}

	// 更新密码
	if err := database.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return gin.H{"error": "Failed to update password"}, http.StatusInternalServerError
	}

	// 删除验证码
	database.RDBVerifyCode.Del(database.Ctx, redisKey)

	return gin.H{"message": "Password reset successful"}, http.StatusOK
}

// RefreshTokens 刷新 Access Token
func RefreshTokens(refreshToken string) (gin.H, int) {
	// 验证 Refresh Token
	claims, err := jwtutils.ValidateToken(refreshToken)
	if err != nil {
		return gin.H{"error": "Invalid refresh token"}, http.StatusUnauthorized
	}

	// 检查 Redis 中是否存在该 Refresh Token
	refreshTokenKey := fmt.Sprintf("refresh_token:%d", claims.UserID)
	storedToken, err := database.RDBToken.Get(database.Ctx, refreshTokenKey).Result()
	if err != nil {
		return gin.H{"error": "Refresh token expired or invalid"}, http.StatusUnauthorized
	}

	if storedToken != refreshToken {
		return gin.H{"error": "Refresh token mismatch"}, http.StatusUnauthorized
	}

	// 生成新的 Access Token
	newAccessToken, err := jwtutils.GenerateAccessToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return gin.H{"error": "Failed to generate new access token"}, http.StatusInternalServerError
	}

	// 可选：生成新的 Refresh Token（Token 轮换）
	newRefreshToken, err := jwtutils.GenerateRefreshToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return gin.H{"error": "Failed to generate new refresh token"}, http.StatusInternalServerError
	}

	// 更新 Redis 中的 Refresh Token
	err = database.RDBToken.Set(database.Ctx, refreshTokenKey, newRefreshToken, config.AppConfig.JWTRefreshTokenExpiresIn).Err()
	if err != nil {
		return gin.H{"error": "Failed to update refresh token"}, http.StatusInternalServerError
	}

	return gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}, http.StatusOK
}

// Logout 登出，删除 Refresh Token
func Logout(userID uint) (gin.H, int) {
	refreshTokenKey := fmt.Sprintf("refresh_token:%d", userID)
	err := database.RDBToken.Del(database.Ctx, refreshTokenKey).Err()
	if err != nil {
		return gin.H{"error": "Failed to logout"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Logout successful"}, http.StatusOK
}
