package dto

import (
	"gin-web/app/ext"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
)

// ListLinkRequest 友链列表（全量）
type ListLinkRequest struct {
	Name  string `form:"name" json:"name" binding:"max=100" example:"golang"`
	State *int8  `form:"state" json:"state" binding:"omitempty,oneof=0 1" example:"1"`
}

func ListLinkReqToDTO(c *gin.Context) (*ListLinkReqDTO, error) {
	var req ListLinkRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}
	return &ListLinkReqDTO{
		Name:  req.Name,
		State: req.State,
	}, nil
}

type ListLinkReqDTO struct {
	Name  string
	State *int8
}

type ListLinkRespDTO struct {
	LinkPOs []*model.Link
}

type ListLinkResponse struct {
	List []*LinkVO `json:"list"`
}

type LinkVO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Sort        int64  `json:"sort"`
	State       int8   `json:"state"`
	UpdateTime  uint32 `json:"update_time"`
	UpdateUser  string `json:"update_user"`
}

func (l *ListLinkRespDTO) ToVO() *ListLinkResponse {
	list := make([]*LinkVO, 0, len(l.LinkPOs))
	for _, po := range l.LinkPOs {
		list = append(list, &LinkVO{
			ID:          po.ID,
			Name:        po.Name,
			URL:         po.URL,
			Logo:        po.Logo,
			Description: po.Description,
			Sort:        po.Sort,
			State:       po.State,
			UpdateTime:  po.ModifiedOn,
			UpdateUser:  po.ModifiedBy,
		})
	}
	return &ListLinkResponse{List: list}
}

// AddLinkRequest 添加友链
type AddLinkRequest struct {
	Name        string `json:"name" binding:"required,max=100" example:"golang"`
	URL         string `json:"url" binding:"required,url,max=255" example:"https://go.dev"`
	Logo        string `json:"logo" binding:"omitempty,max=255"`
	Description string `json:"description" binding:"omitempty,max=255"`
	Sort        int64  `json:"sort" binding:"gte=0" example:"0"`
	State       *int8  `json:"state" binding:"omitempty,oneof=0 1" example:"1"`
}

func AddLinkReqToDTO(c *gin.Context) (*AddLinkReqDTO, error) {
	var req AddLinkRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	state := int8(model.StateEnabled)
	if req.State != nil {
		state = *req.State
	}

	return &AddLinkReqDTO{
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Sort:        req.Sort,
		State:       state,
		Username:    ext.GetUsername(c),
	}, nil
}

type AddLinkReqDTO struct {
	Name        string
	URL         string
	Logo        string
	Description string
	Sort        int64
	State       int8
	Username    string
}

// EditLinkReqDTO 编辑友链
type EditLinkReqDTO struct {
	ID          int64
	Name        string
	URL         string
	Logo        string
	Description string
	Sort        int64
	State       int8
	Username    string
}

func EditLinkReqToDTO(c *gin.Context) (*EditLinkReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}

	var req AddLinkRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	state := int8(model.StateEnabled)
	if req.State != nil {
		state = *req.State
	}

	return &EditLinkReqDTO{
		ID:          id,
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Sort:        req.Sort,
		State:       state,
		Username:    ext.GetUsername(c),
	}, nil
}
