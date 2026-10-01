package handler

import (
	"gin-web/app/ext"
	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/errcode"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.Services
}

func NewAuthHandler(service *service.Services) *AuthHandler {
	return &AuthHandler{service}
}

// Login godoc
// @Summary 登录
// @Description 博主登录，返回JWT
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/user/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	reqDTO, err := dto.LoginReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Auth.Login(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Logout godoc
// @Summary 登出
// @Description token加入黑名单
// @Tags 认证
// @Produce json
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /v1/user/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token := ext.ExtractToken(c)
	if token == "" {
		response.Fail(c, errcode.TokenEmpty)
		return
	}

	if err := h.service.Auth.Logout(c.Request.Context(), token); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Info godoc
// @Summary 当前用户信息
// @Tags 认证
// @Produce json
// @Success 200 {object} response.Response
// @Router /v1/user/info [get]
func (h *AuthHandler) Info(c *gin.Context) {
	respDTO, err := h.service.Auth.Info(c.Request.Context(), ext.GetUser(c))
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// ChangePassword godoc
// @Summary 修改密码
// @Description 修改成功后所有旧token失效
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/user/password [patch]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	reqDTO, err := dto.ChangePasswordReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Auth.ChangePassword(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}
