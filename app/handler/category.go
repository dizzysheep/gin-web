package handler

import (
	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.Services
}

func NewCategoryHandler(service *service.Services) *CategoryHandler {
	return &CategoryHandler{service}
}

// List godoc
// @Summary 分类列表
// @Description 全量分类列表（含文章数），前台可传 state=1 只看启用
// @Tags 分类
// @Produce json
// @Param name query string false "分类名称(模糊)"
// @Param state query int false "状态 0禁用 1启用"
// @Success 200 {object} response.Response
// @Router /v1/category [get]
func (h *CategoryHandler) List(c *gin.Context) {
	reqDTO, err := dto.ListCategoryReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}
	if c.FullPath() == "/api/v1/category" {
		state := int8(1)
		reqDTO.State = &state
	}

	respDTO, err := h.service.Category.List(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Add godoc
// @Summary 管理-新增分类
// @Tags 分类
// @Accept json
// @Produce json
// @Param request body dto.AddCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/category [post]
func (h *CategoryHandler) Add(c *gin.Context) {
	reqDTO, err := dto.AddCategoryReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Category.Add(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Edit godoc
// @Summary 管理-编辑分类
// @Tags 分类
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param request body dto.AddCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/category/{id} [patch]
func (h *CategoryHandler) Edit(c *gin.Context) {
	reqDTO, err := dto.EditCategoryReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Category.Edit(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Del godoc
// @Summary 管理-删除分类
// @Description 分类下存在文章时拒绝删除
// @Tags 分类
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/category/{id} [delete]
func (h *CategoryHandler) Del(c *gin.Context) {
	reqDTO, err := dto.IDReqDTOFromRequest(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Category.Del(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}
