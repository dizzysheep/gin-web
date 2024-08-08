package handler

import (
	"gin-web/app/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Index(c *gin.Context) {
	response.Ok(c, "success")
}
