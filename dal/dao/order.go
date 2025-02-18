package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/model"
	"github.com/hd2yao/go-mall/logic/do"
)

type OrderDao struct {
	ctx context.Context
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{ctx: ctx}
}

func (od *OrderDao) CreateOrder(tx *gorm.DB, order *do.Order) error {
	orderModel := new(model.Order)
	err := util.CopyProperties(orderModel, order)
	if err != nil {
		return errcode.ErrCoverData.WithCause(err)
	}

	err = tx.WithContext(od.ctx).Create(orderModel).Error
	if err != nil {
		return err
	}

	// 填充 orderId
	order.ID = orderModel.ID
	order.Address.OrderId = orderModel.ID
	for _, item := range order.Items {
		item.OrderId = orderModel.ID
	}

	// 创建订单项
	err = od.createOrderItems(tx, order.Items)
	if err != nil {
		return err
	}

	// 创建订单地址
	err = od.createOrderAddress(tx, order.Address)
	return err
}

func (od *OrderDao) createOrderItems(tx *gorm.DB, orderItems []*do.OrderItem) error {
	orderItemModels := make([]*model.OrderItem, 0, len(orderItems))
	err := util.CopyProperties(&orderItemModels, &orderItems)
	if err != nil {
		return errcode.ErrCoverData.WithCause(err)
	}
	return tx.WithContext(od.ctx).Create(orderItemModels).Error
}

func (od *OrderDao) createOrderAddress(tx *gorm.DB, orderAddress *do.OrderAddress) error {
	orderAddressModel := new(model.OrderAddress)
	err := util.CopyProperties(orderAddressModel, orderAddress)
	if err != nil {
		return errcode.ErrCoverData.WithCause(err)
	}
	return tx.WithContext(od.ctx).Create(orderAddressModel).Error
}
