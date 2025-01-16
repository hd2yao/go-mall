package router

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/controller"
	"github.com/hd2yao/go-mall/common/middleware"
)

// 存放 User 模块的路由
func registerUserRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /user 开头
	g := rg.Group("/user")
	// 刷新Token
	g.GET("token/refresh", controller.RefreshUserToken)
	// 注册用户
	g.POST("register", controller.RegisterUser)
	// 登录
	g.POST("login", controller.LoginUser)
	// 登出用户
	g.DELETE("logout", middleware.AuthUser(), controller.LogoutUser)
	// 申请重置密码
	g.POST("password/apply-reset", controller.PasswordResetApply)
}
