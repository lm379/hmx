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

		// 回调路由 (Callbacks)
		callbacks := v1.Group("/callback")
		{
			callbacks.POST("/transcode", HandleTranscodeCallback)
		}

		// 作品路由 (Operas)
		operas := v1.Group("/operas")
		{
			operas.GET("/", middleware.TryAuthMiddleware(), HandleGetOperas)
			operas.GET("/:id", middleware.TryAuthMiddleware(), HandleGetOperaByID)
			operas.GET("/:id/similar", middleware.TryAuthMiddleware(), HandleGetSimilarOperas)          // 相似作品推荐
			operas.GET("/:id/summary", middleware.TryAuthMiddleware(), HandleGetVideoSummary)           // 视频AI字幕摘要
			operas.POST("/:id/request-summary", middleware.AuthMiddleware(), HandleRequestOperaSummary) // 用户请求生成AI摘要
			operas.POST("/:id/history", middleware.TryAuthMiddleware(), HandleRecordHistory)            // 允许游客记录，但优先获取用户信息
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

		// 新闻资讯路由 (News)
		news := v1.Group("/news")
		{
			news.GET("/", HandleGetNews)
			news.GET("/:id", HandleGetNewsByID)
		}

		// 评论路由 (Comments)
		comments := v1.Group("/comments")
		comments.Use(middleware.AuthMiddleware())
		{
			comments.DELETE("/:id", HandleDeleteComment)
			comments.POST("/:id/like", HandleToggleCommentLike)
		}

		// 推荐路由 (Recommendations)
		recommendations := v1.Group("/recommendations")
		{
			recommendations.GET("/", middleware.TryAuthMiddleware(), HandleGetRecommendations) // 个性化推荐（支持游客）
		}

		// 搜索路由 (Search)
		search := v1.Group("/search")
		{
			search.GET("/operas", HandleSearchOperas)   // 搜索作品
			search.GET("/artists", HandleSearchArtists) // 搜索艺术家
		}

		// 知识库问答路由 (QA)
		qa := v1.Group("/qa")
		{
			qa.POST("/ask", middleware.TryAuthMiddleware(), HandleAsk)
			qa.GET("/history", middleware.TryAuthMiddleware(), HandleGetQAHistory)
			qa.POST("/:qa_id/feedback", middleware.TryAuthMiddleware(), HandleSubmitFeedback)
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
			admin.POST("/operas/:id/embedding", HandleGenerateEmbedding)           // 生成单个向量
			admin.POST("/operas/batch-embedding", HandleBatchGenerateEmbeddings)   // 批量生成向量
			admin.POST("/operas/:id/generate-summary", HandleGenerateOperaSummary) // 生成单个AI摘要
			admin.POST("/operas/batch-summary", HandleBatchGenerateOperaSummaries) // 批量生成AI摘要

			// 任务状态查询
			admin.GET("/tasks/batch/:id", HandleGetBatchTaskStatus) // 获取批量任务状态
			admin.GET("/tasks/:id", HandleGetTaskStatus)            // 获取单个任务状态
			admin.GET("/tasks/queue/status", HandleGetQueueStatus)  // 获取队列状态

			admin.GET("/artists", HandleGetArtists)
			admin.POST("/artists", HandleAdminCreateArtist)
			admin.PUT("/artists/:id", HandleAdminUpdateArtist)
			admin.DELETE("/artists/:id", HandleAdminDeleteArtist)
			admin.GET("/users", HandleAdminGetUsers)
			admin.PUT("/users/:id/role", HandleAdminUpdateUserRole)

			admin.GET("/news", HandleAdminGetNews)
			admin.POST("/news", HandleAdminCreateNews)
			admin.GET("/news/:id", HandleAdminGetNewsByID)
			admin.PUT("/news/:id", HandleAdminUpdateNews)
			admin.DELETE("/news/:id", HandleAdminDeleteNews)

			// 搜索管理路由
			adminSearch := admin.Group("/search")
			{
				adminSearch.POST("/reindex/operas", HandleReindexOperas)   // 重新索引所有作品
				adminSearch.POST("/reindex/artists", HandleReindexArtists) // 重新索引所有艺术家
				adminSearch.GET("/stats", HandleGetSearchStats)            // 获取搜索统计
			}

			// 知识库管理路由
			knowledge := admin.Group("/knowledge")
			{
				knowledge.POST("/documents", HandleAdminUploadKnowledge)                               // 上传知识库文档
				knowledge.GET("/documents", HandleAdminGetKnowledgeDocuments)                          // 获取文档列表
				knowledge.GET("/documents/:doc_id", HandleAdminGetKnowledgeDocumentDetail)             // 获取文档详情
				knowledge.DELETE("/documents/:doc_id", HandleAdminDeleteKnowledgeDocument)             // 删除文档
				knowledge.PUT("/documents/:doc_id/reactivate", HandleAdminReactivateKnowledgeDocument) // 重新激活文档
				knowledge.GET("/stats", HandleAdminGetKnowledgeStats)                                  // 知识库统计
				knowledge.POST("/import-operas", HandleAdminImportOperas)                              // 批量导入作品字幕
			}
		}
	}

	// 静态文件服务 (Embedded Frontend)
	// debug 模式下 staticFiles 为 nil，不提供静态文件服务
	if staticFiles != nil {
		staticFS, err := fs.Sub(staticFiles, "static")
		if err == nil {
			// 提供 js 目录下的静态资源
			jsFS, err := fs.Sub(staticFS, "js")
			if err == nil {
				r.StaticFS("/js", http.FS(jsFS))
			}

			// 提供 css 目录下的静态资源
			cssFS, err := fs.Sub(staticFS, "css")
			if err == nil {
				r.StaticFS("/css", http.FS(cssFS))
			}

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
