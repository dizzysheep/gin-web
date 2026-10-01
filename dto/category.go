package dto

import (
	"gin-web/app/ext"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
)

// CategoryVO 分类（文章内嵌展示）
type CategoryVO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func CategoryPOToVO(po *model.Category) *CategoryVO {
	if po == nil {
		return nil
	}
	return &CategoryVO{ID: po.ID, Name: po.Name}
}

// ListCategoryRequest 分类列表查询（全量，不分页传 page_size=0）
type ListCategoryRequest struct {
	Name  string `form:"name" json:"name" binding:"max=100" example:"后端"`
	State *int8  `form:"state" json:"state" binding:"omitempty,oneof=0 1" example:"1"`
}

func ListCategoryReqToDTO(c *gin.Context) (*ListCategoryReqDTO, error) {
	var req ListCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &ListCategoryReqDTO{
		Name:  req.Name,
		State: req.State,
	}, nil
}

type ListCategoryReqDTO struct {
	Name  string
	State *int8
}

type ListCategoryRespDTO struct {
	CategoryPOs   []*model.Category
	ArticleCounts map[int64]int64
}

type ListCategoryResponse struct {
	List []*CategoryItemVO `json:"list"`
}

type CategoryItemVO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ParentID     int64  `json:"parent_id"`
	Sort         int64  `json:"sort"`
	State        int8   `json:"state"`
	ArticleCount int64  `json:"article_count"`
	UpdateTime   uint32 `json:"update_time"`
	UpdateUser   string `json:"update_user"`
}

func (l *ListCategoryRespDTO) ToVO() *ListCategoryResponse {
	list := make([]*CategoryItemVO, 0, len(l.CategoryPOs))
	for _, po := range l.CategoryPOs {
		list = append(list, &CategoryItemVO{
			ID:           po.ID,
			Name:         po.Name,
			ParentID:     po.ParentID,
			Sort:         po.Sort,
			State:        po.State,
			ArticleCount: l.ArticleCounts[po.ID],
			UpdateTime:   po.ModifiedOn,
			UpdateUser:   po.ModifiedBy,
		})
	}
	return &ListCategoryResponse{List: list}
}

// AddCategoryRequest 添加分类
type AddCategoryRequest struct {
	Name     string `json:"name" binding:"required,max=100" example:"后端"`
	ParentID int64  `json:"parent_id" binding:"gte=0" example:"0"`
	Sort     int64  `json:"sort" binding:"gte=0" example:"0"`
	State    *int8  `json:"state" binding:"omitempty,oneof=0 1" example:"1"`
}

func AddCategoryReqToDTO(c *gin.Context) (*AddCategoryReqDTO, error) {
	var req AddCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	state := int8(model.StateEnabled)
	if req.State != nil {
		state = *req.State
	}

	return &AddCategoryReqDTO{
		Name:     req.Name,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		State:    state,
		Username: ext.GetUsername(c),
	}, nil
}

type AddCategoryReqDTO struct {
	Name     string
	ParentID int64
	Sort     int64
	State    int8
	Username string
}

// EditCategoryReqDTO 编辑分类
type EditCategoryReqDTO struct {
	ID       int64
	Name     string
	ParentID int64
	Sort     int64
	State    int8
	Username string
}

func EditCategoryReqToDTO(c *gin.Context) (*EditCategoryReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}

	var req AddCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	state := int8(model.StateEnabled)
	if req.State != nil {
		state = *req.State
	}

	return &EditCategoryReqDTO{
		ID:       id,
		Name:     req.Name,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		State:    state,
		Username: ext.GetUsername(c),
	}, nil
}
