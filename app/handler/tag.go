package handler

import (
	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	service *service.Services
}

func NewTagHandler(service *service.Services) *TagHandler {
	return &TagHandler{service}
}

// List godoc
// @Summary 标签列表
// @Description 分页返回标签（含关联文章数）
// @Tags 标签
// @Produce json
// @Param page query int false "页码(从1起，默认1)"
// @Param page_size query int false "每页数量(默认10，最大100)"
// @Param name query string false "标签名称(模糊)"
// @Param state query int false "状态 0禁用 1启用"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tag [get]
func (h *TagHandler) List(c *gin.Context) {
	reqDTO, err := dto.ListTagReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}
	if c.FullPath() == "/api/v1/tag" {
		state := int8(1)
		reqDTO.State = &state
	}

	respDTO, err := h.service.Tag.List(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Add godoc
// @Summary 管理-新增标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param request body dto.AddTagRequest true "标签信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tag [post]
func (h *TagHandler) Add(c *gin.Context) {
	reqDTO, err := dto.AddTagReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Tag.Add(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Edit godoc
// @Summary 管理-编辑标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param request body dto.AddTagRequest true "标签信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tag/{id} [patch]
func (h *TagHandler) Edit(c *gin.Context) {
	reqDTO, err := dto.EditTagReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Tag.Edit(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Del godoc
// @Summary 管理-删除标签
// @Description 标签已关联文章时拒绝删除
// @Tags 标签
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tag/{id} [delete]
func (h *TagHandler) Del(c *gin.Context) {
	reqDTO, err := dto.IDReqDTOFromRequest(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Tag.Del(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}
