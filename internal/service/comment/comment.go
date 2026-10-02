package comment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gin-web/core/redis"
	"gin-web/core/xtime"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"gin-web/internal/errcode"
	"gin-web/internal/model"
	"gin-web/internal/service/option"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	// rateLimitKeyPrefix 评论IP限频key
	rateLimitKeyPrefix = "blog:comment:limit:"
	// rateLimitWindow 限频窗口
	rateLimitWindow = time.Minute
	// rateLimitMax 窗口内最大评论数
	rateLimitMax = 5
)

type commentService struct {
	daos   *dao.Daos
	option option.OptionService
}

func NewCommentService(daos *dao.Daos, optionSvc option.OptionService) CommentService {
	return &commentService{daos: daos, option: optionSvc}
}

// ListByArticle 文章评论列表（仅审核通过，两级树）
func (s *commentService) ListByArticle(ctx context.Context, reqDTO *dto.ListCommentReqDTO) (*dto.ListCommentRespDTO, error) {
	conditions := common.GormConditions{
		&common.EqCond{Field: "article_id", Value: reqDTO.ArticleID},
		&common.EqCond{Field: "state", Value: model.CommentStatePass},
		&common.OrderCond{Columns: []string{"created_on ASC", "id ASC"}},
	}

	pos, err := s.daos.Comment.SelectMany(ctx, conditions, nil)
	if err != nil {
		return nil, errors.Wrap(err, "select comment list")
	}

	return &dto.ListCommentRespDTO{CommentPOs: pos}, nil
}

// Add 发表评论（IP限频 + 按配置决定是否待审核）
func (s *commentService) Add(ctx context.Context, reqDTO *dto.AddCommentReqDTO) error {
	if !s.allowByIP(reqDTO.IP) {
		return errcode.NewCustomError(errcode.ErrCommentTooFast)
	}

	article, err := s.daos.Article.SelectOne(ctx, reqDTO.ArticleID)
	if err != nil || article.IsDraft != 0 || article.State != model.StateEnabled {
		return errcode.NewCustomError(errcode.ErrArticleNotFound)
	}

	if reqDTO.ParentID > 0 {
		parent, err := s.daos.Comment.SelectOne(ctx, reqDTO.ParentID)
		if err != nil || parent.ArticleID != reqDTO.ArticleID {
			return errcode.NewCustomError(errcode.ErrCommentNotFound)
		}
		if parent.ParentID > 0 {
			reqDTO.ParentID = parent.ParentID
		}
	}

	state := model.CommentStatePass
	if strings.TrimSpace(s.option.Value(ctx, "comment_audit")) == "1" {
		state = model.CommentStatePending
	}

	po := &model.Comment{
		ArticleID: reqDTO.ArticleID,
		ParentID:  reqDTO.ParentID,
		Nickname:  reqDTO.Nickname,
		Email:     reqDTO.Email,
		Content:   reqDTO.Content,
		IsAdmin:   0,
		IP:        reqDTO.IP,
	}
	po.State = state
	po.CreatedOn = uint32(xtime.GetTimestamp())
	po.ModifiedOn = uint32(xtime.GetTimestamp())
	po.CreatedBy = reqDTO.Nickname
	po.ModifiedBy = reqDTO.Nickname

	if err := s.daos.Comment.InsertOne(ctx, po); err != nil {
		return errors.Wrap(err, "add comment fail")
	}
	return nil
}

// AdminList 管理端评论列表（含各状态）
func (s *commentService) AdminList(ctx context.Context, reqDTO *dto.AdminListCommentReqDTO) (*dto.AdminListCommentRespDTO, error) {
	var (
		total      int64
		commentPOs []*model.Comment
		eg         errgroup.Group
	)

	conditions := common.GormConditions{
		&common.EqCond{Field: "state", Value: reqDTO.State},
	}
	if reqDTO.ArticleID > 0 {
		conditions = append(conditions, &common.EqCond{Field: "article_id", Value: reqDTO.ArticleID})
	}

	pager := &common.Pagination{Offset: reqDTO.Offset, PageSize: reqDTO.PageSize}

	eg.Go(func() error {
		count, err := s.daos.Comment.Count(ctx, conditions)
		if err != nil {
			return errors.Wrap(err, "select comment count")
		}
		total = count
		return nil
	})

	eg.Go(func() error {
		pos, err := s.daos.Comment.SelectMany(ctx, conditions, pager)
		if err != nil {
			return errors.Wrap(err, "select comment list")
		}
		commentPOs = pos
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	reqDTO.Pager.Total = total
	return &dto.AdminListCommentRespDTO{
		CommentPOs: commentPOs,
		Pager:      reqDTO.Pager,
	}, nil
}

// Audit 评论审核
func (s *commentService) Audit(ctx context.Context, reqDTO *dto.AuditCommentReqDTO) error {
	if _, err := s.daos.Comment.SelectOne(ctx, reqDTO.ID); err != nil {
		return errcode.NewCustomError(errcode.ErrCommentNotFound)
	}

	cols := map[string]interface{}{
		"state":       reqDTO.State,
		"modified_on": xtime.GetTimestamp(),
	}
	if err := s.daos.Comment.UpdateColumns(ctx, reqDTO.ID, cols); err != nil {
		return errors.Wrap(err, "audit comment fail")
	}
	return nil
}

// Del 删除评论（软删）
func (s *commentService) Del(ctx context.Context, reqDTO *dto.IDReqDTO) error {
	if _, err := s.daos.Comment.SelectOne(ctx, reqDTO.ID); err != nil {
		return errcode.NewCustomError(errcode.ErrCommentNotFound)
	}

	if err := s.daos.Comment.DeleteOne(ctx, reqDTO.ID); err != nil {
		return errors.Wrap(err, "delete comment fail")
	}
	return nil
}

// allowByIP IP限频：窗口内超过阈值返回false，Redis异常时放行
func (s *commentService) allowByIP(ip string) bool {
	if ip == "" || redis.RedisClient == nil {
		return true
	}

	ctx := context.Background()
	key := fmt.Sprintf("%s%s", rateLimitKeyPrefix, ip)
	n, err := redis.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return true
	}
	if n == 1 {
		redis.RedisClient.Expire(ctx, key, rateLimitWindow)
	}
	return n <= rateLimitMax
}
