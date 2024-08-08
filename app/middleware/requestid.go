package middleware

import (
	"context"
	"gin-web/app/ext"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := ext.GetRequestIDByGin(c)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		//header上下文设置请求ID
		c.Writer.Header().Set(ext.RequestIDHeader, requestID)

		//设置request.context
		ctx := context.WithValue(c.Request.Context(), ext.RequestIDHeader, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
