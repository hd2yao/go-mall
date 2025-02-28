package domainservice

import (
	"context"
	"time"

	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/dao"
	"github.com/hd2yao/go-mall/dal/model"
	"github.com/hd2yao/go-mall/logic/do"
)

type ReviewDomainSvc struct {
	ctx       context.Context
	reviewDao *dao.ReviewDao
}

func NewReviewDomainSvc(ctx context.Context) *ReviewDomainSvc {
	return &ReviewDomainSvc{
		ctx:       ctx,
		reviewDao: dao.NewReviewDao(ctx),
	}
}

// CreateReview 创建商品评价
func (rds *ReviewDomainSvc) CreateReview(review *do.Review) error {
	// 验证评分范围
	if review.Rating < 1 || review.Rating > 5 {
		return errcode.ErrParams
	}

	reviewModel := new(model.Review)
	err := util.CopyProperties(reviewModel, review)
	if err != nil {
		return errcode.ErrCoverData
	}

	reviewModel.HasImage = len(review.Images) > 0
	err = rds.reviewDao.CreateReview(reviewModel, review.Images)
	if err != nil {
		return errcode.Wrap("CreateReviewError", err)
	}

	return nil
}

// GetReviewById 获取评价详情
func (rds *ReviewDomainSvc) GetReviewById(reviewId uint) (*do.Review, error) {
	reviewModel, images, err := rds.reviewDao.GetReviewById(reviewId)
	if err != nil {
		return nil, errcode.Wrap("GetReviewError", err)
	}

	review := new(do.Review)
	err = util.CopyProperties(review, reviewModel)
	if err != nil {
		return nil, errcode.ErrCoverData
	}
	review.Images = images

	return review, nil
}

// GetUserReviews 获取用户的评价列表
func (rds *ReviewDomainSvc) GetUserReviews(userId int64, pagination *app.Pagination) ([]*do.Review, error) {
	reviewModels, err := rds.reviewDao.GetUserReviews(userId, pagination)
	if err != nil {
		return nil, errcode.Wrap("GetUserReviewsError", err)
	}

	reviews := make([]*do.Review, 0, len(reviewModels))
	err = util.CopyProperties(&reviews, reviewModels)
	if err != nil {
		return nil, errcode.ErrCoverData
	}

	return reviews, nil
}

// GetCommodityReviews 获取商品的评价列表
func (rds *ReviewDomainSvc) GetCommodityReviews(commodityId int64, pagination *app.Pagination) ([]*do.Review, error) {
	reviewModels, err := rds.reviewDao.GetCommodityReviews(commodityId, pagination)
	if err != nil {
		return nil, errcode.Wrap("GetCommodityReviewsError", err)
	}

	reviews := make([]*do.Review, 0, len(reviewModels))
	err = util.CopyProperties(&reviews, reviewModels)
	if err != nil {
		return nil, errcode.ErrCoverData
	}

	return reviews, nil
}

// GetReviewStatistics 获取商品评价统计
func (rds *ReviewDomainSvc) GetReviewStatistics(commodityId int64) (*do.ReviewStatistics, error) {
	stats, err := rds.reviewDao.GetReviewStatistics(commodityId)
	if err != nil {
		return nil, errcode.Wrap("GetReviewStatisticsError", err)
	}

	return &do.ReviewStatistics{
		CommodityId:    commodityId,
		TotalCount:     stats.TotalCount,
		PositiveCount:  stats.PositiveCount,
		NeutralCount:   stats.NeutralCount,
		NegativeCount:  stats.NegativeCount,
		HasImageCount:  stats.HasImageCount,
		AverageRating:  stats.AvgRating,
	}, nil
}

// AdminReviewReply 商家回复评价
func (rds *ReviewDomainSvc) AdminReviewReply(reviewId uint, reply string) error {
	replyTime := time.Now().Unix()
	err := rds.reviewDao.AdminReply(reviewId, reply, replyTime)
	if err != nil {
		return errcode.Wrap("AdminReviewReplyError", err)
	}
	return nil
}

// UpdateReviewStatus 更新评价状态
func (rds *ReviewDomainSvc) UpdateReviewStatus(reviewId uint, status int) error {
	err := rds.reviewDao.UpdateReviewStatus(reviewId, status)
	if err != nil {
		return errcode.Wrap("UpdateReviewStatusError", err)
	}
	return nil
} 