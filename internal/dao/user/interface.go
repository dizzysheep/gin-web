package user

import (
	"context"
	"gin-web/internal/model"
)

type UserDao interface {
	SelectOne(ctx context.Context, id int64) (*model.User, error)
	SelectOneByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateOne(ctx context.Context, user *model.User) error
}
