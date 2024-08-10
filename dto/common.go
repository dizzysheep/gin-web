package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

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
