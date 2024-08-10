package dto

import (
	"gin-web/app/ext"
	"gin-web/core/xtime"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
	"time"
)

// ListArticleRequest -----------列表查询------
type ListArticleRequest struct {
	PageNo   int `form:"page_no" json:"page_no"  binding:"required"`
	PageSize int `form:"page_size" json:"page_size"  binding:"required"`
	//Name     string `form:"name" `
	//State    *int8  `form:"state"`
}

func ListArticleReqToDTO(c *gin.Context) (*ListArticleReqDTO, error) {
	var req ListArticleRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &ListArticleReqDTO{
		Pager: PagerReqToDTO(req.PageNo, req.PageSize),
	}, nil
}

type ListArticleReqDTO struct {
	Name  string
	State *int8
	*Pager
}

type ListArticleRespDTO struct {
	Pager *Pager
	POs   []*model.Article
}

func (l *ListArticleRespDTO) ToVO() *ListArticleResponse {
	list := make([]*ArticleVO, 0, len(l.POs))
	for _, po := range l.POs {
		list = append(list, ArticlePOToVO(po))
	}
	return &ListArticleResponse{
		Pager: l.Pager,
		List:  list,
	}
}

type ListArticleResponse struct {
	Pager *Pager       `json:"pager"`
	List  []*ArticleVO `json:"list"`
}

// ArticleVO -----------VO------
type ArticleVO struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	State      int8   `json:"state"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	Tag        *TagVO `json:"tag"`
	UpdateUser string `json:"update_user"`
	UpdateTime string `json:"update_time"`
}

func ArticlePOToVO(po *model.Article) *ArticleVO {
	if po == nil {
		return nil
	}
	return &ArticleVO{
		ID:         po.ID,
		Title:      po.Title,
		State:      po.State,
		Desc:       po.Desc,
		Content:    po.Content,
		Tag:        TagPOToVO(po.Tag),
		UpdateUser: po.UpdateUser,
		UpdateTime: time.Unix(po.UpdatedAt, 0).Format(xtime.DATE_TIME_FMT),
	}
}

// AddArticleRequest -----------列表查询------
type AddArticleRequest struct {
	TagID   int64  `json:"tag_id" binding:"required,gt=0"`
	Title   string `json:"title" binding:"required"`
	Desc    string `json:"desc" binding:"required"`
	Content string `json:"content" binding:"required"`
	State   int8   `json:"state" binding:"gte=0"`
}

func AddArticleReqToDTO(c *gin.Context) (*AddArticleReqDTO, error) {
	var req AddArticleRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &AddArticleReqDTO{
		TagID:    req.TagID,
		Title:    req.Title,
		Desc:     req.Desc,
		Content:  req.Content,
		State:    req.State,
		Username: ext.GetUsername(c),
	}, nil
}

type AddArticleReqDTO struct {
	TagID    int64
	Title    string
	Desc     string
	Content  string
	State    int8
	Username string
}
