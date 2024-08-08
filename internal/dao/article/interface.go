package article

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
)

type Reader interface {
	Count(ctx context.Context, conditions common.GormConditions) (int64, error)
	SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Article, error)
	SelectOne(ctx context.Context, id int64) (*model.Article, error)
}

type Writer interface {
	InsertOne(ctx context.Context, model *model.Article) error
	UpdateOne(ctx context.Context, model *model.Article) error
}

type ArticleDao interface {
	Reader
	Writer
}
