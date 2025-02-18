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
}
