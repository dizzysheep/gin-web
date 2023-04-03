package routers

import (
	v1 "gin-web/controllers/v1"
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
)

func SetupV1(r *gin.Engine) {

	//获取token
	apiv1 := r.Group("/v1")

	{
		//获取token
		apiv1.GET("/auth", v1.GetAuth)
		//api测试
		apiv1.GET("/test", v1.GetTest)

		//推送kafka消息
		//apiv1.GET("/kafka/send_msg", kafka.SendMsg)
		//apiv1.GET("/kafka/consumer_msg", kafka.ConsumerMsg)
	}
	{
		apiv1.Use(middleware.JWT())

		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags/add", v1.AddTag)
		//更新指定标签
		apiv1.POST("/tags/edit", v1.EditTag)
		//删除指定标签
		apiv1.GET("/tags/delete", v1.DeleteTag)

		//获取单个文章
		apiv1.GET("/article", v1.GetArticle)
		//获取多个文章
		apiv1.GET("/articles", v1.GetArticles)
		//新建文章
		apiv1.POST("/article/add", v1.AddArticle)
		//更新文章
		apiv1.POST("/article/edit", v1.EditArticle)
		//删除文章
		apiv1.GET("/article/delete", v1.DeleteArticle)
	}

}
