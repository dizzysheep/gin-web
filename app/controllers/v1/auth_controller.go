package v1

import (
	"gin-web/core/ginc"
	"gin-web/core/jwt"
	"gin-web/internal/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) Router(router *gin.RouterGroup) {
	//获取token
	router.POST("/auth", ac.GetAuth)
}

func (ac *AuthController) GetAuth(c *gin.Context) {
	var req dto.AuthReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ginc.Fail(c, err.Error())
		return
	}

	if err := service.NewAuthService().CheckAuth(req); err != nil {
		ginc.Fail(c, err.Error())
		return
	}

	token, err := util.GenerateToken(req.Username, req.Password)
	if err != nil {
		ginc.Fail(c, "生成token失败")
		return
	}
	data := map[string]interface{}{"token": token}
	ginc.Ok(c, data)
}
