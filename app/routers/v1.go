package routers

import (
	v1 "gin-web/app/controllers/v1"
	"gin-web/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupV1(r *gin.Engine) {

	//获取token
	apiv1 := r.Group("/v1")
	{
		v1.NewAuthController().Router(apiv1)
		v1.NewTestController().Router(apiv1)
	}
	apiJWTv1 := r.Group("/v1")
	{
		apiJWTv1.Use(middleware.JWT())
		v1.NewTagController().Router(apiJWTv1)
		v1.NewArticleController().Router(apiJWTv1)
	}

}
