package router

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/controller"
	"github.com/hd2yao/go-mall/common/middleware"
)

func registerReviewRoute(rg *gin.RouterGroup) {
	g := rg.Group("/review")

	// 需要登录的接口
	g.Use(middleware.AuthUser())
	{
		// 创建商品评价
		g.POST("/create", controller.CreateReview)
		// 获取用户的评价列表
		g.GET("/user/list", controller.GetUserReviews)
	}

	// 不需要登录的接口
	{
		// 获取评价详情
		g.GET("/detail", controller.GetReviewById)
		// 获取商品的评价列表
		g.GET("/commodity/list", controller.GetCommodityReviews)
		// 获取商品评价统计
		g.GET("/statistics", controller.GetReviewStatistics)
	}

	// 管理员接口
	// 需要鉴权 middleware.AuthAdmin()
	admin := g.Group("/admin")
	{
		// 商家回复评价
		admin.POST("/reply", controller.AdminReviewReply)
		// 更新评价状态
		admin.POST("/status", controller.UpdateReviewStatus)
	}
}
