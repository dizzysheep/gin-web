package ext

import (
	"context"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
	"strings"
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

// GetUser 当前登录用户（JWT中间件注入）
func GetUser(c *gin.Context) *model.User {
	value, ok := c.Get(UserInfoKey)
	if !ok {
		return nil
	}
	userInfo, ok := value.(*model.User)
	if !ok {
		return nil
	}
	return userInfo
}

func GetUsername(c *gin.Context) string {
	user := GetUser(c)
	if user == nil {
		return ""
	}
	return user.Username
}

// ExtractToken 从请求头提取Bearer token（登出使用）
func ExtractToken(c *gin.Context) string {
	tokenHeader := c.Request.Header.Get("Authorization")
	checkToken := strings.Split(tokenHeader, " ")
	if len(checkToken) != 2 || checkToken[0] != "Bearer" {
		return ""
	}
	return checkToken[1]
}
