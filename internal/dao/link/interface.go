package link

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
)

type Reader interface {
	Count(ctx context.Context, conditions common.GormConditions) (int64, error)
	SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Link, error)
	SelectOne(ctx context.Context, id int64) (*model.Link, error)
}

type Writer interface {
	InsertOne(ctx context.Context, link *model.Link) error
	UpdateOne(ctx context.Context, link *model.Link) error
	DeleteOne(ctx context.Context, id int64) error
}

type LinkDao interface {
	Reader
	Writer
}
