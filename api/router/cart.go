package router

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/controller"
	"github.com/hd2yao/go-mall/common/middleware"
)

// 存放购物车模块的路由

func registerCartRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /cart/ 开头, 并且都需要身份验证
	g := rg.Group("/cart/")
	g.Use(middleware.AuthUser())
	// 添加到购物车
	g.POST("add-item", controller.AddCartItem)
	// 查看购物项账单 -- 确认下单前用来显示商品和支付金额明细
	g.GET("/item/check-bill", controller.CheckCartItemBill)
}
