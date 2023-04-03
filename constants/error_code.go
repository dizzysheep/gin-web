package constants

type ErrorCode int

const (
	SUCCESS     ErrorCode = 1 // 成功
	ServerError ErrorCode = 0 // 失败

	ParamsValidatorError ErrorCode = 10001
	PageEmpty            ErrorCode = 21001
	LimitEmpty           ErrorCode = 21002

	TokenNotExist  ErrorCode = 21003
	TokenValidFail ErrorCode = 21004
	TokenExpired   ErrorCode = 21005
)

var StatusCode = map[ErrorCode]string{
	ServerError:          "systemError",
	SUCCESS:              "success",
	ParamsValidatorError: "paramError",
	PageEmpty:            "分页不能为空",
	LimitEmpty:           "分页数量不能为空",
	TokenNotExist:        "token不能为空",
	TokenValidFail:       "token不合法",
	TokenExpired:         "token授权已过期，请重新登录",
}

//转化成对于中文错误信息
func (code ErrorCode) Transform() string {
	return StatusCode[code]
}
