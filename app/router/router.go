package router

import (
	"gin-web/app/handler"
	"gin-web/core/config"
	"github.com/gin-gonic/gin"
)

func UseIn(g *gin.Engine, handler *handler.Handlers) {
	// 本地存储驱动时托管上传文件（生产环境建议由Nginx直接托管）
	if config.GetString("upload.driver") == "" || config.GetString("upload.driver") == "local" {
		urlPrefix := config.GetString("upload.urlPrefix")
		if urlPrefix != "" {
			g.Static(urlPrefix, config.GetString("upload.path"))
		}
	}

	api := g.Group("/api")
	api.GET("/health", handler.Health.Index)
	g.GET("/rss.xml", handler.Article.RSS)
	g.GET("/sitemap.xml", handler.Article.Sitemap)
	UseV1(api, handler)
}
