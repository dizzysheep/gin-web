package errcode

import "github.com/pkg/errors"

//go:generate stringer -type ErrCode -linecomment
type ErrCode int

const (
	SuccessOk ErrCode = 1 //ok

	//ErrFail 200 -- 通用错误码
	ErrFail          ErrCode = 200000 //服务器错误
	ErrInvalidParams ErrCode = 200001 // 非法请求参数
	ErrNoAccess      ErrCode = 200002 // 无访问权限
	ErrNoFound       ErrCode = 200003 // 找不到资源
	ErrDb            ErrCode = 200004 // 数据库出错
	ErrCache         ErrCode = 200005 // 缓存出错
	ErrLimitExceeded ErrCode = 200006 // 请求异常，请稍后重试

	TokenExpired ErrCode = 210001 // token已过期
	TokenInValid ErrCode = 210002 // 无效token
	TokenEmpty   ErrCode = 210003 // token为空
)

func (e ErrCode) Code() int {
	return int(e)
}

type CustomError struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (c *CustomError) Error() string {
	return c.Code.String()
}

func NewCustomError(code ErrCode) error {
	return errors.Wrap(&CustomError{
		Code:    code,
		Message: code.String(),
	}, "")
}

func NewCustomErrorWithMessage(code ErrCode, message string) error {
	return errors.Wrap(&CustomError{
		Code:    code,
		Message: code.String() + ":" + message,
	}, message)
}
