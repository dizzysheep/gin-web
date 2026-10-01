package link

import (
	"context"

	"gin-web/core/xtime"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"gin-web/internal/errcode"
	"gin-web/internal/model"
	"github.com/pkg/errors"
)

type linkService struct {
	daos *dao.Daos
}

func NewLinkService(daos *dao.Daos) LinkService {
	return &linkService{daos}
}

// List 友链列表（全量，前台传 state=1 只返回启用）
func (s *linkService) List(ctx context.Context, reqDTO *dto.ListLinkReqDTO) (*dto.ListLinkRespDTO, error) {
	conditions := common.GormConditions{}
	if reqDTO.Name != "" {
		conditions = append(conditions, &common.LikeCond{Field: "name", Value: reqDTO.Name})
	}
	if reqDTO.State != nil {
		conditions = append(conditions, &common.EqCond{Field: "state", Value: *reqDTO.State})
	}
	conditions = append(conditions, &common.OrderCond{Columns: []string{"sort ASC", "id ASC"}})

	pos, err := s.daos.Link.SelectMany(ctx, conditions, nil)
	if err != nil {
		return nil, errors.Wrap(err, "select link list")
	}
	return &dto.ListLinkRespDTO{LinkPOs: pos}, nil
}

func (s *linkService) Add(ctx context.Context, reqDTO *dto.AddLinkReqDTO) error {
	now := uint32(xtime.GetTimestamp())
	po := &model.Link{
		Name:        reqDTO.Name,
		URL:         reqDTO.URL,
		Logo:        reqDTO.Logo,
		Description: reqDTO.Description,
		Sort:        reqDTO.Sort,
	}
	po.State = reqDTO.State
	po.CreatedBy = reqDTO.Username
	po.ModifiedBy = reqDTO.Username
	po.CreatedOn = now
	po.ModifiedOn = now

	if err := s.daos.Link.InsertOne(ctx, po); err != nil {
		return errors.Wrap(err, "add link fail")
	}
	return nil
}

func (s *linkService) Edit(ctx context.Context, reqDTO *dto.EditLinkReqDTO) error {
	po, err := s.daos.Link.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errcode.NewCustomError(errcode.ErrNoFound)
	}

	po.Name = reqDTO.Name
	po.URL = reqDTO.URL
	po.Logo = reqDTO.Logo
	po.Description = reqDTO.Description
	po.Sort = reqDTO.Sort
	po.State = reqDTO.State
	po.ModifiedBy = reqDTO.Username
	po.ModifiedOn = uint32(xtime.GetTimestamp())

	if err := s.daos.Link.UpdateOne(ctx, po); err != nil {
		return errors.Wrap(err, "edit link fail")
	}
	return nil
}

func (s *linkService) Del(ctx context.Context, reqDTO *dto.IDReqDTO) error {
	if _, err := s.daos.Link.SelectOne(ctx, reqDTO.ID); err != nil {
		return errcode.NewCustomError(errcode.ErrNoFound)
	}

	if err := s.daos.Link.DeleteOne(ctx, reqDTO.ID); err != nil {
		return errors.Wrap(err, "delete link fail")
	}
	return nil
}
