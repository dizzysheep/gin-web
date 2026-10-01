package dto

import (
	"gin-web/app/ext"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

func LoginReqToDTO(c *gin.Context) (*LoginReqDTO, error) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &LoginReqDTO{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

type LoginReqDTO struct {
	Username string
	Password string
}

type LoginRespDTO struct {
	Jwt string
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (l *LoginRespDTO) ToVO() *LoginResponse {
	return &LoginResponse{l.Jwt}
}

// ChangePasswordRequest 修改密码
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=32"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

func ChangePasswordReqToDTO(c *gin.Context) (*ChangePasswordReqDTO, error) {
	var req ChangePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &ChangePasswordReqDTO{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
		User:        ext.GetUser(c),
	}, nil
}

type ChangePasswordReqDTO struct {
	OldPassword string
	NewPassword string
	User        *model.User
}

// UserInfoRespDTO 当前用户信息
type UserInfoRespDTO struct {
	User *model.User
}

type UserInfoResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	Role      int8   `json:"role"`
	LastLogin uint32 `json:"last_login"`
}

func (u *UserInfoRespDTO) ToVO() *UserInfoResponse {
	if u.User == nil {
		return &UserInfoResponse{}
	}
	return &UserInfoResponse{
		ID:        u.User.ID,
		Username:  u.User.Username,
		Nickname:  u.User.Nickname,
		Avatar:    u.User.Avatar,
		Email:     u.User.Email,
		Role:      u.User.Role,
		LastLogin: u.User.LastLogin,
	}
}
