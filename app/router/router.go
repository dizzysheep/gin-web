package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// 管理后台前端（web/dist）：构建产物存在时由 gin 托管，开发时用 `npm run dev`（Vite 代理 /api）
	webDist := "web/dist"
	if _, err := os.Stat(webDist); err == nil {
		g.Static("/assets", filepath.Join(webDist, "assets"))
		g.NoRoute(func(c *gin.Context) {
			p := c.Request.URL.Path
			// API 与上传路径未命中时返回 404，其余路径回退到 SPA 入口（支持前端路由刷新）
			if strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "/uploads") {
				c.JSON(http.StatusNotFound, gin.H{"code": 200003, "msg": "找不到资源"})
				return
			}
			c.File(filepath.Join(webDist, "index.html"))
		})
	}
}
