package middleware

import (
	"gin-web/core/config"
	"github.com/gin-gonic/gin"
)

func UserIn(g *gin.Engine) {
	g.Use(Recovery())
	g.Use(Cors())
	g.Use(SetRequestID())
	if config.IsDevEnv {
		g.Use(gin.Logger())
	} else {
		g.Use(GinLogger())
	}
}
