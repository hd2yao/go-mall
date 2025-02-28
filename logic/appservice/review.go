package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/logic/do"
	"github.com/hd2yao/go-mall/logic/domainservice"
)

type ReviewAppSvc struct {
	ctx            context.Context
	reviewDomainSvc *domainservice.ReviewDomainSvc
}

func NewReviewAppSvc(ctx context.Context) *ReviewAppSvc {
	return &ReviewAppSvc{
		ctx:            ctx,
		reviewDomainSvc: domainservice.NewReviewDomainSvc(ctx),
	}
}

// CreateReview 创建商品评价
func (ras *ReviewAppSvc) CreateReview(req *request.ReviewCreate, userId int64) error {
	review := new(do.Review)
	err := util.CopyProperties(review, req)
	if err != nil {
		return errcode.ErrCoverData
	}
	review.UserId = userId
	review.Status = 0 // 待审核状态

	return ras.reviewDomainSvc.CreateReview(review)
}

// GetReviewById 获取评价详情
func (ras *ReviewAppSvc) GetReviewById(reviewId uint) (*reply.Review, error) {
	review, err := ras.reviewDomainSvc.GetReviewById(reviewId)
	if err != nil {
		return nil, err
	}

	replyReview := new(reply.Review)
	err = util.CopyProperties(replyReview, review)
	if err != nil {
		return nil, errcode.ErrCoverData
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
	err = util.CopyProperties(&replyReviews, reviews)
	if err != nil {
		return nil, errcode.ErrCoverData
	}

	return replyReviews, nil
}

// GetCommodityReviews 获取商品的评价列表
func (ras *ReviewAppSvc) GetCommodityReviews(commodityId int64, pagination *app.Pagination) ([]*reply.Review, error) {
	reviews, err := ras.reviewDomainSvc.GetCommodityReviews(commodityId, pagination)
	if err != nil {
		return nil, err
	}

	replyReviews := make([]*reply.Review, 0, len(reviews))
	err = util.CopyProperties(&replyReviews, reviews)
	if err != nil {
		return nil, errcode.ErrCoverData
	}

	return replyReviews, nil
}

// GetReviewStatistics 获取商品评价统计
func (ras *ReviewAppSvc) GetReviewStatistics(commodityId int64) (*reply.ReviewStatistics, error) {
	stats, err := ras.reviewDomainSvc.GetReviewStatistics(commodityId)
	if err != nil {
		return nil, err
	}

	replyStats := new(reply.ReviewStatistics)
	err = util.CopyProperties(replyStats, stats)
	if err != nil {
		return nil, errcode.ErrCoverData
	}

	return replyStats, nil
}

// AdminReviewReply 商家回复评价
func (ras *ReviewAppSvc) AdminReviewReply(req *request.ReviewReply) error {
	return ras.reviewDomainSvc.AdminReviewReply(req.ReviewId, req.Reply)
}

// UpdateReviewStatus 更新评价状态
func (ras *ReviewAppSvc) UpdateReviewStatus(reviewId uint, status int) error {
	return ras.reviewDomainSvc.UpdateReviewStatus(reviewId, status)
} 