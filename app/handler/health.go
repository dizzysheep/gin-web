package handler

import (
	"context"

	"gin-web/app/response"
	"gin-web/core/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
	Db     string `json:"db" example:"ok"`
	Redis  string `json:"redis" example:"ok"`
}

// Index godoc
// @Summary 健康检查
// @Description 探测MySQL与Redis连通性
// @Tags 系统
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *HealthHandler) Index(c *gin.Context) {
	resp := &HealthResponse{Status: "ok", Db: "ok", Redis: "ok"}

	if h.db != nil {
		if err := h.db.WithContext(c.Request.Context()).Raw("SELECT 1").Scan(&struct{}{}).Error; err != nil {
			resp.Db = "fail"
			resp.Status = "fail"
		}
	}

	if redis.RedisClient != nil {
		if err := redis.RedisClient.Ping(context.Background()).Err(); err != nil {
			resp.Redis = "fail"
			resp.Status = "fail"
		}
	}

	response.Ok(c, resp)
}
