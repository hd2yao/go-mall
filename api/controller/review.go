package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/logic/appservice"
)

// CreateReview 创建商品评价
func CreateReview(c *gin.Context) {
	requestData := new(request.ReviewCreate)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	err := appservice.NewReviewAppSvc(c).CreateReview(requestData, c.GetInt64("user_id"))
	if err != nil {
		if errors.Is(err, errcode.ErrReviewParams) {
			app.NewResponse(c).Error(errcode.ErrReviewParams)
		} else if errors.Is(err, errcode.ErrReviewUnsupportedScene) {
			app.NewResponse(c).Error(errcode.ErrReviewUnsupportedScene)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).SuccessOk()
}

// GetReviewById 获取评价详情
func GetReviewById(c *gin.Context) {
	reviewIdStr := c.Param("review_id")
	reviewId, _ := strconv.ParseInt(reviewIdStr, 10, 64)

	replyReview, err := appservice.NewReviewAppSvc(c).GetReviewById(reviewId, c.GetInt64("user_id"))
	if err != nil {
		if errors.Is(err, errcode.ErrReviewParams) {
			app.NewResponse(c).Error(errcode.ErrReviewParams)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).Success(replyReview)
}

// GetUserReviews 获取用户的评价列表
func GetUserReviews(c *gin.Context) {
	pagination := app.NewPagination(c)
	replyReviews, err := appservice.NewReviewAppSvc(c).GetUserReviews(c.GetInt64("user_id"), pagination)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	app.NewResponse(c).Success(replyReviews)
}

// GetCommodityReviews 获取商品的评价列表
func GetCommodityReviews(c *gin.Context) {
	commodityIdStr := c.Param("commodity_id")
	commodityId, _ := strconv.ParseInt(commodityIdStr, 10, 64)
	pagination := app.NewPagination(c)

	replyReviews, err := appservice.NewReviewAppSvc(c).GetCommodityReviews(commodityId, pagination)
	if err != nil {
		if errors.Is(err, errcode.ErrReviewParams) {
			app.NewResponse(c).Error(errcode.ErrReviewParams)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).Success(replyReviews)
}

// GetReviewStatistics 获取商品评价统计
func GetReviewStatistics(c *gin.Context) {
	commodityIdStr := c.Param("commodity_id")
	commodityId, _ := strconv.ParseInt(commodityIdStr, 10, 64)

	commodityReviewStatistics, err := appservice.NewReviewAppSvc(c).GetReviewStatistics(commodityId)
	if err != nil {
		if errors.Is(err, errcode.ErrReviewParams) {
			app.NewResponse(c).Error(errcode.ErrReviewParams)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).Success(commodityReviewStatistics)
}

// AdminReviewReply 商家回复评价
func AdminReviewReply(c *gin.Context) {
	requestData := new(request.ReviewReply)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	err := appservice.NewReviewAppSvc(c).AdminReviewReply(requestData)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	app.NewResponse(c).SuccessOk()
}

// UpdateReviewStatus 审核更新评价状态
func UpdateReviewStatus(c *gin.Context) {
	requestData := new(request.ReviewApprove)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	err := appservice.NewReviewAppSvc(c).UpdateReviewStatus(requestData.ReviewId, requestData.Status)
	if err != nil {
		if errors.Is(err, errcode.ErrReviewStatusCanNotChanged) {
			app.NewResponse(c).Error(errcode.ErrReviewStatusCanNotChanged)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).SuccessOk()
}
