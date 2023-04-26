package log

import (
	"bytes"
	"fmt"
	"gin-web/core/config"
	"gin-web/core/ginc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const StartMicroTime = "StartMircoTime"

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	TraceID string
	// Request Body
	RequestBody string
	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// ServerIP equals Context's ServerIP method.
	ServerIP string
	// Method is the HTTP method given to the request.
	Method string
	// UserAgent
	UserAgent string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// BodySize is the size of the Response Body
	BodySize int
	// Response Body
	Body string
	// isError
	isError bool
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
}

// 处理返回数据
func (param *LogFormatterParams) getBody(c *gin.Context, b []byte) {
	if param.StatusCode < http.StatusBadRequest {
		return
	}

	param.RequestBody = getRequestBody(c)
	if len(b) == 0 { // 没有数据返回
		return
	}

	param.Body = string(b)
	return
}

func getRequestBody(c *gin.Context) string {
	method := c.Request.Method
	if method == "GET" {
		return ""
	}

	if data, err := c.GetRawData(); err != nil {
		return string(data)
	}

	return ""
}

// Logger is the logrus logger handler
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(StartMicroTime, time.Now().UnixNano()/1000000)

		defer func() {
			if err := recover(); err != nil {
				RecoveryWithWriter(c, err)
			}
		}()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		start := time.Now()
		if raw != "" {
			path = path + "?" + raw
		}

		bWriter := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bWriter

		//Trace
		traceID := ginc.WithTraceID(c)
		ginc.WithStartTime(c, start)

		c.Next()

		param := &LogFormatterParams{
			TraceID:      traceID,
			Keys:         c.Keys,
			ServerIP:     ServerIP,
			TimeStamp:    time.Now(),
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			Path:         path,
			StatusCode:   c.Writer.Status(),
			UserAgent:    c.Request.UserAgent(),
			BodySize:     c.Writer.Size(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}
		param.Latency = param.TimeStamp.Sub(start)
		param.getBody(c, bWriter.body.Bytes())

		param.LogFormatter()
	}
}

func (param *LogFormatterParams) LogFormatter() {
	entry := logrus.WithFields(logrus.Fields{
		"trace_id":     param.TraceID,
		"app_name":     config.AppName,
		"hostname":     config.Hostname,
		"status":       param.StatusCode,
		"latency":      fmt.Sprintf("%13v", param.Latency),
		"client_ip":    param.ClientIP,
		"server_ip":    ServerIP,
		"timestamp":    param.TimeStamp,
		"method":       param.Method,
		"path":         param.Path,
		"size":         param.BodySize,
		"user_agent":   param.UserAgent,
		"request_body": param.RequestBody,
		"env":          config.Env,
	})

	if len(param.ErrorMessage) > 0 {
		entry.Error(param.ErrorMessage)
		return
	}

	if param.StatusCode >= http.StatusOK && param.StatusCode < http.StatusBadRequest {
		entry.Info(param.Body)
	} else if param.StatusCode >= http.StatusBadRequest && param.StatusCode < http.StatusInternalServerError {
		entry.Warn(param.Body)
	} else {
		entry.Error(param.Body)
	}
}
