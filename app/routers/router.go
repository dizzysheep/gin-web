package routers

import (
	middleware2 "gin-web/app/middleware"
	_ "gin-web/pkg/docs" // 这里需要引入本地已生成文档
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/yangxx0612/plugins/config"
)

func InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(middleware2.RequestID())

	router.Use(middleware2.Logger())

	router.Use(gin.Recovery())

	gin.SetMode(config.RunMode)

	//swag api工具
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	SetupV1(router)

	SetupV2(router)

	return router
}
