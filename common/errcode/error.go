package errcode

import (
    "encoding/json"
    "fmt"
)

type AppError struct {
    code  int    `json:"code"`
    msg   string `json:"msg"`
    cause error  `json:"cause"`
}

// 实现 error 接口
func (e *AppError) Error() string {
    if e == nil {
        return ""
    }
    formattedErr := struct {
        Code  int    `json:"code"`
        Msg   string `json:"msg"`
        Cause string `json:"cause"`
    }{
        Code: e.Code(),
        Msg:  e.Msg(),
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
    if _, duplicated := codes[code]; duplicated {
        panic(fmt.Sprintf("错误码 %d 不能重复, 请检查后更换", code))
    }
    return &AppError{code: code, msg: msg}
}
