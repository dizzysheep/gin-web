package dto

import (
	"errors"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// ListCommentRequest 文章评论列表
type ListCommentRequest struct {
	ArticleID int64 `uri:"id" json:"article_id" example:"1"`
}

func ListCommentReqToDTO(c *gin.Context) (*ListCommentReqDTO, error) {
	id, err := cast.ToInt64E(c.Param("id"))
	if err != nil || id <= 0 {
		return nil, errors.New("不合法的文章 ID")
	}
	return &ListCommentReqDTO{ArticleID: id}, nil
}

type ListCommentReqDTO struct {
	ArticleID int64
}

type ListCommentRespDTO struct {
	CommentPOs []*model.Comment
}

type ListCommentResponse struct {
	List []*CommentVO `json:"list"`
}

// CommentVO 两级评论树
type CommentVO struct {
	ID         int64        `json:"id"`
	ArticleID  int64        `json:"article_id"`
	ParentID   int64        `json:"parent_id"`
	Nickname   string       `json:"nickname"`
	Content    string       `json:"content"`
	IsAdmin    int8         `json:"is_admin"`
	CreateTime uint32       `json:"create_time"`
	Children   []*CommentVO `json:"children"`
}

func (l *ListCommentRespDTO) ToVO() *ListCommentResponse {
	poMap := make(map[int64]*model.Comment, len(l.CommentPOs))
	for _, po := range l.CommentPOs {
		poMap[po.ID] = po
	}

	childrenMap := make(map[int64][]*CommentVO)
	for _, po := range l.CommentPOs {
		if po.ParentID == 0 {
			continue
		}
		// 二级评论统一挂到一级评论下
		rootID := po.ParentID
		if parent, ok := poMap[po.ParentID]; ok && parent.ParentID != 0 {
			rootID = parent.ParentID
		}
		childrenMap[rootID] = append(childrenMap[rootID], commentPOToVO(po))
	}

	list := make([]*CommentVO, 0)
	for _, po := range l.CommentPOs {
		if po.ParentID != 0 {
			continue
		}
		vo := commentPOToVO(po)
		if children, ok := childrenMap[po.ID]; ok {
			vo.Children = children
		} else {
			vo.Children = []*CommentVO{}
		}
		list = append(list, vo)
	}
	return &ListCommentResponse{List: list}
}

func commentPOToVO(po *model.Comment) *CommentVO {
	return &CommentVO{
		ID:         po.ID,
		ArticleID:  po.ArticleID,
		ParentID:   po.ParentID,
		Nickname:   po.Nickname,
		Content:    po.Content,
		IsAdmin:    po.IsAdmin,
		CreateTime: po.CreatedOn,
	}
}

// AddCommentRequest 发表评论
type AddCommentRequest struct {
	ParentID int64  `json:"parent_id" binding:"omitempty,gte=0" example:"0"`
	Nickname string `json:"nickname" binding:"required,max=50" example:"游客"`
	Email    string `json:"email" binding:"omitempty,email,max=100" example:"a@b.com"`
	Content  string `json:"content" binding:"required,max=2000" example:"写得不错"`
}

func AddCommentReqToDTO(c *gin.Context) (*AddCommentReqDTO, error) {
	articleID, err := cast.ToInt64E(c.Param("id"))
	if err != nil || articleID <= 0 {
		return nil, errors.New("不合法的文章 ID")
	}

	var req AddCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &AddCommentReqDTO{
		ArticleID: articleID,
		ParentID:  req.ParentID,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Content:   req.Content,
		IP:        c.ClientIP(),
	}, nil
}

type AddCommentReqDTO struct {
	ArticleID int64
	ParentID  int64
	Nickname  string
	Email     string
	Content   string
	IP        string
}

// AdminListCommentRequest 管理端评论列表
type AdminListCommentRequest struct {
	Page      int   `form:"page" json:"page" binding:"omitempty,gte=1" example:"1"`
	PageSize  int   `form:"page_size" json:"page_size" binding:"omitempty,gte=1,lte=100" example:"10"`
	ArticleID int64 `form:"article_id" json:"article_id" binding:"omitempty,gte=0" example:"1"`
	State     *int8 `form:"state" json:"state" binding:"omitempty,oneof=0 1 2" example:"0"`
}

func AdminListCommentReqToDTO(c *gin.Context) (*AdminListCommentReqDTO, error) {
	var req AdminListCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &AdminListCommentReqDTO{
		ArticleID: req.ArticleID,
		State:     req.State,
		Pager:     PagerReqToDTO(req.Page, req.PageSize),
	}, nil
}

type AdminListCommentReqDTO struct {
	ArticleID int64
	State     *int8
	*Pager
}

type AdminListCommentRespDTO struct {
	Pager      *Pager
	CommentPOs []*model.Comment
}

type AdminListCommentResponse struct {
	Pager *Pager            `json:"pager"`
	List  []*AdminCommentVO `json:"list"`
}

type AdminCommentVO struct {
	ID         int64  `json:"id"`
	ArticleID  int64  `json:"article_id"`
	ParentID   int64  `json:"parent_id"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Content    string `json:"content"`
	IsAdmin    int8   `json:"is_admin"`
	State      int8   `json:"state"`
	IP         string `json:"ip"`
	CreateTime uint32 `json:"create_time"`
}

func (l *AdminListCommentRespDTO) ToVO() *AdminListCommentResponse {
	list := make([]*AdminCommentVO, 0, len(l.CommentPOs))
	for _, po := range l.CommentPOs {
		list = append(list, &AdminCommentVO{
			ID:         po.ID,
			ArticleID:  po.ArticleID,
			ParentID:   po.ParentID,
			Nickname:   po.Nickname,
			Email:      po.Email,
			Content:    po.Content,
			IsAdmin:    po.IsAdmin,
			State:      po.State,
			IP:         po.IP,
			CreateTime: po.CreatedOn,
		})
	}
	return &AdminListCommentResponse{
		Pager: l.Pager,
		List:  list,
	}
}

// AuditCommentRequest 评论审核
type AuditCommentRequest struct {
	State int8 `json:"state" binding:"required,oneof=1 2" example:"1"`
}

func AuditCommentReqToDTO(c *gin.Context) (*AuditCommentReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}

	var req AuditCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &AuditCommentReqDTO{
		ID:    id,
		State: req.State,
	}, nil
}

type AuditCommentReqDTO struct {
	ID    int64
	State int8
}
