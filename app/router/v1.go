package router

import (
	"gin-web/app/handler"
	"gin-web/app/middleware"
	"github.com/gin-gonic/gin"
)

// UseV1 v1版本路由：公开读免认证，管理写走JWT
func UseV1(r *gin.RouterGroup, handlers *handler.Handlers) {

	apiV1 := r.Group("/v1")
	{
		// 认证（登录免鉴权）
		apiV1.POST("/user/login", handlers.Auth.Login)

		// 公开：文章
		apiV1.GET("/article", handlers.Article.List)
		apiV1.GET("/article/archive", handlers.Article.Archive)
		apiV1.GET("/article/slug/:slug", handlers.Article.DetailBySlug)
		apiV1.GET("/article/:id", handlers.Article.Detail)
		apiV1.GET("/article/:id/comment", handlers.Comment.ListByArticle)
		apiV1.POST("/article/:id/comment", handlers.Comment.Add)

		// 公开：分类/标签/友链/站点配置
		apiV1.GET("/category", handlers.Category.List)
		apiV1.GET("/tag", handlers.Tag.List)
		apiV1.GET("/link", handlers.Link.List)
		apiV1.GET("/common/option", handlers.Common.GetOptions)
	}

	apiJwtV1 := r.Group("/v1").Use(middleware.JWT())
	{
		//agent
		apiJwtV1.GET("/admin/agent/tools", handlers.Agent.Tools)
		apiJwtV1.POST("/admin/agent/chat", handlers.Agent.Chat)

		// 认证
		apiJwtV1.POST("/user/logout", handlers.Auth.Logout)
		apiJwtV1.GET("/user/info", handlers.Auth.Info)
		apiJwtV1.PATCH("/user/password", handlers.Auth.ChangePassword)

		// 文章管理
		apiJwtV1.POST("/article", handlers.Article.Add)
		apiJwtV1.PATCH("/article/:id", handlers.Article.Edit)
		apiJwtV1.DELETE("/article/:id", handlers.Article.Del)
		apiJwtV1.PATCH("/article/:id/publish", handlers.Article.Publish)
		apiJwtV1.PATCH("/article/:id/state", handlers.Article.State)
		apiJwtV1.PATCH("/article/:id/top", handlers.Article.Top)

		// 管理端列表（含草稿）
		apiJwtV1.GET("/admin/article", handlers.Article.AdminList)
		apiJwtV1.GET("/admin/article/:id", handlers.Article.AdminDetail)
		apiJwtV1.GET("/admin/category", handlers.Category.List)
		apiJwtV1.GET("/admin/tag", handlers.Tag.List)
		apiJwtV1.GET("/admin/link", handlers.Link.List)

		// 分类管理
		apiJwtV1.POST("/category", handlers.Category.Add)
		apiJwtV1.PATCH("/category/:id", handlers.Category.Edit)
		apiJwtV1.DELETE("/category/:id", handlers.Category.Del)

		// 标签管理
		apiJwtV1.POST("/tag", handlers.Tag.Add)
		apiJwtV1.PATCH("/tag/:id", handlers.Tag.Edit)
		apiJwtV1.DELETE("/tag/:id", handlers.Tag.Del)

		// 评论管理
		apiJwtV1.GET("/admin/comment", handlers.Comment.AdminList)
		apiJwtV1.PATCH("/admin/comment/:id/state", handlers.Comment.Audit)
		apiJwtV1.DELETE("/admin/comment/:id", handlers.Comment.Del)

		// 友链管理
		apiJwtV1.POST("/link", handlers.Link.Add)
		apiJwtV1.PATCH("/link/:id", handlers.Link.Edit)
		apiJwtV1.DELETE("/link/:id", handlers.Link.Del)

		// 站点配置
		apiJwtV1.PUT("/admin/option", handlers.Common.SaveOptions)

		// 首页工作台
		apiJwtV1.GET("/admin/dashboard/overview", handlers.Dashboard.Overview)
		apiJwtV1.GET("/admin/dashboard/trend", handlers.Dashboard.Trend)

		// 运维管理
		apiJwtV1.GET("/admin/server", handlers.Server.Info)

		// 文件上传
		apiJwtV1.POST("/upload/image", handlers.Upload.UploadImage)
	}
}
