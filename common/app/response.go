package app

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
)

type response struct {
	ctx        *gin.Context
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	RequestId  string      `json:"request_id"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"Pagination,omitempty"`
}

func NewResponse(ctx *gin.Context) *response {
	return &response{ctx: ctx}
}

// SetPagination 设置 response 分页信息
func (r *response) SetPagination(pagination *Pagination) *response {
	r.Pagination = pagination
	return r
}

// Success 带数据的成功响应
func (r *response) Success(data interface{}) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	requestId := ""

	// 获取请求上下文中的 traceid，作为响应中的 requestId
	// traceid 来自于项目全局中间件 StartTrace
	if val, exists := r.ctx.Get("traceId"); exists {
		requestId = val.(string)
	}
	r.RequestId = requestId
	r.Data = data

	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}

// SuccessOk 不带数据的成功响应
// 针对只需要知道成功状态的接口响应，简化接口程序的调用
func (r *response) SuccessOk() {
	r.Success("")
}

// Error 带错误信息的响应
func (r *response) Error(err *errcode.AppError) {
	r.Code = err.Code()
	r.Msg = err.Msg()
	requestId := ""
	if val, exists := r.ctx.Get("traceId"); exists {
		requestId = val.(string)
	}
	r.RequestId = requestId
	// 兜底记一条响应错误，项目自定义的 AppError 中有错误链条，方便出错后排查问题
	logger.New(r.ctx).Error("api_response_error", "err", err)
	r.ctx.JSON(err.HttpStatusCode(), r)
}
