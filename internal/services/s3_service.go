package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/internal/models"
)

// 定义权限错误
var (
	ErrPermissionDenied  = errors.New("permission denied")
	ErrInvalidUploadType = errors.New("invalid upload type")
	ErrOperaNotFound     = errors.New("opera not found")
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
// uploadType: "user_avatar", "artist_avatar", "opera_cover", "news_cover", "education_pdf", "education_cover", "video_upload"
// userID: 当前登录用户的ID（从JWT获取）
// userRole: 当前用户的角色（从JWT获取）
// targetID: 目标资源ID（用户ID、艺术家ID、曲目ID等），video_upload时可为0
// filename: 文件名
// contentType: 文件类型
func GeneratePresignedUploadURL(ctx context.Context, uploadType string, userID uint, userRole models.UserRole, targetID uint, filename, contentType string) (string, string, error) {
	cfg := config.AppConfig

	var objectKey string
	ext := filepath.Ext(filename)

	switch uploadType {
	case "user_avatar":
		// 用户头像：只能上传到自己的目录
		if targetID != userID {
			return "", "", fmt.Errorf("%w: cannot upload to other user's avatar directory", ErrPermissionDenied)
		}
		// 路径: user/avatar/{userID}/uuid.ext
		objectKey = filepath.Join("user", "avatar", strconv.FormatUint(uint64(userID), 10), uuid.New().String()+ext)

	case "artist_avatar":
		// 艺术家头像：需要管理员权限
		if userRole != models.Administrator {
			return "", "", fmt.Errorf("%w: administrator role required for artist avatar upload", ErrPermissionDenied)
		}
		// 路径: artists/avatar/{artistID}/uuid.ext
		objectKey = filepath.Join("artists", "avatar", strconv.FormatUint(uint64(targetID), 10), uuid.New().String()+ext)

	case "opera_cover":
		// 视频封面：需要管理员权限
		if userRole != models.Administrator {
			return "", "", fmt.Errorf("%w: administrator role required for opera cover upload", ErrPermissionDenied)
		}
		// 路径: public/cover/{operaTitle}/{operaTitle}.ext
		// 使用targetID作为operaID，从数据库获取曲目名
		opera, err := GetOperaByID(targetID)
		if err != nil {
			return "", "", fmt.Errorf("%w: %v", ErrOperaNotFound, err)
		}
		// 路径: public/cover/{title}/{title}.ext
		objectKey = filepath.Join("public", "cover", opera.OperaTitle, opera.OperaTitle+ext)

	case "news_cover":
		// 新闻封面由管理员上传，新闻保存时只记录对象存储 key。
		if userRole != models.Administrator {
			return "", "", fmt.Errorf("%w: administrator role required for news cover upload", ErrPermissionDenied)
		}
		objectKey = filepath.Join("public", "news", uuid.New().String()+ext)

	case "education_pdf":
		// 黄梅教育 PDF 由管理员上传，前台用返回的完整 URL 直接加载。
		if userRole != models.Administrator {
			return "", "", fmt.Errorf("%w: administrator role required for education pdf upload", ErrPermissionDenied)
		}
		objectKey = filepath.Join("public", "education", uuid.New().String()+ext)

	case "education_cover":
		// 教育资源封面由管理员上传，通常由前端从 PDF 第一页自动生成。
		if userRole != models.Administrator {
			return "", "", fmt.Errorf("%w: administrator role required for education cover upload", ErrPermissionDenied)
		}
		objectKey = filepath.Join("public", "education", "cover", uuid.New().String()+ext)

	case "video_upload":
		// 视频上传到tmp：需要管理员权限，上传后由腾讯云数据万象处理并回调
		if userRole != models.Administrator {
			return "", "", fmt.Errorf("%w: administrator role required for video upload", ErrPermissionDenied)
		}
		// 路径: tmp/{title}.{ext}
		// filename 应该是前端传入的视频标题 + 扩展名，例如 "天仙配.mp4"
		filename = filepath.Base(filename) // 移除路径部分，只保留文件名
		objectKey = filepath.Join("tmp", filename)

	default:
		return "", "", ErrInvalidUploadType
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
