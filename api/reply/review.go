package reply

import "time"

// Review 评价响应
type Review struct {
	ID             uint      `json:"id"`
	UserId         int64     `json:"user_id"`
	OrderId        int64     `json:"order_id"`
	OrderItemId    int64     `json:"order_item_id"`
	CommodityId    int64     `json:"commodity_id"`
	Rating         int       `json:"rating"`
	Content        string    `json:"content"`
	IsAnonymous   bool      `json:"is_anonymous"`
	HasImage      bool      `json:"has_image"`
	Images        []string  `json:"images"`
	AdminReply    string    `json:"admin_reply"`
	AdminReplyTime *int64    `json:"admin_reply_time"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// ReviewStatistics 商品评价统计响应
type ReviewStatistics struct {
	CommodityId    int64   `json:"commodity_id"`
	TotalCount     int     `json:"total_count"`      // 总评价数
	PositiveCount  int     `json:"positive_count"`   // 好评数(4-5星)
	NeutralCount   int     `json:"neutral_count"`    // 中评数(3星)
	NegativeCount  int     `json:"negative_count"`   // 差评数(1-2星)
	HasImageCount  int     `json:"has_image_count"`  // 有图评价数
	AverageRating  float64 `json:"average_rating"`   // 平均评分
} 