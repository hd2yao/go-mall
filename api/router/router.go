package router

import (
    "github.com/gin-gonic/gin"

    "github.com/hd2yao/go-mall/common/middleware"
)

func RegisterRoutes(engine *gin.Engine) {
    // Use global middlewares
    engine.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
    routeGroup := engine.Group("")
    registerBuildingRoute(routeGroup)
}
