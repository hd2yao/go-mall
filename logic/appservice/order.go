package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
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

// CreateOrder 创建订单
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

// GetUserOrders 查询用户订单
func (oas *OrderAppSvc) GetUserOrders(userId int64, pagination *app.Pagination) ([]*reply.Order, error) {
	orders, err := oas.orderDomainSvc.GetUserOrders(userId, pagination)
	if err != nil {
		return nil, err
	}

	replyOrders := make([]*reply.Order, 0, len(orders))
	if err = util.CopyProperties(&replyOrders, &orders); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	for _, replyOrder := range replyOrders {
		// 订单的前台状态
		replyOrder.FrontStatus = enum.OrderFrontStatus[replyOrder.OrderStatus]
		replyOrder.Address.UserName = util.MaskPhone(replyOrder.Address.UserName)
		replyOrder.Address.UserPhone = util.MaskPhone(replyOrder.Address.UserPhone)
	}
	return replyOrders, nil
}

// GetOrderInfo 订单详情
func (oas *OrderAppSvc) GetOrderInfo(orderNo string, userId int64) (*reply.Order, error) {
	order, err := oas.orderDomainSvc.GetSpecifiedUserOrder(orderNo, userId)
	if err != nil {
		return nil, err
	}

	replyOrder := new(reply.Order)
	if err = util.CopyProperties(replyOrder, order); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	// 订单的前台状态
	replyOrder.FrontStatus = enum.OrderFrontStatus[replyOrder.OrderStatus]
	// 敏感信息脱敏
	replyOrder.Address.UserName = util.MaskRealName(replyOrder.Address.UserName)
	replyOrder.Address.UserPhone = util.MaskPhone(replyOrder.Address.UserPhone)

	return replyOrder, nil
}

// CancelOrder 取消订单
func (oas *OrderAppSvc) CancelOrder(orderNo string, userId int64) error {
	return oas.orderDomainSvc.CancelUserOrder(orderNo, userId)
}
