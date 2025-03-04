package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// ReviewImage 评价图片表
type ReviewImage struct {
	ID        int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 评价图片关联评价主键 ID
	ReviewId  int64                 `gorm:"column:review_id;not null;index:idx_review_id"`        // 评价 ID
	ImageUrl  string                `gorm:"column:image_url;not null"`                            // 图片 URL
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 0-未删除 1-已删除
	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

func (ri *ReviewImage) TableName() string {
	return "review_images"
}
