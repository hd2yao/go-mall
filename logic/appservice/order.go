package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/logic/domainservice"
)

type OrderAppSvc struct {
	ctx            context.Context
	orderDomainSvc *domainservice.OrderDomainSvc
}

func NewOrderAppSvc(ctx context.Context) *OrderAppSvc {
	return &OrderAppSvc{
		ctx:            ctx,
		orderDomainSvc: domainservice.NewOrderDomainSvc(ctx),
	}
}

func (oas *OrderAppSvc) CreateOrder(orderRequest *request.OrderCreate, userId int64) (*reply.OrderCreateReply, error) {
	// 通过购物项 ID 获取用户添加在购物车中的购物项
	cartDomainSvc := domainservice.NewCartDomainSvc(oas.ctx)
	cartItems, err := cartDomainSvc.GetCheckedCartItems(orderRequest.CartItemIdList, userId)
	if err != nil {
		return nil, err
	}

	// 通过用户地址 ID 获取用户地址
	userDomainSvc := domainservice.NewUserDomainSvc(oas.ctx)
	address, err := userDomainSvc.GetUserSingleAddress(userId, orderRequest.UserAddressId)
	if err != nil {
		return nil, err
	}

	// 创建订单
	order, err := oas.orderDomainSvc.CreateOrder(cartItems, address)
	if err != nil {
		return nil, err
	}

	orderReply := new(reply.OrderCreateReply)
	orderReply.OrderNo = order.OrderNo
	return orderReply, nil
}
