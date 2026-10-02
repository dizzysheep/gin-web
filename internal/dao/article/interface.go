package article

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
)

type Reader interface {
	Count(ctx context.Context, conditions common.GormConditions) (int64, error)
	SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Article, error)
	CountByTagID(ctx context.Context, conditions common.GormConditions, tagID int64) (int64, error)
	SelectManyByTagID(ctx context.Context, conditions common.GormConditions, pager *common.Pagination, tagID int64) ([]*model.Article, error)
	SelectOne(ctx context.Context, id int64) (*model.Article, error)
	SelectOneBySlug(ctx context.Context, slug string) (*model.Article, error)
	ExistSlug(ctx context.Context, slug string, excludeID int64) (bool, error)
	Archive(ctx context.Context) ([]*model.ArchiveRow, error)
	CountPublishedByDay(ctx context.Context, startUnix, endUnix int64) ([]*model.ArticleTrendRow, error)
}

type Writer interface {
	Create(ctx context.Context, article *model.Article, tagIDs []int64) error
	Update(ctx context.Context, article *model.Article, tagIDs []int64, withTags bool) error
	UpdateColumns(ctx context.Context, id int64, cols map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
	IncrViewCount(ctx context.Context, id int64, delta int64) error
}

type ArticleDao interface {
	Reader
	Writer
}
