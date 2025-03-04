package request

// ReviewCreate 创建评价请求
type ReviewCreate struct {
	OrderNo     string   `json:"order_no" binding:"required"`           // 订单编号
	OrderItemId int64    `json:"order_item_id" binding:"required"`      // 订单商品ID
	CommodityId int64    `json:"commodity_id" binding:"required"`       // 商品ID
	Rating      int      `json:"rating" binding:"required,min=1,max=5"` // 评分(1-5)
	Content     string   `json:"content"`                               // 评价内容
	IsAnonymous bool     `json:"is_anonymous"`                          // 是否匿名评价
	Images      []string `json:"images"`                                // 评价图片
}

// ReviewReply 商家回复评价请求
type ReviewReply struct {
	ReviewId int64  `json:"review_id" binding:"required"` // 评价ID
	Reply    string `json:"reply" binding:"required"`     // 回复内容
}
