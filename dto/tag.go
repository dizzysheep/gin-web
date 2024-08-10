package dto

import (
	"gin-web/app/ext"
	"gin-web/core/xtime"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
	"time"
)

// ListTagRequest -----------列表查询------
type ListTagRequest struct {
	PageNo   int    `form:"page_no" json:"page_no"  binding:"required"`
	PageSize int    `form:"page_size" json:"page_size" binding:"required"`
	Name     string `form:"name" json:"name" `
	State    *int8  `form:"state" json:"state"`
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
	State *int8
	*Pager
}

type ListTagRespDTO struct {
	Pager  *Pager
	TagPOs []*model.Tag
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

type ListTagResponse struct {
	Pager *Pager   `json:"pager"`
	List  []*TagVO `json:"list"`
}

// AddTagRequest -----------添加标签------
type AddTagRequest struct {
	Name  string `form:"name" json:"name" binding:"required"`
	State *int8  `form:"state" json:"state" binding:"required"`
}

func AddTagReqToDTO(c *gin.Context) (*AddTagReqDTO, error) {
	var req AddTagRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &AddTagReqDTO{
		Name:     req.Name,
		State:    req.State,
		Username: ext.GetUsername(c),
	}, nil
}

type AddTagReqDTO struct {
	Name     string
	Username string
	State    *int8
}

// EditTagReqDTO -----------添加标签------
type EditTagReqDTO struct {
	ID       int64
	State    *int8
	Name     string
	Username string
}

func EditTagReqToDTO(c *gin.Context) (*EditTagReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}

	var req AddTagRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &EditTagReqDTO{
		ID:       id,
		Name:     req.Name,
		State:    req.State,
		Username: ext.GetUsername(c),
	}, nil
}

// TagVO -----------试图层显示------
type TagVO struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	State      int8   `json:"state"`
	UpdateTime string `json:"update_time"`
	UpdateUser string `json:"update_user"`
}

func TagPOToVO(po *model.Tag) *TagVO {
	if po == nil {
		return nil
	}

	return &TagVO{
		ID:         po.ID,
		Name:       po.Name,
		State:      po.State,
		UpdateTime: time.Unix(po.UpdatedAt, 0).Format(xtime.DATE_TIME_FMT),
		UpdateUser: po.UpdateUser,
	}
}
