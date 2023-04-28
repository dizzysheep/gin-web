// 操作请求 gin.Context 信息
package ginc

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// TraceIDKey 请求唯一标识，类型：string
	TraceIDKey = "x-trace-id"
	// StartTimeKey 请求开始时间，类型：time.Time
	StartTimeKey = "x-start-time"
)

// GetTraceID 获取用户请求标识
func GetTraceID(ctx *gin.Context) string {
	id, _ := ctx.Value(TraceIDKey).(string)
	return id
}

// WithTraceID 注入 trace_id
func WithTraceID(ctx *gin.Context) string {
	traceID := NewTraceID()

	ctx.Writer.Header().Set(TraceIDKey, traceID)
	ctx.Set(TraceIDKey, traceID)
	return traceID
}

// WithStartTime 注入 trace_id
func WithStartTime(ctx *gin.Context, startTime time.Time) {
	ctx.Set(StartTimeKey, startTime)
	return
}

// NewTraceID 生成
func NewTraceID() string {
	return uuid.New().String()
}
