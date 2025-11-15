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
	DB  *gorm.DB
	RDB *redis.Client
	Ctx = context.Background()
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
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPass,
		DB:       config.AppConfig.RedisDB,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established.")
}
