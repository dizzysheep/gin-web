package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin-web/internal/errcode"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type FieldError struct {
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
	Field   string `json:"field,omitempty"`
	Value   any    `json:"value,omitempty"`
}

func BuildResponse(code int, msg string, data interface{}) *Response {
	return &Response{Code: code, Data: data, Msg: msg}
}

func Ok(c *gin.Context, data interface{}) {
	res := BuildResponse(errcode.SuccessOk.Code(), errcode.SuccessOk.String(), data)
	c.JSON(http.StatusOK, res)
}

func Fail(c *gin.Context, code errcode.ErrCode) {
	res := BuildResponse(code.Code(), code.String(), nil)
	c.JSON(http.StatusOK, res)
}

func FailErr(c *gin.Context, err error) {
	code := errcode.ErrFail.Code()
	msg := err.Error()

	var customErr *errcode.CustomError
	if errors.As(err, &customErr) {
		code = customErr.Code.Code()
		if customErr.Message != "" {
			msg = customErr.Message
		} else {
			msg = customErr.Error()
		}
	}

	res := BuildResponse(code, msg, nil)
	c.JSON(http.StatusOK, res)
}

func InternalServerError(c *gin.Context) {
	res := BuildResponse(errcode.ErrFail.Code(), errcode.ErrFail.String(), nil)
	c.JSON(http.StatusOK, res)
}

func BadRequest(c *gin.Context, errs ...error) {
	if len(errs) == 0 {
		Fail(c, errcode.ErrInvalidParams)
		return
	}

	var fieldErrs []FieldError
	for _, err := range errs {
		fieldErrs = append(fieldErrs, handleError(err)...)
	}

	res := BuildResponse(errcode.ErrInvalidParams.Code(), errcode.ErrInvalidParams.String(), nil)
	if len(fieldErrs) > 0 {
		res.Msg = fieldErrs[0].Message
	}
	c.JSON(http.StatusBadRequest, res)
}

func handleError(err error) []FieldError {
	switch e := err.(type) {
	case validator.ValidationErrors:
		return translateErrors(e)
	case *json.SyntaxError:
		return []FieldError{{Message: "无效的请求参数"}}
	case *json.UnmarshalTypeError:
		return []FieldError{{
			Message: fmt.Sprintf("`%s`的类型必须是 `%s`", e.Field, e.Type.String()),
		}}
	default:
		return []FieldError{{Message: e.Error()}}
	}
}
