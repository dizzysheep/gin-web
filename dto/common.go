package dto

import (
	"errors"
	"gin-web/app/ext"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// StateRequest 启用/禁用请求体
type StateRequest struct {
	State *int8 `json:"state" binding:"required,oneof=0 1" example:"1"`
}

// StateReqDTO 启用/禁用请求参数
type StateReqDTO struct {
	ID       int64
	State    int8
	Username string
}

// StateReqToDTO 解析启用/禁用请求：路径参数 id + 请求体 state + 当前登录用户
func StateReqToDTO(c *gin.Context) (*StateReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}

	var req StateRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	return &StateReqDTO{
		ID:       id,
		State:    *req.State,
		Username: ext.GetUsername(c),
	}, nil
}

type IDReqDTO struct {
	ID int64
}

func IDReqDTOFromRequest(c *gin.Context) (*IDReqDTO, error) {
	id, err := GetIDByCtx(c)
	if err != nil {
		return nil, err
	}
	return &IDReqDTO{ID: id}, nil
}

func GetIDByCtx(c *gin.Context) (int64, error) {
	id, err := cast.ToInt64E(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("不合法的 ID")
	}

	return id, nil
}
