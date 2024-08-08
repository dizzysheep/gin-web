package auth

import (
	"context"
	"gin-web/dto"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginReqDTO) (*dto.LoginRespDTO, error)
}
