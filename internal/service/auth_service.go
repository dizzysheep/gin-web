package service

import (
	"errors"
	"gin-web/internal/dao"
	"gin-web/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type AuthService struct {
	dao *dao.AuthDao
}

func NewAuthService(ctx *gin.Context) *AuthService {
	return &AuthService{dao.NewAuthDao(GetBlogDB(ctx))}
}

func (svc *AuthService) CheckAuth(req dto.AuthReqDTO) (*dto.AuthDTO, error) {
	user, err := svc.dao.GetInfoByUserName(req.Username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Password != req.Password {
		return nil, errors.New("邮箱或者密码不正确")
	}

	var authDTO dto.AuthDTO
	if err := mapstructure.Decode(user, &authDTO); err != nil {
		return nil, errors.New("转化authDTO失败")
	}

	return &authDTO, nil
}
