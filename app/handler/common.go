package handler

import (
	"gin-web/app/response"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	service *service.Services
}

func NewCommonHandler(service *service.Services) *CommonHandler {
	return &CommonHandler{service}
}

func (h *CommonHandler) GetOptions(c *gin.Context) {
	options := map[string]interface{}{}
	response.Ok(c, options)
}
