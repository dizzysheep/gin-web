package handler

import (
	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	service *service.Services
}

func NewCommonHandler(service *service.Services) *CommonHandler {
	return &CommonHandler{service}
}

// GetOptions godoc
// @Summary 站点配置
// @Description 返回站点标题/描述/备案号等KV配置（Redis缓存）
// @Tags 配置
// @Produce json
// @Success 200 {object} response.Response
// @Router /v1/common/option [get]
func (h *CommonHandler) GetOptions(c *gin.Context) {
	respDTO, err := h.service.Option.Get(c.Request.Context())
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// SaveOptions godoc
// @Summary 管理-保存站点配置
// @Description 批量保存KV配置并清除缓存
// @Tags 配置
// @Accept json
// @Produce json
// @Param request body dto.SaveOptionsRequest true "配置项列表"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/admin/option [put]
func (h *CommonHandler) SaveOptions(c *gin.Context) {
	reqDTO, err := dto.SaveOptionsReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Option.Save(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}
