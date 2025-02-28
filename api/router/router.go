package router

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/common/middleware"
)

func RegisterRoutes(engine *gin.Engine) {
	// 注册全局中间件
	engine.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())

	routeGroup := engine.Group("")
	registerBuildingRoute(routeGroup)
	registerUserRoutes(routeGroup)
	registerCommodityRoutes(routeGroup)
	registerCartRoutes(routeGroup)
	registerOrderRoutes(routeGroup)
}
