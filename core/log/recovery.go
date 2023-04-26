package log

import (
	"gin-web/core/result"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(c *gin.Context, err interface{}) {
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				brokenPipe = true
			}
		}
	}

	httpRequest, _ := httputil.DumpRequest(c.Request, false)
	headers := strings.Split(string(httpRequest), "\r\n")
	for idx, header := range headers {
		current := strings.Split(header, ":")
		if current[0] == "Authorization" {
			headers[idx] = current[0] + ": *"
		}
	}

	//panic转Error级别日志
	Get(c).WithField("is_panic", true).Error(err)

	// If the connection is dead, we can't write a status to it.
	if brokenPipe {
		c.Error(err.(error)) // nolint: errcheck
		c.Abort()
	} else {
		c.JSON(http.StatusInternalServerError, result.Fail("抱歉服务错误，请稍后重试"))
		c.Abort()
	}
}
