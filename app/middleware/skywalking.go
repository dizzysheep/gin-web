package middleware

import (
	"gin-web/core/skywalking"
	"github.com/gin-gonic/gin"
)

func SkywalkingMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if skywalking.Reporter == nil || skywalking.Tracer == nil {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		span, ctx := skywalking.CreateEntrySpan(path, c)
		skywalking.SetContext(c, ctx)

		c.Next()

		if span != nil {
			span.End()
		}
	}
}
