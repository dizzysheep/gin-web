package handler

import (
	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	service *service.Services
}

func NewLinkHandler(service *service.Services) *LinkHandler {
	return &LinkHandler{service}
}

// List godoc
// @Summary 友情链接列表
// @Description 全量友链，前台可传 state=1 只看启用
// @Tags 友链
// @Produce json
// @Param name query string false "站点名称(模糊)"
// @Param state query int false "状态 0禁用 1启用"
// @Success 200 {object} response.Response
// @Router /v1/link [get]
func (h *LinkHandler) List(c *gin.Context) {
	reqDTO, err := dto.ListLinkReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}
	if c.FullPath() == "/api/v1/link" {
		state := int8(1)
		reqDTO.State = &state
	}

	respDTO, err := h.service.Link.List(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Add godoc
// @Summary 管理-新增友链
// @Tags 友链
// @Accept json
// @Produce json
// @Param request body dto.AddLinkRequest true "友链信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/link [post]
func (h *LinkHandler) Add(c *gin.Context) {
	reqDTO, err := dto.AddLinkReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Link.Add(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Edit godoc
// @Summary 管理-编辑友链
// @Tags 友链
// @Accept json
// @Produce json
// @Param id path int true "友链ID"
// @Param request body dto.AddLinkRequest true "友链信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/link/{id} [patch]
func (h *LinkHandler) Edit(c *gin.Context) {
	reqDTO, err := dto.EditLinkReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Link.Edit(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Del godoc
// @Summary 管理-删除友链
// @Tags 友链
// @Produce json
// @Param id path int true "友链ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/link/{id} [delete]
func (h *LinkHandler) Del(c *gin.Context) {
	reqDTO, err := dto.IDReqDTOFromRequest(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Link.Del(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}
