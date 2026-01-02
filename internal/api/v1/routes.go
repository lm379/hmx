package v1

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/api/v1/middleware"
)

func SetupRouter(staticFiles *embed.FS) *gin.Engine {
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
			auth.POST("/refresh", HandleRefreshToken)
			auth.POST("/logout", middleware.AuthMiddleware(), HandleLogout)
		}

		// 用户路由 (User)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", HandleGetMe)
			users.PUT("/me", HandleUpdateUserProfile)
			users.POST("/me/password", HandleUpdatePassword)
			users.POST("/me/email-code", HandleSendCodeToCurrentUser)
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
			operas.GET("/", middleware.TryAuthMiddleware(), HandleGetOperas)
			operas.GET("/:id", middleware.TryAuthMiddleware(), HandleGetOperaByID)
			operas.POST("/:id/history", middleware.TryAuthMiddleware(), HandleRecordHistory) // 允许游客记录，但优先获取用户信息
			operas.POST("/", middleware.AuthMiddleware(), HandleCreateOpera)
			operas.POST("/:id/like", middleware.AuthMiddleware(), HandleToggleLike)
			operas.POST("/:id/favorite", middleware.AuthMiddleware(), HandleToggleFavorite)
			operas.POST("/:id/share", middleware.TryAuthMiddleware(), HandleShare)
			operas.POST("/:id/comments", middleware.AuthMiddleware(), HandleCreateComment)
			operas.GET("/:id/comments", middleware.TryAuthMiddleware(), HandleGetComments)
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
			comments.POST("/:id/like", HandleToggleCommentLike)
		}

		// 管理员路由 (Admin)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/stats", HandleAdminGetDashboardStats)

			admin.GET("/operas", HandleAdminGetOperas)
			admin.PUT("/operas/:id", HandleAdminUpdateOpera)
			admin.DELETE("/operas/:id", HandleDeleteOpera)
			admin.GET("/artists", HandleGetArtists)
			admin.POST("/artists", HandleAdminCreateArtist)
			admin.PUT("/artists/:id", HandleAdminUpdateArtist)
			admin.DELETE("/artists/:id", HandleAdminDeleteArtist)
			admin.GET("/users", HandleAdminGetUsers)
			admin.PUT("/users/:id/role", HandleAdminUpdateUserRole)
		}
	}

	// 静态文件服务 (Embedded Frontend)
	// debug 模式下 staticFiles 为 nil，不提供静态文件服务
	if staticFiles != nil {
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
	}

	return r
}
