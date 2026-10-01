package category

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
)

type categoryService struct {
	daos *dao.Daos
}

func NewCategoryService(daos *dao.Daos) CategoryService {
	return &categoryService{daos}
}

// List 分类列表（全量 + 文章数统计），前台传 state=1 只返回启用分类
func (s *categoryService) List(ctx context.Context, reqDTO *dto.ListCategoryReqDTO) (*dto.ListCategoryRespDTO, error) {
	conditions := common.GormConditions{}
	if reqDTO.Name != "" {
		conditions = append(conditions, &common.LikeCond{Field: "name", Value: reqDTO.Name})
	}
	if reqDTO.State != nil {
		conditions = append(conditions, &common.EqCond{Field: "state", Value: *reqDTO.State})
	}
	conditions = append(conditions, &common.OrderCond{Columns: []string{"sort ASC", "id ASC"}})

	categories, err := s.daos.Category.SelectMany(ctx, conditions, nil)
	if err != nil {
		return nil, errors.Wrap(err, "select category list")
	}

	ids := make([]int64, 0, len(categories))
	for _, c := range categories {
		ids = append(ids, c.ID)
	}
	counts, err := s.daos.Category.CountArticleByIDs(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "count category articles")
	}

	return &dto.ListCategoryRespDTO{
		CategoryPOs:   categories,
		ArticleCounts: counts,
	}, nil
}

func (s *categoryService) Add(ctx context.Context, reqDTO *dto.AddCategoryReqDTO) error {
	if err := s.validateParent(ctx, reqDTO.ParentID, 0); err != nil {
		return err
	}
	exist, err := s.daos.Category.ExistName(ctx, reqDTO.Name, 0)
	if err != nil {
		return errors.Wrap(err, "check category name fail")
	}
	if exist {
		return errcode.NewCustomError(errcode.ErrCategoryExist)
	}

	now := uint32(xtime.GetTimestamp())
	po := &model.Category{
		Name:     reqDTO.Name,
		ParentID: reqDTO.ParentID,
		Sort:     reqDTO.Sort,
	}
	po.State = reqDTO.State
	po.CreatedBy = reqDTO.Username
	po.ModifiedBy = reqDTO.Username
	po.CreatedOn = now
	po.ModifiedOn = now

	if err := s.daos.Category.InsertOne(ctx, po); err != nil {
		return errors.Wrap(err, "add category fail")
	}
	return nil
}

func (s *categoryService) Edit(ctx context.Context, reqDTO *dto.EditCategoryReqDTO) error {
	po, err := s.daos.Category.SelectOne(ctx, reqDTO.ID)
	if err != nil {
		return errcode.NewCustomError(errcode.ErrCategoryNotFound)
	}

	exist, err := s.daos.Category.ExistName(ctx, reqDTO.Name, reqDTO.ID)
	if err != nil {
		return errors.Wrap(err, "check category name fail")
	}
	if exist {
		return errcode.NewCustomError(errcode.ErrCategoryExist)
	}
	if err := s.validateParent(ctx, reqDTO.ParentID, reqDTO.ID); err != nil {
		return err
	}

	po.Name = reqDTO.Name
	po.ParentID = reqDTO.ParentID
	po.Sort = reqDTO.Sort
	po.State = reqDTO.State
	po.ModifiedBy = reqDTO.Username
	po.ModifiedOn = uint32(xtime.GetTimestamp())

	if err := s.daos.Category.UpdateOne(ctx, po); err != nil {
		return errors.Wrap(err, "edit category fail")
	}
	articlecache.InvalidateArticleCaches()
	return nil
}

func (s *categoryService) validateParent(ctx context.Context, parentID, selfID int64) error {
	visited := map[int64]struct{}{}
	for parentID > 0 {
		if parentID == selfID {
			return errcode.NewCustomErrorWithMessage(errcode.ErrInvalidParams, "分类层级不能形成循环")
		}
		if _, ok := visited[parentID]; ok {
			return errcode.NewCustomErrorWithMessage(errcode.ErrInvalidParams, "分类层级不能形成循环")
		}
		visited[parentID] = struct{}{}
		parent, err := s.daos.Category.SelectOne(ctx, parentID)
		if err != nil {
			return errcode.NewCustomError(errcode.ErrCategoryNotFound)
		}
		parentID = parent.ParentID
	}
	return nil
}

// Del 删除分类，分类下存在文章时拒绝删除
func (s *categoryService) Del(ctx context.Context, reqDTO *dto.IDReqDTO) error {
	if _, err := s.daos.Category.SelectOne(ctx, reqDTO.ID); err != nil {
		return errcode.NewCustomError(errcode.ErrCategoryNotFound)
	}

	counts, err := s.daos.Category.CountArticleByIDs(ctx, []int64{reqDTO.ID})
	if err != nil {
		return errors.Wrap(err, "count category articles fail")
	}
	if counts[reqDTO.ID] > 0 {
		return errcode.NewCustomError(errcode.ErrCategoryInUse)
	}

	if err := s.daos.Category.DeleteOne(ctx, reqDTO.ID); err != nil {
		return errors.Wrap(err, "delete category fail")
	}
	articlecache.InvalidateArticleCaches()
	return nil
}
