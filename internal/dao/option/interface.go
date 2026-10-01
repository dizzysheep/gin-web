package option

import (
	"context"
	"gin-web/internal/model"
)

type OptionDao interface {
	SelectAll(ctx context.Context) ([]*model.Option, error)
	SelectOneByKey(ctx context.Context, key string) (*model.Option, error)
	Upsert(ctx context.Context, option *model.Option) error
}
