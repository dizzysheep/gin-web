package router

import (
	"net/http"

	"gin-web/app/handler"
	"gin-web/core/config"
	_ "gin-web/doc/swagger" // swagger文档（swag生成于doc/swagger目录）
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

// swaggerUIPage swagger-ui静态页（CDN加载），数据源为 /swagger/doc.json
const swaggerUIPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>gin-web API 文档</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.19.1/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@4.19.1/swagger-ui-bundle.js"></script>
<script>
    window.onload = function () {
        window.ui = window.SwaggerUIBundle({
            url: "/swagger/doc.json",
            dom_id: "#swagger-ui",
            deepLinking: true,
            persistAuthorization: true
        });
    };
</script>
</body>
</html>`

func UseIn(g *gin.Engine, handler *handler.Handlers) {
	// swagger文档：/swagger/index.html（doc.json为数据源）
	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Param("any") == "/doc.json" {
			doc, err := swag.ReadDoc()
			if err != nil {
				c.String(http.StatusInternalServerError, "read swagger doc fail: "+err.Error())
				return
			}
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(doc))
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIPage))
	})

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
