package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// Review 商品评价表
type Review struct {
	ID             int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	UserId         int64                 `gorm:"column:user_id;not null;index:idx_user_id"`            // 用户 ID
	OrderId        int64                 `gorm:"column:order_id;not null;index:idx_order_id"`          // 订单 ID
	OrderItemId    int64                 `gorm:"column:order_item_id;not null"`                        // 订单商品 ID
	CommodityId    int64                 `gorm:"column:commodity_id;not null;index:idx_commodity_id"`  // 商品 ID
	Rating         int                   `gorm:"column:rating;not null"`                               // 评分(1-5)
	Content        string                `gorm:"column:content;type:text"`                             // 评价内容
	IsAnonymous    bool                  `gorm:"column:is_anonymous;default:false"`                    // 是否匿名评价
	HasImage       bool                  `gorm:"column:has_image;default:false"`                       // 是否包含图片
	AdminReply     string                `gorm:"column:admin_reply;type:text"`                         // 商家回复
	AdminReplyTime *int64                `gorm:"column:admin_reply_time"`                              // 商家回复时间
	Status         int                   `gorm:"column:status;default:0"`                              // 状态：0-待审核 1-已发布 2-已删除
	IsDel          soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 0-未删除 1-已删除
	CreatedAt      time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt      time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

// ReviewImage 评价图片表
type ReviewImage struct {
	ID        int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 评价图片关联评价主键 ID
	ReviewId  int64                 `gorm:"column:review_id;not null;index:idx_review_id"`        // 评价 ID
	ImageUrl  string                `gorm:"column:image_url;not null"`                            // 图片 URL
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 0-未删除 1-已删除
	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

// TableName 指定表名
func (r *Review) TableName() string {
	return "reviews"
}

func (ri *ReviewImage) TableName() string {
	return "review_images"
}
