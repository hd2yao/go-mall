package router

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/controller"
)

func registerCommodityRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /commodity/ 开头
	g := rg.Group("/commodity/")
	// 按层级划分的所有商品分类
	g.GET("category-hierarchy/", controller.GetCategoryHierarchy)
}
