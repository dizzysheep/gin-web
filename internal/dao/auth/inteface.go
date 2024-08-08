package auth

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
)

type Reader interface {
	Count(ctx context.Context, conditions common.GormConditions) (int64, error)
	SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Auth, error)
	SelectOne(ctx context.Context, id int64) (*model.Auth, error)
	SelectOneByWhere(ctx context.Context, conditions common.GormConditions) (*model.Auth, error)
}

type Writer interface {
	InsertOne(ctx context.Context, model *model.Auth) error
	UpdateOne(ctx context.Context, model *model.Auth) error
}

type AuthDao interface {
	Reader
	Writer
}
