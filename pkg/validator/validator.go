package validator

import (
	"regexp"
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	// 简单的邮箱正则
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ValidatePhone 验证手机号格式 (中国大陆手机号)
func ValidatePhone(phone string) bool {
	// 中国大陆手机号正则: 1开头，第二位3-9，后面9位数字
	const phoneRegex = `^1[3-9]\d{9}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}
