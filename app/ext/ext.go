package ext

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "x-request-id"
)

func GetRequestIDByGin(c *gin.Context) string {
	return c.GetHeader(RequestIDHeader)
}

func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDHeader).(string)
	return id
}
