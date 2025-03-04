package dao

import (
	"context"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/model"
	"github.com/hd2yao/go-mall/logic/do"
)

type ReviewDao struct {
	ctx context.Context
}

func NewReviewDao(ctx context.Context) *ReviewDao {
	return &ReviewDao{ctx: ctx}
}

// CreateReview 创建评价
func (rd *ReviewDao) CreateReview(tx *gorm.DB, review *do.Review, images []string) error {
	reviewModel := new(model.Review)
	err := util.CopyProperties(reviewModel, review)
	if err != nil {
		return errcode.ErrCoverData.WithCause(err)
	}

	err = tx.WithContext(rd.ctx).Create(reviewModel).Error
	if err != nil {
		return err
	}

	// 填充回 ID
	review.ID = reviewModel.ID

	// 处理评价图片
	if len(images) > 0 {
		reviewImages := make([]model.ReviewImage, 0, len(images))
		for _, img := range images {
			reviewImages = append(reviewImages, model.ReviewImage{
				ReviewId: review.ID,
				ImageUrl: img,
			})
		}
		if err = tx.WithContext(rd.ctx).Create(&reviewImages).Error; err != nil {
			return err
		}
	}
	return nil

}

// GetReviewById 根据评论 ID 获取评价
func (rd *ReviewDao) GetReviewById(reviewId int64) (*model.Review, error) {
	review := new(model.Review)
	err := DB().WithContext(rd.ctx).First(&review, reviewId).Error
	return review, err
}

// GetReviewImages 根据 reviewId 获取评价图片
func (rd *ReviewDao) GetReviewImages(reviewId int64) ([]*model.ReviewImage, error) {
	reviewImages := make([]*model.ReviewImage, 0)
	err := DB().WithContext(rd.ctx).Where("review_id = ?", reviewId).
		Find(&reviewImages).Error
	return reviewImages, err
}

// GetUserReviews 获取用户的评价列表
func (rd *ReviewDao) GetUserReviews(userId int64, offset, returnSize int) (reviews []*model.Review, totalRows int64, err error) {
	err = DB().WithContext(rd.ctx).Where("user_id = ? AND status = ?", userId, 1).
		Offset(offset).Limit(returnSize).
		Order("created_at DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询满足条件的记录数
	DB().WithContext(rd.ctx).Model(&model.Review{}).Where("user_id = ? AND status = ?", userId, 1).Count(&totalRows)
	return
}

// GetMultiReviewsImages 获取多条评论的图片, 返回以 reviewId 为 Key, 对应的评论图片为值的 Map
func (rd *ReviewDao) GetMultiReviewsImages(reviewIds []int64) (map[int64][]string, error) {
	// 1. 查询所有评论的图片数据
	reviewImageModel := make([]*model.ReviewImage, 0)
	err := DB().WithContext(rd.ctx).Where("review_id IN ?", reviewIds).
		Find(&reviewImageModel).Error
	if err != nil {
		return nil, err
	}

	// 2. 使用 lo.GroupBy 将图片按照 reviewId 分组
	// 将 []*model.ReviewImage 转换为 map[int64][]*model.ReviewImage
	reviewImagesMap := lo.GroupBy(reviewImageModel, func(image *model.ReviewImage) int64 {
		return image.ReviewId
	})

	// 3. 使用 lo.MapValues 处理 map 的值部分
	// 将 map[int64][]*model.ReviewImage 转换为 map[int64][]string
	reviewImages := lo.MapValues(reviewImagesMap, func(images []*model.ReviewImage, index int64) []string {
		// 使用 lo.Map 将每组图片记录转换为 URL 数组
		return lo.Map(images, func(image *model.ReviewImage, index int) string {
			return image.ImageUrl
		})
	})

	// 使用 for 循环实现
	//reviewImagesMap := make(map[int64][]string)
	//for _, image := range reviewImages {
	//	reviewImagesMap[image.ReviewId] = append(reviewImagesMap[image.ReviewId], image.ImageUrl)
	//}

	return reviewImages, nil
}

// GetCommodityReviews 获取商品的评价列表
func (rd *ReviewDao) GetCommodityReviews(commodityId int64, offset, returnSize int) (reviews []*model.Review, totalRows int64, err error) {
	err = DB().WithContext(rd.ctx).Where("commodity_id = ? AND status = ?", commodityId, 1).
		Offset(offset).Limit(returnSize).
		Order("created_at DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询满足条件的记录数
	DB().WithContext(rd.ctx).Model(model.Review{}).Where("commodity_id = ?", commodityId).Count(&totalRows)
	return
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
