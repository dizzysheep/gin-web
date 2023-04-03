package routers

import (
	_ "gin-demo/docs" // 这里需要引入本地已生成文档
	"gin-demo/middleware"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/yangxx0612/plugins/config"
)

func InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(middleware.RequestID())

	router.Use(middleware.Logger())

	router.Use(gin.Recovery())

	gin.SetMode(config.RunMode)

	//swag api工具
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	SetupV1(router)

	SetupV2(router)

	return router
}
