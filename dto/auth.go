package dto

import (
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
	Jwt string
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (l *LoginRespDTO) ToVO() *LoginResponse {
	return &LoginResponse{l.Jwt}
}
