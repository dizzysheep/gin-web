package dto

import (
	"gin-web/app/ext"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
)

// TagVO 标签（文章内嵌展示）
type TagVO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func TagPOToVO(po *model.Tag) *TagVO {
	if po == nil {
		return nil
	}
	return &TagVO{ID: po.ID, Name: po.Name}
}

// ListTagRequest 标签列表查询
type ListTagRequest struct {
	Page     int    `form:"page" json:"page" binding:"omitempty,gte=1" example:"1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"omitempty,gte=1,lte=100" example:"10"`
	Name     string `form:"name" json:"name" binding:"max=100" example:"go"`
	State    *int8  `form:"state" json:"state" binding:"omitempty,oneof=0 1" example:"1"`
}

func ListTagReqToDTO(c *gin.Context) (*ListTagReqDTO, error) {
	var req ListTagRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &ListTagReqDTO{
		Name:  req.Name,
		State: req.State,
		Pager: PagerReqToDTO(req.Page, req.PageSize),
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
	// ArticleCounts 标签ID -> 关联文章数
	ArticleCounts map[int64]int64
}

type ListTagResponse struct {
	Pager *Pager       `json:"pager"`
	List  []*TagItemVO `json:"list"`
}

type TagItemVO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	State        int8   `json:"state"`
	ArticleCount int64  `json:"article_count"`
	UpdateTime   uint32 `json:"update_time"`
	UpdateUser   string `json:"update_user"`
}

func (l *ListTagRespDTO) ToVO() *ListTagResponse {
	list := make([]*TagItemVO, 0, len(l.TagPOs))
	for _, po := range l.TagPOs {
		list = append(list, &TagItemVO{
			ID:           po.ID,
			Name:         po.Name,
			State:        po.State,
			ArticleCount: l.ArticleCounts[po.ID],
			UpdateTime:   po.ModifiedOn,
			UpdateUser:   po.ModifiedBy,
		})
	}
	return &ListTagResponse{
		Pager: l.Pager,
		List:  list,
	}
}

// AddTagRequest 添加标签
type AddTagRequest struct {
	Name  string `form:"name" json:"name" binding:"required,max=100" example:"go"`
	State *int8  `form:"state" json:"state" binding:"required,oneof=0 1" example:"1"`
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
	State    *int8
	Username string
}

// EditTagReqDTO 编辑标签
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
