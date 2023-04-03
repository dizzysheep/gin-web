package v1

import (
	"gin-web/core/ginc"
	"gin-web/core/jwt"
	"gin-web/models"
	"github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context) {
	username := ginc.GetString(c, "username")
	password := ginc.GetString(c, "password")
	if username == "" || password == "" {
		ginc.Fail(c, "缺少参数username|password")
		return
	}
	if !models.NewAuth().CheckAuth(username, password) {
		ginc.Fail(c, "账号或者密码不匹配")
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		ginc.Fail(c, "生成token失败")
		return
	}
	data := map[string]interface{}{"token": token}
	ginc.Ok(c, data)
}
