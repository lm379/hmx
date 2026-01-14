package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/database"
	api "github.com/lm379/hmx/internal/api/v1"
	"github.com/lm379/hmx/internal/queue"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/search"
)

//go:embed static
var staticFiles embed.FS

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	database.InitDB()
	database.InitRedis()

	// 初始化 Meilisearch 索引
	if err := search.InitializeIndexes(); err != nil {
		log.Printf("Warning: Failed to initialize Meilisearch indexes: %v", err)
	}

	// 启动Worker
	workers := queue.StartWorkers()
	log.Println("Workers started")

	services.InitS3(context.Background(), cfg)

	gin.SetMode(cfg.GinMode)

	// debug 模式下不嵌入静态文件
	var r *gin.Engine
	if cfg.GinMode == gin.DebugMode {
		r = api.SetupRouter(nil)
	} else {
		r = api.SetupRouter(&staticFiles)
	}

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down workers...")
		queue.StopWorkers(workers)
		log.Println("Workers stopped")
		os.Exit(0)
	}()

	r.Run(":" + cfg.ServerPort)
}
