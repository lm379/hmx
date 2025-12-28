package v1

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/api/v1/middleware"
)

func SetupRouter(staticFiles embed.FS) *gin.Engine {
	r := gin.Default()
	// ... (CORS)

	v1 := r.Group("/api/v1")
	{
		// 认证路由 (Auth)
		auth := v1.Group("/auth")
		{
			auth.POST("/send-code", HandleSendCode)
			auth.POST("/register", HandleRegister)
			auth.POST("/login", HandleLogin)
			auth.POST("/forget", HandleForgetPassword)
		}

		// 用户路由 (User)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", HandleGetMe)
			users.GET("/me/likes", HandleGetUserLikes)
			users.GET("/me/favorites", HandleGetUserFavorites)
			users.GET("/me/history", HandleGetUserHistory)
		}

		// 上传路由 (Uploads)
		uploads := v1.Group("/uploads")
		uploads.Use(middleware.AuthMiddleware())
		{
			uploads.POST("/presign", HandleRequestUploadURL)
		}

		// 作品路由 (Operas)
		operas := v1.Group("/operas")
		{
			operas.GET("/", HandleGetOperas)
			operas.GET("/:id", HandleGetOperaByID)
			operas.POST("/:id/history", middleware.TryAuthMiddleware(), HandleRecordHistory) // 允许游客记录，但优先获取用户信息
			operas.POST("/", middleware.AuthMiddleware(), HandleCreateOpera)
			operas.POST("/:id/like", middleware.AuthMiddleware(), HandleToggleLike)
			operas.POST("/:id/favorite", middleware.AuthMiddleware(), HandleToggleFavorite)
			operas.POST("/:id/comments", middleware.AuthMiddleware(), HandleCreateComment)
		}

		// 艺术家路由 (Artists)
		artists := v1.Group("/artists")
		{
			artists.GET("/", HandleGetArtists)
			artists.GET("/:id", HandleGetArtistByID)
		}

		// 评论路由 (Comments)
		comments := v1.Group("/comments")
		comments.Use(middleware.AuthMiddleware())
		{
			comments.DELETE("/:id", HandleDeleteComment)
		}

		// 管理员路由 (Admin)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			admin.DELETE("/operas/:id", HandleDeleteOpera)
		}
	}

	// 静态文件服务 (Embedded Frontend)
	staticFS, err := fs.Sub(staticFiles, "static")
	if err == nil {
		// 提供 assets 目录下的静态资源
		assetsFS, err := fs.Sub(staticFS, "assets")
		if err == nil {
			r.StaticFS("/assets", http.FS(assetsFS))
		}

		// 首页路由
		r.GET("/", func(c *gin.Context) {
			data, err := staticFiles.ReadFile("static/index.html")
			if err != nil {
				c.String(http.StatusNotFound, "404 page not found")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		})
		r.NoRoute(func(c *gin.Context) {
			data, err := staticFiles.ReadFile("static/index.html")
			if err != nil {
				c.String(http.StatusNotFound, "404 page not found")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		})
	}

	return r
}
