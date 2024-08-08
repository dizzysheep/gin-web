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

func (h *TagHandler) List(c *gin.Context) {
	reqDTO, err := dto.ListTagReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Tag.List(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}
