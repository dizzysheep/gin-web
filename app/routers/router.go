package routers

import (
	"gin-web/core/log"
	"gin-web/core/middleware"
	_ "gin-web/pkg/docs" // 这里需要引入本地已生成文档
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(middleware.Cors())
	router.Use(log.GinLogger())
	//router.Use(v3.Middleware(router, skywalking.Tracer))

	//swag api工具
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	SetupV1(router)
	SetupV2(router)

	return router
}
