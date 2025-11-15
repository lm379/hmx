package utils

import (
	"strings"

	"github.com/lm379/hmx/config"
)

// GetFullURL 将 S3 存储的相对路径转换为完整URL
// 如果路径已经是完整URL（http/https开头），直接返回
// 否则使用 custom_domain 或 endpoint 拼接
func GetFullURL(path string) string {
	if path == "" {
		return ""
	}

	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	cfg := config.AppConfig

	// 优先使用自定义域名
	if cfg.S3CustomDomain != "" {
		return strings.TrimSuffix(cfg.S3CustomDomain, "/") + "/" + strings.TrimPrefix(path, "/")
	}

	// 否则使用 S3 Endpoint
	return strings.TrimSuffix(cfg.S3Endpoint, "/") + "/" + cfg.S3Bucket + "/" + strings.TrimPrefix(path, "/")
}
