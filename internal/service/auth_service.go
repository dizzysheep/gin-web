package service

import (
	"errors"
	"gin-web/internal/dao"
	"gin-web/internal/dto"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (svc *AuthService) CheckAuth(req dto.AuthReqDTO) error {
	user, err := dao.NewAuthDao().GetInfoByUserName(req.Username)
	if err != nil {
		return errors.New("用户不存在")
	}

	if user.Password != req.Password {
		return errors.New("邮箱或者密码不正确")
	}

	return nil
}
