package do

import "time"

// 演示 DEMO，后期使用时删掉

type DemoOrder struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"userId"`
	BillMoney int64     `json:"billMoney"`
	OrderNo   string    `json:"orderNo"`
	State     int8      `json:"state"`
	IsDel     uint      `json:"is_del"`
	PaidAt    time.Time `json:"paidAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
