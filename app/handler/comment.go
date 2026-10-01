package handler

import (
	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service *service.Services
}

func NewCommentHandler(service *service.Services) *CommentHandler {
	return &CommentHandler{service}
}

// ListByArticle godoc
// @Summary 公开-文章评论列表
// @Description 返回审核通过的两级评论树
// @Tags 评论
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/{id}/comment [get]
func (h *CommentHandler) ListByArticle(c *gin.Context) {
	reqDTO, err := dto.ListCommentReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Comment.ListByArticle(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Add godoc
// @Summary 公开-发表评论
// @Description IP限频（每分钟5条），按站点配置决定是否待审核
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body dto.AddCommentRequest true "评论内容"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/{id}/comment [post]
func (h *CommentHandler) Add(c *gin.Context) {
	reqDTO, err := dto.AddCommentReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Comment.Add(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// AdminList godoc
// @Summary 管理-评论列表
// @Description 含各审核状态
// @Tags 评论
// @Produce json
// @Param page query int false "页码(从1起，默认1)"
// @Param page_size query int false "每页数量(默认10，最大100)"
// @Param article_id query int false "文章ID"
// @Param state query int false "状态 0待审核 1通过 2拒绝"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/admin/comment [get]
func (h *CommentHandler) AdminList(c *gin.Context) {
	reqDTO, err := dto.AdminListCommentReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Comment.AdminList(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Audit godoc
// @Summary 管理-评论审核
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Param request body dto.AuditCommentRequest true "审核状态 1通过 2拒绝"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/admin/comment/{id}/state [patch]
func (h *CommentHandler) Audit(c *gin.Context) {
	reqDTO, err := dto.AuditCommentReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Comment.Audit(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Del godoc
// @Summary 管理-删除评论
// @Tags 评论
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/admin/comment/{id} [delete]
func (h *CommentHandler) Del(c *gin.Context) {
	reqDTO, err := dto.IDReqDTOFromRequest(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Comment.Del(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}
