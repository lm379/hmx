package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// --- 自定义 ENUM 类型 ---
type UserSex string

const (
	Male   UserSex = "Male"
	Female UserSex = "Female"
	Other  UserSex = "Other"
)

type UserRole string

const (
	Administrator UserRole = "Administrator"
	User          UserRole = "User"
)

// Users (用户表)
type Users struct {
	UserID      uint           `gorm:"column:user_id;primaryKey"`
	Username    string         `gorm:"column:username;type:varchar(50);unique;not null"`
	Phone       string         `gorm:"column:phone;type:varchar(15);unique;not null"`
	Email       sql.NullString `gorm:"column:email;type:varchar(255);unique"`
	Password    string         `gorm:"column:password;type:varchar(255);not null"`
	Sex         UserSex        `gorm:"column:sex;type:user_sex;default:Other;not null"`
	Icon        sql.NullString `gorm:"column:icon;type:varchar(255)"`
	Role        UserRole       `gorm:"column:role;type:user_role;default:User;not null"`
	LastLoginAt *time.Time     `gorm:"column:last_login_at"`
	LastIp      string         `gorm:"column:last_ip;type:varchar(50)"`
	CreatedAt   time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 指定表名
func (Users) TableName() string {
	return "users"
}

// UserResponse 用于 API 返回的 DTO，自动处理 URL 拼接
type UserResponse struct {
	UserID      uint       `json:"user_id"`
	Username    string     `json:"username"`
	Phone       string     `json:"phone"`
	Email       *string    `json:"email"`
	Sex         UserSex    `json:"sex"`
	Icon        *string    `json:"icon"` // 完整 URL
	Role        UserRole   `json:"role"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// RegisterInput 注册时绑定的 JSON
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code" binding:"required"` // 邮箱验证码
}

// LoginInput 登录时绑定的 JSON
type LoginInput struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SendCodeInput 发送验证码时绑定的 JSON
type SendCodeInput struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgetPasswordInput 忘记密码时绑定的 JSON
type ForgetPasswordInput struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdateUserProfileRequest DTO (PUT /users/me)
type UpdateUserProfileRequest struct {
	Username string  `json:"username"`
	Phone    string  `json:"phone"`
	Sex      UserSex `json:"sex" binding:"omitempty,oneof=Male Female Other"`
	Icon     string  `json:"icon"` // Object Key for avatar
	// Email modification requires verification code
	Email string `json:"email" binding:"omitempty,email"`
	Code  string `json:"code"` // Code sent to *current* email, required if Email is being changed
}

// UpdatePasswordRequest DTO (POST /users/me/password)
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
