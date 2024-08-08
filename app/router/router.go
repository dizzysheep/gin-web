package router

import (
	"gin-web/app/handler"
	_ "gin-web/pkg/docs" // 这里需要引入本地已生成文档
	"github.com/gin-gonic/gin"
)

func UseIn(g *gin.Engine, handler *handler.Handlers) {
	api := g.Group("/api")
	api.GET("/health", handler.Health.Index)
	UseV1(api, handler)
}
