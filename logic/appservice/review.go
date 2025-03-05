package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/logic/do"
	"github.com/hd2yao/go-mall/logic/domainservice"
)

type ReviewAppSvc struct {
	ctx             context.Context
	reviewDomainSvc *domainservice.ReviewDomainSvc
}

func NewReviewAppSvc(ctx context.Context) *ReviewAppSvc {
	return &ReviewAppSvc{
		ctx:             ctx,
		reviewDomainSvc: domainservice.NewReviewDomainSvc(ctx),
	}
}

// CreateReview 创建商品评价
func (ras *ReviewAppSvc) CreateReview(request *request.ReviewCreate, userId int64) error {
	orderModel, err := domainservice.NewOrderDomainSvc(ras.ctx).GetSpecifiedUserOrder(request.OrderNo, userId)
	if err != nil {
		return err
	}
	if orderModel == nil || orderModel.UserId != userId { // 订单不存在或不是当前用户
		return errcode.ErrReviewParams
	}
	if orderModel.OrderStatus != enum.OrderStatusConfirmReceipt && orderModel.OrderStatus != enum.OrderStatusCompleted { // 订单状态不是确认收货或已完成
		return errcode.ErrReviewUnsupportedScene
	}

	review := new(do.Review)
	if err = util.CopyProperties(review, request); err != nil {
		return errcode.ErrCoverData
	}
	review.UserId = userId
	review.Status = 0 // 待审核状态

	err = ras.reviewDomainSvc.CreateReview(review)
	if err != nil {
		return err
	}
	return nil
}

// GetReviewById 获取评价详情
func (ras *ReviewAppSvc) GetReviewById(reviewId int64, userId int64) (*reply.Review, error) {
	review, err := ras.reviewDomainSvc.GetReviewById(reviewId, userId)
	if err != nil {
		return nil, err
	}

	replyReview := new(reply.Review)
	if err = util.CopyProperties(replyReview, review); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return replyReview, nil
}

// GetUserReviews 获取用户的评价列表
func (ras *ReviewAppSvc) GetUserReviews(userId int64, pagination *app.Pagination) ([]*reply.Review, error) {
	reviews, err := ras.reviewDomainSvc.GetUserReviews(userId, pagination)
	if err != nil {
		return nil, err
	}

	replyReviews := make([]*reply.Review, 0, len(reviews))
	if err = util.CopyProperties(&replyReviews, reviews); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return replyReviews, nil
}

// GetCommodityReviews 获取商品的评价列表
func (ras *ReviewAppSvc) GetCommodityReviews(commodityId int64, pagination *app.Pagination) ([]*reply.Review, error) {
	commodityModel := domainservice.NewCommodityDomainSvc(ras.ctx).GetCommodityInfo(commodityId)
	if commodityModel == nil {
		return nil, errcode.ErrReviewParams
	}

	reviews, err := ras.reviewDomainSvc.GetCommodityReviews(commodityId, pagination)
	if err != nil {
		return nil, err
	}

	replyReviews := make([]*reply.Review, 0, len(reviews))
	if err = util.CopyProperties(&replyReviews, reviews); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return replyReviews, nil
}

// GetReviewStatistics 获取商品评价统计
func (ras *ReviewAppSvc) GetReviewStatistics(commodityId int64) (*reply.ReviewStatistics, error) {
	commodityModel := domainservice.NewCommodityDomainSvc(ras.ctx).GetCommodityInfo(commodityId)
	if commodityModel == nil {
		return nil, errcode.ErrReviewParams
	}

	stats, err := ras.reviewDomainSvc.GetReviewStatistics(commodityId)
	if err != nil {
		return nil, err
	}

	replyStats := new(reply.ReviewStatistics)
	if err = util.CopyProperties(replyStats, stats); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return replyStats, nil
}

// AdminReviewReply 商家回复评价
func (ras *ReviewAppSvc) AdminReviewReply(req *request.ReviewReply) error {
	return ras.reviewDomainSvc.AdminReviewReply(req.ReviewId, req.Reply)
}

// UpdateReviewStatus 更新评价状态
func (ras *ReviewAppSvc) UpdateReviewStatus(reviewId int64, status int) error {
	return ras.reviewDomainSvc.UpdateReviewStatus(reviewId, status)
}
