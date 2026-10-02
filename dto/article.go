package dto

import (
	"errors"

	"gin-web/app/ext"
	"gin-web/internal/model"
	"gin-web/internal/utils"
	"github.com/gin-gonic/gin"
)

// ListArticleRequest 公开文章列表查询
type ListArticleRequest struct {
	Page       int    `form:"page" json:"page" binding:"omitempty,gte=1" example:"1"`
	PageSize   int    `form:"page_size" json:"page_size" binding:"omitempty,gte=1,lte=100" example:"10"`
	CategoryID int64  `form:"category_id" json:"category_id" binding:"gte=0" example:"1"`
	TagID      int64  `form:"tag_id" json:"tag_id" binding:"gte=0" example:"1"`
	Keyword    string `form:"keyword" json:"keyword" binding:"max=100" example:"gin"`
}

func ListArticleReqToDTO(c *gin.Context) (*ListArticleReqDTO, error) {
	var req ListArticleRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &ListArticleReqDTO{
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
		Keyword:    req.Keyword,
		Pager:      PagerReqToDTO(req.Page, req.PageSize),
	}, nil
}

type ListArticleReqDTO struct {
	CategoryID int64
	TagID      int64
	Keyword    string
	*Pager
}

// AdminListArticleRequest 管理端文章列表查询（含草稿）
type AdminListArticleRequest struct {
	ListArticleRequest
	State         *int8  `form:"state" json:"state" binding:"omitempty,oneof=0 1" example:"1"`
	IsDraft       *int8  `form:"is_draft" json:"is_draft" binding:"omitempty,oneof=0 1" example:"0"`
	PublishedFrom string `form:"published_from" json:"published_from" example:"2026-10-01"`
	PublishedTo   string `form:"published_to" json:"published_to" example:"2026-10-02"`
}

func AdminListArticleReqToDTO(c *gin.Context) (*AdminListArticleReqDTO, error) {
	var req AdminListArticleRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}
	publishedFrom, err := utils.ParseDateToTimestamp(req.PublishedFrom, false)
	if err != nil {
		return nil, err
	}
	publishedTo, err := utils.ParseDateToTimestamp(req.PublishedTo, true)
	if err != nil {
		return nil, err
	}
	if publishedFrom != nil && publishedTo != nil && *publishedFrom >= *publishedTo {
		return nil, errors.New("发布时间起始日期不能晚于结束日期")
	}

	return &AdminListArticleReqDTO{
		ListArticleReqDTO: ListArticleReqDTO{
			CategoryID: req.CategoryID,
			TagID:      req.TagID,
			Keyword:    req.Keyword,
			Pager:      PagerReqToDTO(req.Page, req.PageSize),
		},
		State:         req.State,
		IsDraft:       req.IsDraft,
		PublishedFrom: publishedFrom,
		PublishedTo:   publishedTo,
	}, nil
}

type AdminListArticleReqDTO struct {
	ListArticleReqDTO
	State         *int8
	IsDraft       *int8
	PublishedFrom *int64
	PublishedTo   *int64
}

// ArticleItem 列表项（PO + 聚合的标签/分类）
type ArticleItem struct {
	PO       *model.Article
	Category *model.Category
	Tags     []*model.Tag
}

type ListArticleRespDTO struct {
	Pager *Pager
	Items []*ArticleItem
}

type ListArticleResponse struct {
	Pager *Pager           `json:"pager"`
	List  []*ArticleListVO `json:"list"`
}

type ArticleListVO struct {
	ID          int64       `json:"id"`
	CategoryID  int64       `json:"category_id"`
	Category    *CategoryVO `json:"category"`
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Desc        string      `json:"desc"`
	Cover       string      `json:"cover"`
	ViewCount   int64       `json:"view_count"`
	IsTop       int8        `json:"is_top"`
	IsDraft     int8        `json:"is_draft"`
	State       int8        `json:"state"`
	PublishedOn uint32      `json:"published_on"`
	Tags        []*TagVO    `json:"tags"`
}

func ArticleItemToVO(item *ArticleItem) *ArticleListVO {
	if item == nil || item.PO == nil {
		return nil
	}
	po := item.PO
	vo := &ArticleListVO{
		ID:          po.ID,
		CategoryID:  po.CategoryID,
		Title:       po.Title,
		Slug:        po.Slug,
		Desc:        po.Desc,
		Cover:       po.Cover,
		ViewCount:   po.ViewCount,
		IsTop:       po.IsTop,
		IsDraft:     po.IsDraft,
		State:       po.State,
		PublishedOn: po.PublishedOn,
		Category:    CategoryPOToVO(item.Category),
		Tags:        make([]*TagVO, 0, len(item.Tags)),
	}
	for _, tag := range item.Tags {
		vo.Tags = append(vo.Tags, TagPOToVO(tag))
	}
	return vo
}

func (l *ListArticleRespDTO) ToVO() *ListArticleResponse {
	list := make([]*ArticleListVO, 0, len(l.Items))
	for _, item := range l.Items {
		list = append(list, ArticleItemToVO(item))
	}
	return &ListArticleResponse{
		Pager: l.Pager,
		List:  list,
	}
}

// ArticleDetailRespDTO 文章详情
type ArticleDetailRespDTO struct {
	PO        *model.Article
	Category  *model.Category
	Tags      []*model.Tag
	ViewDelta int64 // 当日未落库的浏览量增量
}

type ArticleDetailResponse struct {
	ID          int64       `json:"id"`
	CategoryID  int64       `json:"category_id"`
	Category    *CategoryVO `json:"category"`
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Desc        string      `json:"desc"`
	Cover       string      `json:"cover"`
	ContentMd   string      `json:"content_md"`
	ContentHtml string      `json:"content_html"`
	ViewCount   int64       `json:"view_count"`
	IsTop       int8        `json:"is_top"`
	IsDraft     int8        `json:"is_draft"`
	State       int8        `json:"state"`
	PublishedOn uint32      `json:"published_on"`
	Tags        []*TagVO    `json:"tags"`
}

func (d *ArticleDetailRespDTO) ToVO() *ArticleDetailResponse {
	if d.PO == nil {
		return &ArticleDetailResponse{}
	}
	po := d.PO
	vo := &ArticleDetailResponse{
		ID:          po.ID,
		CategoryID:  po.CategoryID,
		Category:    CategoryPOToVO(d.Category),
		Title:       po.Title,
		Slug:        po.Slug,
		Desc:        po.Desc,
		Cover:       po.Cover,
		ContentMd:   po.ContentMd,
		ContentHtml: po.ContentHtml,
		ViewCount:   po.ViewCount + d.ViewDelta,
		IsTop:       po.IsTop,
		IsDraft:     po.IsDraft,
		State:       po.State,
		PublishedOn: po.PublishedOn,
		Tags:        make([]*TagVO, 0, len(d.Tags)),
	}
	for _, tag := range d.Tags {
		vo.Tags = append(vo.Tags, TagPOToVO(tag))
	}
	return vo
}

// AddArticleRequest 新建文章
type AddArticleRequest struct {
	Title      string  `json:"title" binding:"required,max=100" example:"文章标题"`
	Desc       string  `json:"desc" binding:"max=255" example:"文章简述"`
	ContentMd  string  `json:"content_md" binding:"required" example:"# markdown"`
	Cover      string  `json:"cover" binding:"omitempty,max=255" example:"/uploads/20261001/xx.png"`
	CategoryID int64   `json:"category_id" binding:"gte=0" example:"1"`
	TagIDs     []int64 `json:"tag_ids" binding:"omitempty,max=10,dive,gt=0" example:"1,2"`
	Slug       string  `json:"slug" binding:"omitempty,max=150" example:"hello-world"`
	IsDraft    bool    `json:"is_draft" example:"false"`
}

func AddArticleReqToDTO(c *gin.Context) (*AddArticleReqDTO, error) {
	var req AddArticleRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &AddArticleReqDTO{
		Title:      req.Title,
		Desc:       req.Desc,
		ContentMd:  req.ContentMd,
		Cover:      req.Cover,
		CategoryID: req.CategoryID,
		TagIDs:     req.TagIDs,
		Slug:       req.Slug,
		IsDraft:    req.IsDraft,
		Username:   ext.GetUsername(c),
	}, nil
}

type AddArticleReqDTO struct {
	Title      string
	Desc       string
	ContentMd  string
	Cover      string
	CategoryID int64
	TagIDs     []int64
	Slug       string
	IsDraft    bool
	Username   string
}

// EditArticleRequest 编辑文章
type EditArticleRequest struct {
	AddArticleRequest
}

func EditArticleReqToDTO(c *gin.Context) (*EditArticleReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}

	var req EditArticleRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &EditArticleReqDTO{
		ID: id,
		AddArticleReqDTO: AddArticleReqDTO{
			Title:      req.Title,
			Desc:       req.Desc,
			ContentMd:  req.ContentMd,
			Cover:      req.Cover,
			CategoryID: req.CategoryID,
			TagIDs:     req.TagIDs,
			Slug:       req.Slug,
			IsDraft:    req.IsDraft,
			Username:   ext.GetUsername(c),
		},
	}, nil
}

type EditArticleReqDTO struct {
	ID int64
	AddArticleReqDTO
}

// PublishArticleRequest 发布/下架
type PublishArticleRequest struct {
	Publish *bool `json:"publish" binding:"required" example:"true"`
}

func PublishArticleReqToDTO(c *gin.Context) (*PublishArticleReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}

	var req PublishArticleRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &PublishArticleReqDTO{
		ID:       id,
		Publish:  *req.Publish,
		Username: ext.GetUsername(c),
	}, nil
}

type PublishArticleReqDTO struct {
	ID       int64
	Publish  bool
	Username string
}

// TopArticleRequest 置顶/取消置顶
type TopArticleRequest struct {
	IsTop *int8 `json:"is_top" binding:"required,oneof=0 1" example:"1"`
}

func TopArticleReqToDTO(c *gin.Context) (*TopArticleReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}

	var req TopArticleRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &TopArticleReqDTO{
		ID:       id,
		IsTop:    *req.IsTop,
		Username: ext.GetUsername(c),
	}, nil
}

type TopArticleReqDTO struct {
	ID       int64
	IsTop    int8
	Username string
}

// ArticleSlugReqDTO 按slug查询
type ArticleSlugReqDTO struct {
	Slug string
}

// ArchiveRespDTO 归档
type ArchiveRespDTO struct {
	Rows []*model.ArchiveRow
}

type ArchiveResponse struct {
	List []*ArchiveVO `json:"list"`
}

type ArchiveVO struct {
	Month string `json:"month" example:"2026-09"`
	Total int64  `json:"total" example:"3"`
}

func (a *ArchiveRespDTO) ToVO() *ArchiveResponse {
	list := make([]*ArchiveVO, 0, len(a.Rows))
	for _, row := range a.Rows {
		list = append(list, &ArchiveVO{Month: row.Month, Total: row.Total})
	}
	return &ArchiveResponse{List: list}
}
