package tag

import (
	"context"

	"gin-web/core/xtime"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"gin-web/internal/errcode"
	"gin-web/internal/model"
	articlecache "gin-web/internal/service/article"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type tagService struct {
	daos *dao.Daos
}

func NewTagService(daos *dao.Daos) TagService {
	return &tagService{daos}
}

func (s *tagService) List(ctx context.Context, reqDTO *dto.ListTagReqDTO) (*dto.ListTagRespDTO, error) {
	var (
		total  int64
		tagPOs []*model.Tag
		eg     errgroup.Group
	)

	conditions := common.GormConditions{
		&common.LikeCond{Field: "name", Value: reqDTO.Name},
		&common.EqCond{Field: "state", Value: reqDTO.State},
	}

	pager := &common.Pagination{Offset: reqDTO.Offset, PageSize: reqDTO.PageSize}

	eg.Go(func() error {
		count, err := s.daos.Tag.Count(ctx, conditions)
		if err != nil {
			return errors.Wrap(err, "select tag count")
		}
		total = count
		return nil
	})

	eg.Go(func() error {
		pos, err := s.daos.Tag.SelectMany(ctx, conditions, pager)
		if err != nil {
			return errors.Wrap(err, "select tag list")
		}
		tagPOs = pos
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(tagPOs))
	for _, t := range tagPOs {
		ids = append(ids, t.ID)
	}
	counts, err := s.daos.Tag.CountArticleByIDs(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "count tag articles")
	}

	reqDTO.Pager.Total = total
	return &dto.ListTagRespDTO{
		TagPOs:        tagPOs,
		ArticleCounts: counts,
		Pager:         reqDTO.Pager,
	}, nil
}

func (s *tagService) Add(ctx context.Context, reqDTO *dto.AddTagReqDTO) error {
	exists, err := s.daos.Tag.ExistName(ctx, reqDTO.Name, 0)
	if err != nil {
		return errors.Wrap(err, "check tag name fail")
	}
	if exists {
		return errcode.NewCustomErrorWithMessage(errcode.ErrInvalidParams, "标签名称已存在")
	}
	tagPO := &model.Tag{Name: reqDTO.Name}
	tagPO.State = *reqDTO.State
	tagPO.CreatedBy = reqDTO.Username
	tagPO.ModifiedBy = reqDTO.Username
	tagPO.CreatedOn = uint32(xtime.GetTimestamp())
	tagPO.ModifiedOn = uint32(xtime.GetTimestamp())
	if err := s.daos.Tag.InsertOne(ctx, tagPO); err != nil {
		return errors.Wrap(err, "add tag fail")
	}
	return nil
}

func (s *tagService) Edit(ctx context.Context, reqDTO *dto.EditTagReqDTO) error {
	po, err := s.daos.Tag.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errcode.NewCustomError(errcode.ErrTagNotFound)
	}
	exists, err := s.daos.Tag.ExistName(ctx, reqDTO.Name, reqDTO.ID)
	if err != nil {
		return errors.Wrap(err, "check tag name fail")
	}
	if exists {
		return errcode.NewCustomErrorWithMessage(errcode.ErrInvalidParams, "标签名称已存在")
	}

	po.Name = reqDTO.Name
	po.State = *reqDTO.State
	po.ModifiedBy = reqDTO.Username
	po.ModifiedOn = uint32(xtime.GetTimestamp())
	if err := s.daos.Tag.UpdateOne(ctx, po); err != nil {
		return errors.Wrap(err, "edit tag fail")
	}
	articlecache.InvalidateArticleCaches()

	return nil
}

// Del 删除标签，标签已关联文章时拒绝删除
func (s *tagService) Del(ctx context.Context, reqDTO *dto.IDReqDTO) error {
	if _, err := s.daos.Tag.SelectOne(ctx, reqDTO.ID); err != nil {
		return errcode.NewCustomError(errcode.ErrTagNotFound)
	}

	count, err := s.daos.Tag.CountRelationsByTagID(ctx, reqDTO.ID)
	if err != nil {
		return errors.Wrap(err, "count tag relations fail")
	}
	if count > 0 {
		return errcode.NewCustomError(errcode.ErrTagInUse)
	}

	if err := s.daos.Tag.DeleteOne(ctx, reqDTO.ID); err != nil {
		return errors.Wrap(err, "delete tag fail")
	}
	articlecache.InvalidateArticleCaches()
	return nil
}
