package services

import (
	"context"
	"log"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/lm379/hmx/config"
)

var (
	s3PresignClient *s3.PresignClient
	s3BucketName    string
)

// InitS3 初始化 S3 客户端和 Presign 客户端
func InitS3(ctx context.Context, cfg config.Config) {
	log.Println("Initializing S3 Service...")

	// 加载 AWS 配置
	awsCfg, err := aws_config.LoadDefaultConfig(ctx,
		aws_config.WithRegion(cfg.S3Region),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, "")),
		aws_config.WithBaseEndpoint(cfg.S3Endpoint),
	)
	if err != nil {
		log.Fatalf("Failed to load S3 config: %v", err)
	}

	// 创建 S3 客户端
	s3Client := s3.NewFromConfig(awsCfg)

	// 创建 Presign 客户端
	s3PresignClient = s3.NewPresignClient(s3Client)
	s3BucketName = cfg.S3Bucket

	log.Println("S3 Service initialized successfully.")
}

// GeneratePresignedUploadURL 生成预签名的 PUT URL
// uploadType 建议为: "videos", "avatars", "covers"
// useUUID 控制是否在文件名前添加 UUID 前缀
func GeneratePresignedUploadURL(ctx context.Context, uploadType, filename, contentType string, useUUID bool) (string, string, error) {
	cfg := config.AppConfig

	var basePath string
	switch uploadType {
	case "avatars":
		basePath = cfg.S3AvatarPath
	case "videos":
		basePath = cfg.S3VideoPath
	default:
		basePath = uploadType // Fallback to uploadType if not configured
	}

	// 生成唯一的文件路径 (Object Key)
	var objectKey string
	if useUUID {
		// 添加 UUID 前缀以避免文件名冲突
		ext := filepath.Ext(filename)
		baseFilename := filename[:len(filename)-len(ext)]
		objectKey = filepath.Join(basePath, (uuid.New().String() + "-" + baseFilename + ext))
	} else {
		// 将传入文件名中的路径部分去掉，只保留文件名
		filename = filepath.Base(filename)
		objectKey = filepath.Join(basePath, filename)
	}

	// 创建 PutObject 请求
	request := &s3.PutObjectInput{
		Bucket:      aws.String(s3BucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}

	// 生成预签名 URL
	expiresIn := cfg.S3PresignExpires
	presignedURL, err := s3PresignClient.PresignPutObject(ctx, request, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		return "", "", err
	}

	return presignedURL.URL, objectKey, nil
}
