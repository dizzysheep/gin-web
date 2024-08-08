package router

import (
	"gin-web/app/handler"
	"github.com/gin-gonic/gin"
)

func UseV1(r *gin.RouterGroup, handlers *handler.Handlers) {

	apiv1 := r.Group("/v1")
	{
		apiv1.POST("/user/login", handlers.Auth.Login)
		apiv1.GET("/common/option", handlers.Common.GetOptions)

		apiv1.GET("/tags", handlers.Tag.List)

	}

}
