package handler

import (
	"gin-web/app/response"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	service *service.Services
}

func NewServerHandler(service *service.Services) *ServerHandler {
	return &ServerHandler{service}
}

// Info godoc
// @Summary 运维-服务器信息
// @Description 返回后端所在主机的系统信息与 CPU/内存/SWAP/磁盘占用，结果缓存 3 秒
// @Tags 运维
// @Produce json
// @Success 200 {object} response.Response
// @Router /v1/admin/server [get]
func (h *ServerHandler) Info(c *gin.Context) {
	response.Ok(c, h.service.Server.Info())
}
