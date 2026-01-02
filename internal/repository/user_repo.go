package repository

import (
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/pkg/pagination"
	"gorm.io/gorm"
)

// UserRepo 用户相关的数据库操作
type UserRepo struct{}

// NewUserRepo 创建用户仓库实例
func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// getDB 获取数据库实例
func (r *UserRepo) getDB() *gorm.DB {
	return database.DB
}

// GetAll 获取所有用户 (带分页)
func (r *UserRepo) GetAll(pagination *pagination.Pagination) ([]models.Users, int64, error) {
	var users []models.Users
	var total int64
	db := r.getDB().Model(&models.Users{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Scopes(pagination.Paginate()).Find(&users).Error
	return users, total, err
}

// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(id uint) (*models.Users, error) {
	var user models.Users
	err := r.getDB().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(username string) (*models.Users, error) {
	var user models.Users
	err := r.getDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepo) GetByEmail(email string) (*models.Users, error) {
	var user models.Users
	err := r.getDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *UserRepo) GetByPhone(phone string) (*models.Users, error) {
	var user models.Users
	err := r.getDB().Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepo) Create(user *models.Users) error {
	return r.getDB().Create(user).Error
}

// Update 更新用户
func (r *UserRepo) Update(id uint, updates map[string]interface{}) error {
	return r.getDB().Model(&models.Users{}).Where("user_id = ?", id).Updates(updates).Error
}

// UpdateRole 更新用户角色
func (r *UserRepo) UpdateRole(id uint, role string) error {
	return r.getDB().Model(&models.Users{}).Where("user_id = ?", id).Update("role", role).Error
}

// Delete 删除用户
func (r *UserRepo) Delete(id uint) error {
	return r.getDB().Delete(&models.Users{}, id).Error
}

// Exists 检查用户是否存在
func (r *UserRepo) Exists(id uint) (bool, error) {
	var count int64
	err := r.getDB().Model(&models.Users{}).Where("user_id = ?", id).Count(&count).Error
	return count > 0, err
}

// UsernameExists 检查用户名是否存在
func (r *UserRepo) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.getDB().Model(&models.Users{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// EmailExists 检查邮箱是否存在
func (r *UserRepo) EmailExists(email string) (bool, error) {
	var count int64
	err := r.getDB().Model(&models.Users{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// PhoneExists 检查手机号是否存在
func (r *UserRepo) PhoneExists(phone string) (bool, error) {
	var count int64
	err := r.getDB().Model(&models.Users{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}

// GetByEmailOrUsername 根据邮箱或用户名获取用户
func (r *UserRepo) GetByEmailOrUsername(account string) (*models.Users, error) {
	var user models.Users
	err := r.getDB().Where("email = ? OR username = ?", account, account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckDuplicate 检查邮箱、手机号或用户名是否已存在
func (r *UserRepo) CheckDuplicate(email, phone, username string) (bool, error) {
	var count int64
	err := r.getDB().Model(&models.Users{}).
		Where("email = ? OR phone = ? OR username = ?", email, phone, username).
		Count(&count).Error
	return count > 0, err
}

// UpdateLoginInfo 更新用户登录信息
func (r *UserRepo) UpdateLoginInfo(id uint, lastLoginAt interface{}, lastLoginIP string) error {
	return r.getDB().Model(&models.Users{}).Where("user_id = ?", id).Updates(map[string]interface{}{
		"last_login_at": lastLoginAt,
		"last_ip":       lastLoginIP,
	}).Error
}

// UpdatePassword 更新用户密码
func (r *UserRepo) UpdatePassword(id uint, hashedPassword string) error {
	return r.getDB().Model(&models.Users{}).Where("user_id = ?", id).Update("password", hashedPassword).Error
}
