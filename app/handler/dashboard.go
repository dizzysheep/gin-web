package handler

import (
	"strconv"

	"gin-web/app/response"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *service.Services
}

func NewDashboardHandler(service *service.Services) *DashboardHandler {
	return &DashboardHandler{service}
}

// Overview godoc
// @Summary 管理-首页数据概览
// @Description 返回首页数据概览卡片所需的文章、分类、标签和待审评论数量
// @Tags 首页
// @Produce json
// @Success 200 {object} response.Response
// @Router /v1/admin/dashboard/overview [get]
func (h *DashboardHandler) Overview(c *gin.Context) {
	result, err := h.service.Dashboard.Overview(c.Request.Context())
	if err != nil {
		response.FailErr(c, err)
		return
	}
	response.Ok(c, result)
}

// Trend godoc
// @Summary 管理-首页内容趋势
// @Description 返回指定天数内每日发布文章数量，缺少数据的日期会补零
// @Tags 首页
// @Produce json
// @Param days query int false "统计天数，默认7，最大31"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/admin/dashboard/trend [get]
func (h *DashboardHandler) Trend(c *gin.Context) {
	days := 7
	if raw := c.Query("days"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 31 {
			response.BadRequest(c)
			return
		}
		days = parsed
	}

	result, err := h.service.Dashboard.Trend(c.Request.Context(), days)
	if err != nil {
		response.FailErr(c, err)
		return
	}
	response.Ok(c, result)
}
