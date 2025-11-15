package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/api/handlers"
	"github.com/lm379/hmx/api/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// ... (CORS)

	v1 := r.Group("/api/v1")
	{
		// 认证路由 (Auth)
		auth := v1.Group("/auth")
		{
			auth.POST("/send-code", handlers.HandleSendCode)
			auth.POST("/register", handlers.HandleRegister)
			auth.POST("/login", handlers.HandleLogin)
			auth.POST("/forget", handlers.HandleForgetPassword)
		}

		// 用户路由 (User)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", handlers.HandleGetMe)
			users.GET("/me/likes", handlers.HandleGetUserLikes)
			users.GET("/me/favorites", handlers.HandleGetUserFavorites)
		}

		// 上传路由 (Uploads)
		uploads := v1.Group("/uploads")
		uploads.Use(middleware.AuthMiddleware())
		{
			uploads.POST("/presign", handlers.HandleRequestUploadURL)
		}

		// 作品路由 (Operas)
		operas := v1.Group("/operas")
		{
			operas.GET("/", handlers.HandleGetOperas)
			operas.GET("/:id", handlers.HandleGetOperaByID)
			operas.POST("/", middleware.AuthMiddleware(), handlers.HandleCreateOpera)
			operas.POST("/:id/like", middleware.AuthMiddleware(), handlers.HandleToggleLike)
			operas.POST("/:id/favorite", middleware.AuthMiddleware(), handlers.HandleToggleFavorite)
			operas.POST("/:id/comments", middleware.AuthMiddleware(), handlers.HandleCreateComment)
		}

		// 评论路由 (Comments)
		comments := v1.Group("/comments")
		comments.Use(middleware.AuthMiddleware())
		{
			comments.DELETE("/:id", handlers.HandleDeleteComment)
		}

		// 管理员路由 (Admin)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			admin.DELETE("/operas/:id", handlers.HandleDeleteOpera)
		}
	}

	return r
}
