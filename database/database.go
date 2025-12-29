package database

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/lm379/hmx/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB            *gorm.DB
	RDB           *redis.Client // 默认 Redis 客户端（0号DB）
	RDBToken      *redis.Client // Refresh Token Redis 客户端（0号DB）
	RDBVerifyCode *redis.Client // 验证码 Redis 客户端（1号DB）
	RDBCache      *redis.Client // 缓存 Redis 客户端（2号DB）
	Ctx           = context.Background()
)

// InitDB 初始化 GORM
func InitDB() {
	var err error
	dsn := "host=" + config.AppConfig.DBAddress +
		" user=" + config.AppConfig.DBUser +
		" password=" + config.AppConfig.DBPassword +
		" dbname=" + config.AppConfig.DBName +
		" port=" + strconv.Itoa(config.AppConfig.DBPort) +
		" sslmode=disable TimeZone=Asia/Shanghai"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info), // 打印 SQL 日志
		DisableForeignKeyConstraintWhenMigrating: true,                                // 禁用自动创建外键
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// 验证数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database ready (using existing schema from SQL file).")
}

// InitRedis 初始化 Redis
func InitRedis() {
	// 0号DB: Refresh Token
	RDBToken = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPass,
		DB:       0,
	})

	_, err := RDBToken.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis DB 0 (Token): %v", err)
	}
	log.Println("Redis DB 0 (Refresh Token) connection established.")

	// 1号DB: 验证码
	RDBVerifyCode = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPass,
		DB:       1,
	})

	_, err = RDBVerifyCode.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis DB 1 (VerifyCode): %v", err)
	}
	log.Println("Redis DB 1 (Verify Code) connection established.")

	// 2号DB: 缓存
	RDBCache = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPass,
		DB:       2,
	})

	_, err = RDBCache.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis DB 2 (Cache): %v", err)
	}
	log.Println("Redis DB 2 (Cache) connection established.")

	// 默认客户端指向 0号DB
	RDB = RDBToken

	log.Println("All Redis connections established.")
}
