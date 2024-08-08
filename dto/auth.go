package dto

import (
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
	AuthPO *model.Auth
}

func (l *LoginRespDTO) ToVO() *AuthVO {
	return AuthPOToVO(l.AuthPO)
}

type AuthVO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func AuthPOToVO(po *model.Auth) *AuthVO {
	if po == nil {
		return nil
	}
	return &AuthVO{
		ID:       po.ID,
		Username: po.Username,
	}
}
