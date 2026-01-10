package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GinMode                  string        `mapstructure:"GIN_MODE"`
	ServerPort               string        `mapstructure:"SERVER_PORT"`
	DBUser                   string        `mapstructure:"DB_USER"`
	DBPassword               string        `mapstructure:"DB_PASSWORD"`
	DBAddress                string        `mapstructure:"DB_ADDRESS"`
	DBPort                   int           `mapstructure:"DB_PORT"`
	DBName                   string        `mapstructure:"DB_NAME"`
	RedisAddr                string        `mapstructure:"REDIS_ADDR"`
	RedisPass                string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB                  int           `mapstructure:"REDIS_DB"`
	JWTSecret                string        `mapstructure:"JWT_SECRET_KEY"`
	JWTAccessTokenExpiresIn  time.Duration `mapstructure:"JWT_ACCESS_TOKEN_EXPIRES_IN"`
	JWTRefreshTokenExpiresIn time.Duration `mapstructure:"JWT_REFRESH_TOKEN_EXPIRES_IN"`

	S3Endpoint       string        `mapstructure:"S3_ENDPOINT"`
	S3Region         string        `mapstructure:"S3_REGION"`
	S3AccessKey      string        `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretKey      string        `mapstructure:"S3_SECRET_ACCESS_KEY"`
	S3Bucket         string        `mapstructure:"S3_BUCKET_NAME"`
	S3PresignExpires time.Duration `mapstructure:"S3_PRESIGN_EXPIRES_IN_MINUTES"`
	S3CustomDomain   string        `mapstructure:"S3_CUSTOM_DOMAIN"`
	S3AvatarPath     string        `mapstructure:"S3_AVATAR_PATH"`
	S3VideoPath      string        `mapstructure:"S3_VIDEO_PATH"`

	SMTPHost        string `mapstructure:"SMTP_HOST"`
	SMTPPort        int    `mapstructure:"SMTP_PORT"`
	SMTPUser        string `mapstructure:"SMTP_USER"`
	SMTPPass        string `mapstructure:"SMTP_PASS"`
	SMTPCodeExpires uint   `mapstructure:"SMTP_CODE_EXPIRES"`

	// AI服务配置
	// Embedding服务
	EmbeddingAPIKey  string `mapstructure:"EMBEDDING_API_KEY"`
	EmbeddingBaseURL string `mapstructure:"EMBEDDING_BASE_URL"`
	EmbeddingModel   string `mapstructure:"EMBEDDING_MODEL"`

	// 视频总结服务（基于AI字幕）
	SubtitleAPIKey  string `mapstructure:"SUBTITLE_API_KEY"`
	SubtitleBaseURL string `mapstructure:"SUBTITLE_BASE_URL"`
	SubtitleModel   string `mapstructure:"SUBTITLE_MODEL"`
}

var AppConfig Config

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.SetDefault("S3_PRESIGN_EXPIRES_IN_MINUTES", 15)
	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRES_IN", "15m")
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRES_IN", "168h")
	viper.SetDefault("S3_AVATAR_PATH", "avatars")
	viper.SetDefault("S3_VIDEO_PATH", "videos")

	// AI服务默认配置
	viper.SetDefault("EMBEDDING_BASE_URL", "https://api.openai.com/v1")
	viper.SetDefault("EMBEDDING_MODEL", "text-embedding-3-small")
	viper.SetDefault("SUBTITLE_BASE_URL", "https://api.openai.com/v1")
	viper.SetDefault("SUBTITLE_MODEL", "gpt-4.1")

	AppConfig.S3PresignExpires = viper.GetDuration("S3_PRESIGN_EXPIRES_IN_MINUTES") * time.Minute
	AppConfig.JWTAccessTokenExpiresIn = viper.GetDuration("JWT_ACCESS_TOKEN_EXPIRES_IN")
	AppConfig.JWTRefreshTokenExpiresIn = viper.GetDuration("JWT_REFRESH_TOKEN_EXPIRES_IN")

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Println("Configuration loaded successfully.")
}
