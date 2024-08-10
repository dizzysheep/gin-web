package middleware

import (
	"context"
	"errors"
	"gin-web/app/ext"
	"gin-web/app/response"
	"gin-web/core/jwt"
	"gin-web/internal/errcode"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			response.FailErr(c, err)
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			if errors.Is(err, jwt.TokenExpired) {
				response.Fail(c, errcode.TokenExpired)
				c.Abort()
				return
			}

			response.Fail(c, errcode.TokenInValid)
			c.Abort()
			return
		}

		c.Set(ext.UserInfoKey, claims.UserInfo)
		c.Request.WithContext(context.WithValue(c.Request.Context(), ext.UserInfoKey, claims.UserInfo))
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	tokenHeader := c.Request.Header.Get("Authorization")
	if tokenHeader == "" {
		return "", errcode.NewCustomError(errcode.TokenEmpty)
	}

	checkToken := strings.Split(tokenHeader, " ")
	if len(checkToken) == 0 {
		return "", errcode.NewCustomError(errcode.TokenInValid)
	}

	if len(checkToken) != 2 || checkToken[0] != "Bearer" {
		return "", errcode.NewCustomError(errcode.TokenInValid)
	}

	return checkToken[1], nil
}
