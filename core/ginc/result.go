package ginc

import (
	"gin-web/core/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ok(c *gin.Context, data interface{}, errs ...interface{}) {
	res := result.Ok(data, errs...).SetRequestId(c.GetString("X-Request-ID"))
	c.JSON(http.StatusOK, res)
}

func Fail(c *gin.Context, errs ...interface{}) {
	res := result.Fail(errs...).SetRequestId(c.GetString("X-Request-ID"))
	c.JSON(http.StatusOK, res)
}

func Page(c *gin.Context, page *result.PageInfo, errs ...interface{}) {
	res := result.Page(page, errs...).SetRequestId(c.GetString("X-Request-ID"))
	c.JSON(http.StatusOK, res)
}
