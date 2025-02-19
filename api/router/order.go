package router

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/controller"
	"github.com/hd2yao/go-mall/common/middleware"
)

func registerOrderRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /order/ 开头, 并且都需要身份验证
	g := rg.Group("/order/")
	g.Use(middleware.AuthUser())
	// 创建订单
	g.POST("create", controller.OrderCreate)
	// 用户订单列表
	g.GET("user-order/", controller.UserOrders)
	// 订单详情
	g.GET(":order_no/info", controller.OrderInfo)
	// 取消订单
	g.PATCH(":order_no/cancel", controller.OrderCancel)
	// 发起订单支付
	g.POST("create-pay", controller.CreateOrderPay)
}
