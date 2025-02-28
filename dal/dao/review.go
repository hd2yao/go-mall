package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/dal/model"
)

type ReviewDao struct {
	ctx context.Context
}

func NewReviewDao(ctx context.Context) *ReviewDao {
	return &ReviewDao{ctx: ctx}
}

// CreateReview 创建评价
func (rd *ReviewDao) CreateReview(review *model.Review, images []string) error {
	return DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(review).Error; err != nil {
			return err
		}

		if len(images) > 0 {
			reviewImages := make([]model.ReviewImage, 0, len(images))
			for _, img := range images {
				reviewImages = append(reviewImages, model.ReviewImage{
					ReviewId: int64(review.ID),
					ImageUrl: img,
				})
			}
			if err := tx.Create(&reviewImages).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetReviewById 根据ID获取评价
func (rd *ReviewDao) GetReviewById(reviewId uint) (*model.Review, []string, error) {
	var review model.Review
	var images []model.ReviewImage

	err := DB().First(&review, reviewId).Error
	if err != nil {
		return nil, nil, err
	}

	if review.HasImage {
		if err := DB().Where("review_id = ?", reviewId).Find(&images).Error; err != nil {
			return nil, nil, err
		}
	}

	imageUrls := make([]string, 0, len(images))
	for _, img := range images {
		imageUrls = append(imageUrls, img.ImageUrl)
	}

	return &review, imageUrls, nil
}

// GetUserReviews 获取用户的评价列表
func (rd *ReviewDao) GetUserReviews(userId int64, pagination *app.Pagination) ([]*model.Review, error) {
	var reviews []*model.Review
	offset := (pagination.Page - 1) * pagination.Size

	err := DB().Where("user_id = ? AND status = ?", userId, 1).
		Offset(offset).
		Limit(pagination.Size).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, err
}

// GetCommodityReviews 获取商品的评价列表
func (rd *ReviewDao) GetCommodityReviews(commodityId int64, pagination *app.Pagination) ([]*model.Review, error) {
	var reviews []*model.Review
	offset := (pagination.Page - 1) * pagination.Size

	err := DB().Where("commodity_id = ? AND status = ?", commodityId, 1).
		Offset(offset).
		Limit(pagination.Size).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, err
}

// GetReviewStatistics 获取商品评价统计
func (rd *ReviewDao) GetReviewStatistics(commodityId int64) (*struct {
	TotalCount    int
	PositiveCount int
	NeutralCount  int
	NegativeCount int
	HasImageCount int
	AvgRating     float64
}, error) {
	var stats struct {
		TotalCount    int
		PositiveCount int
		NeutralCount  int
		NegativeCount int
		HasImageCount int
		AvgRating     float64
	}

	err := DB().Model(&model.Review{}).
		Where("commodity_id = ? AND status = ?", commodityId, 1).
		Select("COUNT(*) as total_count, " +
			"SUM(CASE WHEN rating >= 4 THEN 1 ELSE 0 END) as positive_count, " +
			"SUM(CASE WHEN rating = 3 THEN 1 ELSE 0 END) as neutral_count, " +
			"SUM(CASE WHEN rating <= 2 THEN 1 ELSE 0 END) as negative_count, " +
			"SUM(CASE WHEN has_image = 1 THEN 1 ELSE 0 END) as has_image_count, " +
			"AVG(rating) as avg_rating").
		Scan(&stats).Error

	return &stats, err
}

// UpdateReviewStatus 更新评价状态
func (rd *ReviewDao) UpdateReviewStatus(reviewId uint, status int) error {
	return DB().Model(&model.Review{}).
		Where("id = ?", reviewId).
		Update("status", status).Error
}

// AdminReply 商家回复评价
func (rd *ReviewDao) AdminReply(reviewId uint, reply string, replyTime int64) error {
	return DB().Model(&model.Review{}).
		Where("id = ?", reviewId).
		Updates(map[string]interface{}{
			"admin_reply":      reply,
			"admin_reply_time": replyTime,
		}).Error
}
