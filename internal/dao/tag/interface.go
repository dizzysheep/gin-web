package tag

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
)

type Reader interface {
	Count(ctx context.Context, conditions common.GormConditions) (int64, error)
	SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Tag, error)
	SelectOne(ctx context.Context, id int64) (*model.Tag, error)
	SelectManyByIDs(ctx context.Context, ids []int64) ([]*model.Tag, error)
	CountArticleByIDs(ctx context.Context, ids []int64) (map[int64]int64, error)
	SelectRelationsByArticleIDs(ctx context.Context, articleIDs []int64) ([]*model.ArticleTag, error)
	CountRelationsByTagID(ctx context.Context, tagID int64) (int64, error)
	ExistName(ctx context.Context, name string, excludeID int64) (bool, error)
}

type Writer interface {
	InsertOne(ctx context.Context, tag *model.Tag) error
	UpdateOne(ctx context.Context, tag *model.Tag) error
	DeleteOne(ctx context.Context, id int64) error
}

type TagDao interface {
	Reader
	Writer
}
