package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/logic/appservice"
)

// CreateReview 创建商品评价
func CreateReview(c *gin.Context) {
	userId := app.GetUserId(c)
	var req request.ReviewCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	err := appservice.NewReviewAppSvc(c).CreateReview(&req, userId)
	if err != nil {
		app.NewResponse(c).Error(err)
		return
	}

	app.NewResponse(c).Success(nil)
}

// GetReviewById 获取评价详情
func GetReviewById(c *gin.Context) {
	var req struct {
		ReviewId uint `form:"review_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	review, err := appservice.NewReviewAppSvc(c).GetReviewById(req.ReviewId)
	if err != nil {
		app.NewResponse(c).Error(err)
		return
	}

	app.NewResponse(c).Success(review)
}

// GetUserReviews 获取用户的评价列表
func GetUserReviews(c *gin.Context) {
	userId := app.GetUserId(c)
	pagination := app.GetPagination(c)

	reviews, err := appservice.NewReviewAppSvc(c).GetUserReviews(userId, pagination)
	if err != nil {
		app.NewResponse(c).Error(err)
		return
	}

	app.NewResponse(c).Success(reviews)
}

// GetCommodityReviews 获取商品的评价列表
func GetCommodityReviews(c *gin.Context) {
	var req struct {
		CommodityId int64 `form:"commodity_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	pagination := app.GetPagination(c)
	reviews, err := appservice.NewReviewAppSvc(c).GetCommodityReviews(req.CommodityId, pagination)
	if err != nil {
		app.NewResponse(c).Error(err)
		return
	}

	app.NewResponse(c).Success(reviews)
}

// GetReviewStatistics 获取商品评价统计
func GetReviewStatistics(c *gin.Context) {
	var req struct {
		CommodityId int64 `form:"commodity_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	stats, err := appservice.NewReviewAppSvc(c).GetReviewStatistics(req.CommodityId)
	if err != nil {
		app.NewResponse(c).Error(err)
		return
	}

	app.NewResponse(c).Success(stats)
}

// AdminReviewReply 商家回复评价
func AdminReviewReply(c *gin.Context) {
	var req request.ReviewReply
	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	err := appservice.NewReviewAppSvc(c).AdminReviewReply(&req)
	if err != nil {
		app.NewResponse(c).Error(err)
		return
	}

	app.NewResponse(c).Success(nil)
}

// UpdateReviewStatus 更新评价状态
func UpdateReviewStatus(c *gin.Context) {
	var req struct {
		ReviewId uint `json:"review_id" binding:"required"`
		Status   int  `json:"status" binding:"required,oneof=0 1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	err := appservice.NewReviewAppSvc(c).UpdateReviewStatus(req.ReviewId, req.Status)
	if err != nil {
		app.NewResponse(c).Error(err)
		return
	}

	app.NewResponse(c).Success(nil)
} 