package errcode

import "net/http"

var codes = map[int]struct{}{}

// 此处为公共的错误码, 预留 10000000 ~ 10000099 间的 100 个错误码
var (
	Success            = newError(0, "success")
	ErrServer          = newError(10000000, "服务器内部错误")
	ErrParams          = newError(10000001, "参数错误, 请检查")
	ErrNotFound        = newError(10000002, "资源未找到")
	ErrPanic           = newError(10000003, "(*^__^*)系统开小差了,请稍后重试") // 无预期的panic错误
	ErrToken           = newError(10000004, "Token无效")
	ErrForbidden       = newError(10000005, "未授权") // 访问一些未授权的资源时的错误
	ErrTooManyRequests = newError(10000006, "请求过多")
	ErrCoverData       = newError(10000007, "ConvertDataError") // 数据转换错误
)

// 各个业务模块自定义的错误码, 从 10000100 开始, 可以按照不同的业务模块划分不同的号段

// 用户模块相关错误码 10000100 ~ 10000199
var (
	ErrUserInvalid        = newError(10000101, "用户异常")
	ErrUserNameOccupied   = newError(10000102, "用户名已被占用")
	ErrUserNotRight       = newError(10000103, "用户名或密码不正确")
	ErrPasswordComplexity = newError(10000104, "密码复杂度不满足要求")
)

// 商品模块相关错误码 10000200 ~ 1000299
var (
	ErrCommodityNotExists = newError(10000200, "商品不存在")
	ErrCommodityStockOut  = newError(10000201, "库存不足")
)

// 购物车模块相关错误码 10000300 ～ 1000399
var (
	ErrCartItemParam = newError(10000300, "购物项参数异常")
	ErrCartWrongUser = newError(10000301, "用户购物信息不匹配")
)

// 订单模块相关错误码 10000500 ~ 10000599
var (
	ErrOrderParams              = newError(10000500, "订单参数异常")
	ErrOrderCanNotBeChanged     = newError(10000501, "订单不可修改")
	ErrOrderUnsupportedPayScene = newError(10000502, "支付场景暂不支持")
)

// 评价模块相关错误码 10000600 ~ 10000699
var (
	ErrReviewParams              = newError(10000600, "评价参数异常")
	ErrReviewStatusCanNotChanged = newError(10000601, "评价状态不可修改")
	ErrReviewUnsupportedScene    = newError(10000602, "评价场景暂不支持")
)

// HttpStatusCode 返回 HTTP 状态码
func (e *AppError) HttpStatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ErrServer.Code(), ErrPanic.Code():
		return http.StatusInternalServerError
	case ErrParams.Code(), ErrUserInvalid.Code(), ErrUserNameOccupied.Code(), ErrUserNotRight.Code(), ErrPasswordComplexity.Code(),
		ErrCommodityNotExists.Code(), ErrCommodityStockOut.Code(), ErrCartItemParam.Code(), ErrOrderParams.Code(),
		ErrReviewParams.Code(), ErrReviewUnsupportedScene.Code():
		return http.StatusBadRequest
	case ErrNotFound.Code():
		return http.StatusNotFound
	case ErrTooManyRequests.Code():
		return http.StatusTooManyRequests
	case ErrToken.Code():
		return http.StatusUnauthorized
	case ErrForbidden.Code(), ErrCartWrongUser.Code(), ErrOrderCanNotBeChanged.Code(), ErrReviewStatusCanNotChanged.Code():
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
