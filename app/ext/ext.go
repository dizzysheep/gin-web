package ext

import (
	"context"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "x-request-id"
	UserInfoKey     = "user-info"
)

func GetRequestIDByGin(c *gin.Context) string {
	return c.GetHeader(RequestIDHeader)
}

func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDHeader).(string)
	return id
}

func GetUsername(c *gin.Context) string {
	value, ok := c.Get(UserInfoKey)
	if !ok {
		return ""
	}
	userInfo, ok := value.(*model.Auth)
	if !ok {
		return ""
	}
	return userInfo.Username
}
