package errcode

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

type AppError struct {
	code     int    `json:"code"`
	msg      string `json:"msg"`
	cause    error  `json:"cause"`    // 错误链
	occurred string `json:"occurred"` // 错误发生的位置（函数、文件、行号）
}

// 实现 error 接口
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	formattedErr := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Cause    string `json:"cause"`
		Occurred string `json:"occurred"`
	}{
		Code:     e.Code(),
		Msg:      e.Msg(),
		Occurred: e.occurred,
	}
	if e.cause != nil {
		formattedErr.Cause = e.cause.Error()
	}
	errByte, _ := json.Marshal(formattedErr)
	return string(errByte)
}

func (e *AppError) String() string {
	return e.Error()
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Msg() string {
	return e.msg
}

func newError(code int, msg string) *AppError {
	if code > -1 {
		if _, duplicated := codes[code]; duplicated {
			panic(fmt.Sprintf("预定义错误码 %d 不能重复, 请检查后更换", code))
		}
		codes[code] = struct{}{}
	}

	return &AppError{code: code, msg: msg}
}

// WithCause 在逻辑执行中出现错误, 比如 dao 层返回的数据库查询错误
// 可以在领域层返回预定义的错误前附加上导致错误的基础错误
// 如果业务模块预定义的错误码比较详细, 可以使用这个方法, 反之错误码定义的比较笼统建议使用Wrap方法包装底层错误生成项目自定义 Error
// 并将其记录到日志后再使用预定义错误码返回接口响应
func (e *AppError) WithCause(err error) *AppError {
	newErr := e.Clone()
	newErr.cause = err
	newErr.occurred = getAppErrOccurredInfo()
	return newErr
}

func (e *AppError) UnWrap() error {
	return e.cause
}

// Is 与上面的 UnWrap 一起让 *AppError 支持 errors.Is(err, target)
func (e *AppError) Is(target error) bool {
	targetErr, ok := target.(*AppError)
	if !ok {
		return false
	}
	return targetErr.Code() == e.Code()
}

func (e *AppError) Clone() *AppError {
	return &AppError{
		code:     e.code,
		msg:      e.msg,
		cause:    e.cause,
		occurred: e.occurred,
	}
}

// Wrap 用于逻辑中包装底层函数返回的 error 和 WithCause 一样都是为了记录错误链条
// 该方法生成的 error 用于日志记录, 返回响应请使用预定义好的error
func Wrap(msg string, err error) *AppError {
	if err == nil {
		return nil
	}
	appErr := &AppError{code: -1, msg: msg, cause: err}
	appErr.occurred = getAppErrOccurredInfo()
	return appErr
}

// getAppErrOccurredInfo 获取项目中调用Wrap或者WithCause方法时的程序位置, 方便排查问题
func getAppErrOccurredInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	file = path.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("func: %s, file: %s, line: %d", funcName, file, line)
}
