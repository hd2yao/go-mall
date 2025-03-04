package domainservice

import (
	"context"
	"time"

	"github.com/samber/lo"

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
func (rds *ReviewDomainSvc) CreateReview(review *do.Review) (err error) {
	// 验证评分范围
	if review.Rating < 1 || review.Rating > 5 {
		return errcode.ErrParams
	}

	tx := dao.DBMaster().Begin()
	panicked := true
	defer func() {
		if err != nil || panicked {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	review.HasImage = len(review.Images) > 0
	err = rds.reviewDao.CreateReview(tx, review, review.Images)
	if err != nil {
		return errcode.Wrap("CreateReviewError", err)
	}

	panicked = false
	return nil
}

// GetReviewById 获取评价详情
func (rds *ReviewDomainSvc) GetReviewById(reviewId int64, userId int64) (*do.Review, error) {
	reviewModel, err := rds.reviewDao.GetReviewById(reviewId)
	if err != nil {
		return nil, errcode.Wrap("GetReviewByIdError", err)
	}
	if reviewModel == nil || reviewModel.UserId != userId {
		return nil, errcode.ErrReviewParams
	}

	review := new(do.Review)
	err = util.CopyProperties(review, reviewModel)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	// 获取评价图片
	reviewImages, err := rds.reviewDao.GetReviewImages(reviewId)
	if err != nil {
		return nil, errcode.Wrap("GetReviewByIdError", err)
	}
	reviewImageUrls := lo.Map(reviewImages, func(reviewImage *model.ReviewImage, index int) string {
		return reviewImage.ImageUrl
	})
	review.Images = reviewImageUrls

	return review, nil
}

// GetUserReviews 获取用户的评价列表
func (rds *ReviewDomainSvc) GetUserReviews(userId int64, pagination *app.Pagination) ([]*do.Review, error) {
	offset := pagination.Offset()
	size := pagination.GetPageSize()

	reviewModels, totalRow, err := rds.reviewDao.GetUserReviews(userId, offset, size)
	if err != nil {
		return nil, errcode.Wrap("GetUserReviewsError", err)
	}
	pagination.SetTotalRows(int(totalRow))
	reviews := make([]*do.Review, 0, len(reviewModels))
	if err = util.CopyProperties(&reviews, reviewModels); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	// 提取所有评价 ID
	reviewIds := lo.Map(reviewModels, func(reviewModel *model.Review, index int) int64 {
		return reviewModel.ID
	})
	// 获取评价图片
	reviewImages, err := rds.reviewDao.GetMultiReviewsImages(reviewIds)
	if err != nil {
		return nil, errcode.Wrap("GetUserReviewsError", err)
	}

	// 填充 Review 中的 Images
	for _, review := range reviews {
		review.Images = reviewImages[review.ID]
	}

	return reviews, nil
}

// GetCommodityReviews 获取商品的评价列表
func (rds *ReviewDomainSvc) GetCommodityReviews(commodityId int64, pagination *app.Pagination) ([]*do.Review, error) {
	offset := pagination.Offset()
	size := pagination.GetPageSize()

	// 获取商品评价列表
	reviewModels, totalRow, err := rds.reviewDao.GetCommodityReviews(commodityId, offset, size)
	if err != nil {
		return nil, errcode.Wrap("GetCommodityReviewsError", err)
	}
	pagination.SetTotalRows(int(totalRow))
	reviews := make([]*do.Review, 0, len(reviewModels))
	if err = util.CopyProperties(&reviews, reviewModels); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	// 提取所有评价 ID
	reviewIds := lo.Map(reviewModels, func(reviewModel *model.Review, index int) int64 {
		return reviewModel.ID
	})
	// 获取评价图片
	reviewImages, err := rds.reviewDao.GetMultiReviewsImages(reviewIds)
	if err != nil {
		return nil, errcode.Wrap("GetCommodityReviewsError", err)
	}

	// 填充 Review 中的 Images
	for _, review := range reviews {
		review.Images = reviewImages[review.ID]
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
		CommodityId:   commodityId,
		TotalCount:    stats.TotalCount,
		PositiveCount: stats.PositiveCount,
		NeutralCount:  stats.NeutralCount,
		NegativeCount: stats.NegativeCount,
		HasImageCount: stats.HasImageCount,
		AverageRating: stats.AvgRating,
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
