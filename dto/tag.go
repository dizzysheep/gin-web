package dto

import (
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
)

type CreateTagReqDTO struct {
	Name      string `json:"name" binding:"required"`
	State     int    `json:"state"  binding:"gte=0"`
	CreatedBy string `json:"created_by"  binding:"required"`
}

type UpdateTagReqDTO struct {
	Id         int    `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	State      int    `json:"state"  binding:"gte=0"`
	ModifiedBy string `json:"modified_by"  binding:"required"`
}

// ListTagRequest -----------列表查询------
type ListTagRequest struct {
	PageNo   int    `form:"page_no"  binding:"required"`
	PageSize int    `form:"page_size"  binding:"required"`
	Name     string `form:"name" `
	State    *int   `form:"state"`
}

func ListTagReqToDTO(c *gin.Context) (*ListTagReqDTO, error) {
	var req ListTagRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &ListTagReqDTO{
		Name:  req.Name,
		State: req.State,
		Pager: PagerReqToDTO(req.PageNo, req.PageSize),
	}, nil
}

type ListTagReqDTO struct {
	Name  string
	State *int
	*Pager
}

type ListTagRespDTO struct {
	Pager  *Pager
	TagPOs []*model.Tag
}

type ListTagResponse struct {
	Pager *Pager
	List  []*TagVO
}

func (l *ListTagRespDTO) ToVO() *ListTagResponse {
	list := make([]*TagVO, 0, len(l.TagPOs))
	for _, po := range l.TagPOs {
		list = append(list, TagPOToVO(po))
	}
	return &ListTagResponse{
		Pager: l.Pager,
		List:  list,
	}
}

type TagVO struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	State int8   `json:"state"`
}

func TagPOToVO(po *model.Tag) *TagVO {
	return &TagVO{
		ID:    po.ID,
		Name:  po.Name,
		State: po.State,
	}
}
