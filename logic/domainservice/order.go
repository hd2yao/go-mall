package domainservice

import (
	"context"

	"github.com/samber/lo"

	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/dao"
	"github.com/hd2yao/go-mall/logic/do"
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

// CreateOrder 创建订单
func (ods *OrderDomainSvc) CreateOrder(items []*do.ShoppingCartItem, userAddress *do.UserAddressInfo) (*do.Order, error) {
	// 计算订单商品的总价、优惠金额等结算信息
	billInfo, err := NewCartBillChecker(items, userAddress.UserId).GetBill()
	if err != nil {
		return nil, errcode.Wrap("CreateOrderError", err)
	}
	if billInfo.OriginalTotalPrice <= 0 {
		return nil, errcode.ErrCartItemParam
	}

	order := do.OrderNew()
	order.UserId = userAddress.UserId
	order.OrderNo = util.GenOrderNo(order.UserId)
	order.BillMoney = billInfo.OriginalTotalPrice
	order.PayMoney = billInfo.TotalPrice
	order.OrderStatus = enum.OrderStatusCreated
	if err = util.CopyProperties(&order.Items, &items); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	if err = util.CopyProperties(&order.Address, &userAddress); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	// 手动开启事务
	tx := dao.DBMaster().Begin()
	panicked := true
	defer func() { // 控制事务的提交和回滚，保证事务的完整性
		// db.Transaction 内部其实就是这么实现的
		if err != nil || panicked {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 下面的步骤如果很多可以再使用责任链模式把步骤组织起来
	// 1. 创建订单
	err = ods.orderDao.CreateOrder(tx, order)
	if err != nil {
		return nil, err
	}
	// 2. 删除购物车中的购买的购物项
	cartDao := dao.NewCartDao(ods.ctx)
	cartItemIds := lo.Map(items, func(item *do.ShoppingCartItem, index int) int64 {
		return item.CartItemId
	})
	err = cartDao.DeleteMultiCartItemInTx(tx, cartItemIds)
	if err != nil {
		return nil, err
	}
	// 3. 记录 Coupon 使用信息 并 锁定优惠券，等支付成功后再核销
	if billInfo.Coupon.CouponId > 0 {
		// couponDao.LockCoupon(tx, coupon)
	}
	// 4. 记录满减券使用信息
	if billInfo.Discount.DiscountId > 0 {
		// discountDao.recordDiscount(tx, discount)
	}
	// 5. 减少订单购买商品的库存 -- 会锁行记录，把这一步放到创建订单步骤的最后，减少行记录加锁的时间
	commodityDao := dao.NewCommodityDao(ods.ctx)
	err = commodityDao.ReduceStuckInOrderCreate(tx, order.Items)
	if err != nil {
		return nil, err
	}

	// 记得设置，让事务能正常提交
	panicked = false

	return order, nil
}

// GetUserOrders 获取用户订单
func (ods *OrderDomainSvc) GetUserOrders(userId int64, pagination *app.Pagination) ([]*do.Order, error) {
	offset := pagination.Offset()
	size := pagination.GetPageSize()

	// 查询用户订单
	orderModels, totalRow, err := ods.orderDao.GetUserOrders(userId, offset, size)
	if err != nil {
		return nil, errcode.Wrap("GetUserOrdersError", err)
	}
	pagination.SetTotalRows(int(totalRow))
	orders := make([]*do.Order, 0, len(orderModels))
	if err = util.CopyProperties(&orders, &orderModels); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	// 提取所有订单 ID
	orderIds := lo.Map(orders, func(order *do.Order, index int) int64 {
		return order.ID
	})
	// 查询订单的地址
	ordersAddressMap, err := ods.orderDao.GetMultiOrdersAddress(orderIds)
	if err != nil {
		return nil, errcode.Wrap("GetUserOrdersError", err)
	}
	// 查询订单明细
	ordersItemMap, err := ods.orderDao.GetMultiOrdersItems(orderIds)
	if err != nil {
		return nil, errcode.Wrap("GetUserOrdersError", err)
	}

	// 填充 Order 中的 Address 和 Items
	for _, order := range orders {
		order.Address = new(do.OrderAddress) // 先初始化
		if err = util.CopyProperties(order.Address, ordersAddressMap[order.ID]); err != nil {
			return nil, errcode.ErrCoverData.WithCause(err)
		}
		orderItems := ordersItemMap[order.ID]
		if err = util.CopyProperties(&order.Items, &orderItems); err != nil {
			return nil, errcode.ErrCoverData.WithCause(err)
		}
	}

	return orders, nil
}
