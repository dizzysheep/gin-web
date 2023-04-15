package v1

import (
	"gin-web/core/ginc"
	"github.com/gin-gonic/gin"
)

type TestController struct {
}

func NewTestController() *TestController {
	return &TestController{}
}

func (tc *TestController) Router(router *gin.RouterGroup) {
	router.GET("/test", tc.GetTest) //获取单个文章
}

func (tc *TestController) GetTest(c *gin.Context) {
	ginc.Ok(c, nil)
}
