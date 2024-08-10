package tag

import (
	"context"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
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
		err    error
		total  int64
		tagPOs []*model.Tag
		eg     errgroup.Group
	)

	conditions := common.GormConditions{}
	pager := &common.Pagination{Offset: reqDTO.Offset, PageSize: reqDTO.PageSize}

	eg.Go(func() error {
		total, err = s.daos.Tag.Count(ctx, conditions)
		if err != nil {
			return errors.Wrap(err, "select tag count")
		}
		return nil
	})

	eg.Go(func() error {
		tagPOs, err = s.daos.Tag.SelectMany(ctx, conditions, pager)
		if err != nil {
			return errors.Wrap(err, "select tag list")
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	reqDTO.Pager.Total = total
	return &dto.ListTagRespDTO{
		TagPOs: tagPOs,
		Pager:  reqDTO.Pager,
	}, nil
}

func (s *tagService) Add(ctx context.Context, reqDTO *dto.AddTagReqDTO) error {
	tagPO := &model.Tag{
		Name:       reqDTO.Name,
		State:      *reqDTO.State,
		CreateUser: reqDTO.Username,
		UpdateUser: reqDTO.Username,
	}
	if err := s.daos.Tag.InsertOne(ctx, tagPO); err != nil {
		return errors.Wrap(err, "add tag fail")
	}
	return nil
}

func (s *tagService) Edit(ctx context.Context, reqDTO *dto.EditTagReqDTO) error {
	po, err := s.daos.Tag.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errors.Wrap(err, "select tag fail")
	}

	po.Name = reqDTO.Name
	po.State = *reqDTO.State
	po.UpdateUser = reqDTO.Username
	if err := s.daos.Tag.UpdateOne(ctx, po); err != nil {
		return errors.Wrap(err, "edit tag fail")
	}

	return nil
}

func (s *tagService) Del(ctx context.Context, reqDTO *dto.IDReqDTO) error {
	return nil
}
