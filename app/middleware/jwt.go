package middleware

import (
	"gin-web/core/jwt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"strings"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			//response.Fail(c, "token不存在")
			c.Abort()
			return
		}

		checkToken := strings.Split(tokenHeader, " ")
		if len(checkToken) == 0 {
			//ginc.Fail(c, "无效token")
			c.Abort()
			return
		}

		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			//ginc.Fail(c, "无效token")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(checkToken[1])
		if err != nil {
			if err == jwt.TokenExpired {
				//ginc.Fail(c, "token已过期")
				c.Abort()
				return
			}

			//ginc.Fail(c, "无效token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Request.Header.Set("user_id", cast.ToString(claims.UserId))
		c.Next()
	}
}
