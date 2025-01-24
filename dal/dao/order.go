package dao

import "context"

type OrderDao struct {
	ctx context.Context
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{ctx: ctx}
}
