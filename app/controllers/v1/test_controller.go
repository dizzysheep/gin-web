package v1

import (
	"gin-web/core/ginc"
	"gin-web/core/netlib"
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
	url := "https://ug.baidu.com/mcp/pc/pcsearch"
	params := map[string]interface{}{
		"invoke_info": map[string][]interface{}{
			"pos_1": {},
			"pos_2": {},
			"pos_3": {},
		},
	}
	req, _ := netlib.Post(url).SetCtx(c).JSONBody(params)
	res, _ := req.ToMap()
	ginc.Ok(c, res)
}
