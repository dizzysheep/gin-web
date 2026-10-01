package auth

import (
	"context"
	"gin-web/dto"
	"gin-web/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginReqDTO) (*dto.LoginRespDTO, error)
	Logout(ctx context.Context, token string) error
	Info(ctx context.Context, user *model.User) (*dto.UserInfoRespDTO, error)
	ChangePassword(ctx context.Context, req *dto.ChangePasswordReqDTO) error
}
