package reply

type DemoOrder struct {
	UserId    int64  `json:"userId"`
	BillMoney int64  `json:"billMoney"`
	OrderNo   string `json:"orderNo"`
	State     int8   `json:"state"`
	PaidAt    string `json:"paidAt"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
