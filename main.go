package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/api"
	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/services"
)

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	database.InitDB()
	database.InitRedis()

	services.InitS3(context.Background(), cfg)

	gin.SetMode(cfg.GinMode)

	r := api.SetupRouter()

	r.Run(":" + cfg.ServerPort)
}
