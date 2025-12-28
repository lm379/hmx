package main

import (
	"context"
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/database"
	api "github.com/lm379/hmx/internal/api/v1"
	"github.com/lm379/hmx/internal/services"
)

//go:embed static
var staticFiles embed.FS

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	database.InitDB()
	database.InitRedis()

	services.InitS3(context.Background(), cfg)

	gin.SetMode(cfg.GinMode)

	r := api.SetupRouter(staticFiles)

	r.Run(":" + cfg.ServerPort)
}
