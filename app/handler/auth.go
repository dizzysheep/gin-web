package handler

import (
	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.Services
}

func NewAuthHandler(service *service.Services) *AuthHandler {
	return &AuthHandler{service}
}

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
