package middleware

import (
	"gin-web/core/ginc"
	"gin-web/core/jwt"
	"gin-web/pkg/constants"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			ginc.Fail(c, constants.TokenNotExist)
			c.Abort()
			return
		}

		checkToken := strings.Split(tokenHeader, " ")
		if len(checkToken) == 0 {
			ginc.Fail(c, constants.TokenValidFail)
			c.Abort()
			return
		}

		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			ginc.Fail(c, constants.TokenValidFail)
			c.Abort()
			return
		}

		claims, err := util.ParseToken(checkToken[1])
		if err != nil {
			if err == util.TokenExpired {
				ginc.Fail(c, constants.TokenExpired)
				c.Abort()
				return
			}

			ginc.Fail(c, constants.TokenValidFail)
			c.Abort()
			return
		}

		c.Set("username", claims)
		c.Next()
	}
}
