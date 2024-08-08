package tag

import (
	"context"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"github.com/pkg/errors"
	"sync"
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
		errs   []error
		total  int64
		tagPOs []*model.Tag
		wg     sync.WaitGroup
	)

	conditions := common.GormConditions{}
	pager := &common.Pagination{Offset: reqDTO.Offset, PageSize: reqDTO.PageSize}

	wg.Add(1)
	go func() {
		defer wg.Done()
		total, err = s.daos.Tag.Count(ctx, conditions)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "查询总数失败"))
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tagPOs, err = s.daos.Tag.SelectMany(ctx, conditions, pager)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "查询总数失败"))
			return
		}
	}()

	wg.Wait()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	reqDTO.Pager.Total = total
	return &dto.ListTagRespDTO{
		TagPOs: tagPOs,
		Pager:  reqDTO.Pager,
	}, nil
}

// CreateTag ---创建标签tag---
//func (svc *tagService) CreateTag(req dto.CreateTagReqDTO) (int, error) {
//if svc.dao.ExistTagByName(req.Name) {
//	return 0, errors.New("标签已存在")
//}
//
//return svc.dao.CreateTag(&model.Tag{
//	Name:      req.Name,
//	CreatedBy: req.CreatedBy,
//	State:     req.State,
//})
//}

//func (svc *tagService) UpdateTag(tag dto.UpdateTagReqDTO) error {
//	return svc.dao.UpdateInfoById(tag.Id, map[string]interface{}{
//		"name":        tag.Name,
//		"state":       tag.State,
//		"modified_by": tag.ModifiedBy,
//	})
//}
//
//func (svc *tagService) DeleteTag(id int) error {
//	return svc.dao.Delete(id).Error
//
//}
