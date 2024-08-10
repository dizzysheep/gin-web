package router

import (
	"gin-web/app/handler"
	"gin-web/app/middleware"
	"github.com/gin-gonic/gin"
)

func UseV1(r *gin.RouterGroup, handlers *handler.Handlers) {

	apiV1 := r.Group("/v1")
	{
		apiV1.POST("/user/login", handlers.Auth.Login)
		apiV1.GET("/common/option", handlers.Common.GetOptions)

	}

	apiJwtV1 := r.Group("/v1").Use(middleware.JWT())
	{
		//标签管理
		apiJwtV1.GET("/tag", handlers.Tag.List)
		apiJwtV1.POST("/tag", handlers.Tag.Add)
		apiJwtV1.PATCH("/tag/:id", handlers.Tag.Edit)
		apiJwtV1.DELETE("/tag/:id", handlers.Tag.Del)

		//文章管理
		apiJwtV1.GET("/article", handlers.Article.List)
		apiJwtV1.POST("/article", handlers.Article.Add)
	}
}
