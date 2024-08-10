package handler

import (
	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service *service.Services
}

func NewArticleHandler(service *service.Services) *ArticleHandler {
	return &ArticleHandler{service}
}

func (h *ArticleHandler) List(c *gin.Context) {
	reqDTO, err := dto.ListArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Article.List(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

func (h *ArticleHandler) Add(c *gin.Context) {
	reqDTO, err := dto.AddArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Article.Add(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}
