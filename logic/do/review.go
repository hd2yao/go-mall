package do

import "time"

// Review 评价领域对象
type Review struct {
	ID             int64
	UserId         int64
	OrderNo        string
	OrderItemId    int64
	CommodityId    int64
	Rating         int
	Content        string
	IsAnonymous    bool
	HasImage       bool
	AdminReply     string
	AdminReplyTime *int64
	Status         int
	Images         []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ReviewStatistics 商品评价统计
type ReviewStatistics struct {
	CommodityId   int64
	TotalCount    int     // 总评价数
	PositiveCount int     // 好评数(4-5星)
	NeutralCount  int     // 中评数(3星)
	NegativeCount int     // 差评数(1-2星)
	AverageRating float64 // 平均评分
	HasImageCount int     // 有图评价数
}
