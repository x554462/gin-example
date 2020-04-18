package exception

import (
	"github.com/pkg/errors"
)

type ErrorWrapInterface interface {
	Error() string
	Unwrap() error
	GetCode() int
	GetMsg() string
	GetFileLine() string
}

type ErrorWrap struct {
	Msg      string
	Err      error
	Code     int
	fileLine string
}

func (e ErrorWrap) Error() string { return e.Msg }
func (e ErrorWrap) Unwrap() error { return e.Err }
func (e ErrorWrap) GetCode() int {
	if e.Code == 0 {
		if err, ok := e.Unwrap().(ErrorWrapInterface); ok {
			return err.GetCode()
		}
	}
	return e.Code
}
func (e ErrorWrap) GetMsg() string { return e.Msg }

func (e ErrorWrap) GetFileLine() string {
	if e.fileLine == "" {
		if err, ok := e.Unwrap().(ErrorWrapInterface); ok {
			return err.GetFileLine()
		}
	}
	return e.fileLine
}

var (
	RootError          = ErrorWrap{Msg: "服务器内部错误", Err: errors.New("服务器内部错误"), Code: 500}
	LoadConfigError    = ErrorLoadConfig{ErrorWrap{Msg: "配置错误", Err: RootError, Code: 400}}
	ResponseError      = ErrorResponse{ErrorWrap{Msg: "服务器开小差", Err: RootError, Code: 400}}
	UnauthorizedError  = ErrorResponse{ErrorWrap{Msg: "权限不足", Err: RootError, Code: 401}}
	ValidateError      = ErrorValidate{ErrorWrap{Msg: "参数有误", Err: RootError, Code: 400}}
	RuntimeError       = ErrorRuntime{ErrorWrap{Msg: "服务器内部错误", Err: RootError, Code: 500}}
	JsonRuntimeError   = ErrorJsonRuntime{ErrorWrap{Msg: "Json decode error", Err: RuntimeError, Code: 500}}
	ModelRuntimeError  = ErrorJsonRuntime{ErrorWrap{Msg: "Model runtime error", Err: RuntimeError, Code: 400}}
	ModelNotFoundError = ErrorModelNotFound{ErrorWrap{Msg: "数据未找到", Err: ModelRuntimeError, Code: 400}}
)

// 配置错误
type ErrorLoadConfig struct {
	ErrorWrap
}

// 响应错误
type ErrorResponse struct {
	ErrorWrap
}

// 权限不足权限
type ErrorUnauthorized struct {
	ErrorWrap
}

// 运行错误
type ErrorRuntime struct {
	ErrorWrap
}

// 运行错误
type ErrorJsonRuntime struct {
	ErrorWrap
}

// 数据未找到错误
type ErrorModelNotFound struct {
	ErrorWrap
}

// 数据校验错误
type ErrorValidate struct {
	ErrorWrap
}
