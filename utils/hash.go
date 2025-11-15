package utils

import "golang.org/x/crypto/bcrypt"

// BcryptCost 是 bcrypt 的 cost 因子，可以通过配置修改
var BcryptCost = 12

// HashPassword 使用 bcrypt 哈希密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost) // 使用可配置的 cost
	return string(bytes), err
}

// CheckPasswordHash 检查密码哈希是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
