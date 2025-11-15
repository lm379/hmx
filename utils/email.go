package utils

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/lm379/hmx/config"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	cfg := config.AppConfig
	m := gomail.NewMessage()

	// 设置发件人（可以添加显示名称）
	m.SetHeader("From", cfg.SMTPUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// 创建拨号器
	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)

	// 配置 TLS
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
	}

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
