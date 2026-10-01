package auth

import (
	"context"

	"gin-web/core/jwt"
	"gin-web/core/xtime"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/errcode"
	"gin-web/internal/model"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	daos *dao.Daos
}

func NewAuthService(daos *dao.Daos) AuthService {
	return &authService{daos}
}

func (s *authService) Login(ctx context.Context, req *dto.LoginReqDTO) (*dto.LoginRespDTO, error) {
	user, err := s.daos.User.SelectOneByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.NewCustomError(errcode.ErrPasswordWrong)
		}
		return nil, errors.Wrap(err, "AuthService.Login")
	}

	// 密码使用 bcrypt 校验，数据库中存储的是哈希值
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errcode.NewCustomError(errcode.ErrPasswordWrong)
	}

	// 记录最后登录时间
	user.LastLogin = uint32(xtime.GetTimestamp())
	if err := s.daos.User.UpdateOne(ctx, user); err != nil {
		return nil, errors.Wrap(err, "update last login fail")
	}

	token, err := jwt.GenerateToken(user, jwt.CurrentVer(user.ID))
	if err != nil {
		return nil, err
	}

	return &dto.LoginRespDTO{
		Jwt: token,
	}, nil
}

// Logout 登出：token 加入Redis黑名单
func (s *authService) Logout(ctx context.Context, token string) error {
	return jwt.Logout(token)
}

// Info 当前用户信息
func (s *authService) Info(ctx context.Context, user *model.User) (*dto.UserInfoRespDTO, error) {
	return &dto.UserInfoRespDTO{User: user}, nil
}

// ChangePassword 修改密码：校验旧密码后更新，并递增版本号使所有旧token失效
func (s *authService) ChangePassword(ctx context.Context, req *dto.ChangePasswordReqDTO) error {
	user := req.User
	if user == nil {
		return errcode.NewCustomError(errcode.TokenInValid)
	}

	dbUser, err := s.daos.User.SelectOne(ctx, user.ID)
	if err != nil {
		return errors.Wrap(err, "select user fail")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.OldPassword)); err != nil {
		return errcode.NewCustomError(errcode.ErrPasswordWrong)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "hash password fail")
	}

	dbUser.Password = string(hash)
	if err := s.daos.User.UpdateOne(ctx, dbUser); err != nil {
		return errors.Wrap(err, "update password fail")
	}

	if err := jwt.BumpVer(dbUser.ID); err != nil {
		return errors.Wrap(err, "bump token version fail")
	}

	return nil
}
