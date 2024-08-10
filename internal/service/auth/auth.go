package auth

import (
	"context"
	"gin-web/core/jwt"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"github.com/pkg/errors"
)

type authService struct {
	daos *dao.Daos
}

func NewAuthService(daos *dao.Daos) AuthService {
	return &authService{daos}
}

func (s *authService) Login(ctx context.Context, req *dto.LoginReqDTO) (*dto.LoginRespDTO, error) {
	userCond := common.GormConditions{
		&common.EqCond{
			Field: "username",
			Value: req.Username,
		},
	}
	user, err := s.daos.Auth.SelectOneByWhere(ctx, userCond)
	if err != nil {
		return nil, errors.Wrap(err, "AuthService.Login")
	}

	if user.Password != req.Password {
		return nil, errors.New("邮箱或者密码不正确")
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginRespDTO{
		Jwt: token,
	}, nil
}
