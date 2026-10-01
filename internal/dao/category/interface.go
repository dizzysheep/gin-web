package category

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
)

type Reader interface {
	Count(ctx context.Context, conditions common.GormConditions) (int64, error)
	SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Category, error)
	SelectOne(ctx context.Context, id int64) (*model.Category, error)
	SelectManyByIDs(ctx context.Context, ids []int64) ([]*model.Category, error)
	ExistName(ctx context.Context, name string, excludeID int64) (bool, error)
	CountArticleByIDs(ctx context.Context, ids []int64) (map[int64]int64, error)
}

type Writer interface {
	InsertOne(ctx context.Context, category *model.Category) error
	UpdateOne(ctx context.Context, category *model.Category) error
	DeleteOne(ctx context.Context, id int64) error
}

type CategoryDao interface {
	Reader
	Writer
}
