package comment

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
)

type Reader interface {
	Count(ctx context.Context, conditions common.GormConditions) (int64, error)
	SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Comment, error)
	SelectOne(ctx context.Context, id int64) (*model.Comment, error)
}

type Writer interface {
	InsertOne(ctx context.Context, comment *model.Comment) error
	UpdateColumns(ctx context.Context, id int64, cols map[string]interface{}) error
	DeleteOne(ctx context.Context, id int64) error
}

type CommentDao interface {
	Reader
	Writer
}
