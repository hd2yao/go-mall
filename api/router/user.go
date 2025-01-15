package router

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/controller"
)

// 存放 User 模块的路由
func registerUserRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /user 开头
	g := rg.Group("/user")
	// 注册用户
	g.POST("register", controller.RegisterUser)
}
