package article

import (
	"context"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type articleService struct {
	daos *dao.Daos
}

func NewArticleService(daos *dao.Daos) ArticleService {
	return &articleService{daos}
}

func (s *articleService) List(ctx context.Context, reqDTO *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error) {
	var (
		total      int64
		err        error
		articlePOs []*model.Article
		eg         errgroup.Group
	)

	conditions := common.GormConditions{}
	pager := &common.Pagination{Offset: reqDTO.Offset, PageSize: reqDTO.PageSize}

	eg.Go(func() error {
		total, err = s.daos.Article.Count(ctx, conditions)
		if err != nil {
			return errors.Wrap(err, "select article count")
		}
		return nil
	})

	eg.Go(func() error {
		articlePOs, err = s.daos.Article.SelectMany(ctx, conditions, pager)
		if err != nil {
			return errors.Wrap(err, "select article list")
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	reqDTO.Pager.Total = total
	return &dto.ListArticleRespDTO{
		POs:   articlePOs,
		Pager: reqDTO.Pager,
	}, nil
}

func (s *articleService) Add(ctx context.Context, reqDTO *dto.AddArticleReqDTO) error {
	po := &model.Article{
		TagID:      reqDTO.TagID,
		Title:      reqDTO.Title,
		Desc:       reqDTO.Desc,
		Content:    reqDTO.Content,
		State:      reqDTO.State,
		CreateUser: reqDTO.Username,
		UpdateUser: reqDTO.Username,
	}

	if err := s.daos.Article.InsertOne(ctx, po); err != nil {
		return errors.Wrap(err, "add tag fail")
	}

	return nil
}
