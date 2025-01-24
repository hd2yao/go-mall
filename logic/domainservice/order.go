package domainservice

import (
	"context"

	"github.com/hd2yao/go-mall/dal/dao"
)

type OrderDomainSvc struct {
	ctx      context.Context
	orderDao *dao.OrderDao
}

func NewOrderDomainSvc(ctx context.Context) *OrderDomainSvc {
	return &OrderDomainSvc{
		ctx:      ctx,
		orderDao: dao.NewOrderDao(ctx),
	}
}
