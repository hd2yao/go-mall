package domainservice

import (
	"math"

	"github.com/samber/lo"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/logic/do"
)

type CartBillChecker struct {
	UserId        int64
	checkingItems []*do.ShoppingCartItem
	Coupon        struct { // 可用的优惠券
		CouponId      int64
		CouponName    string
		DiscountMoney int // 减免金额, 单位: 分
		Threshold     int // 使用门槛, 单位: 分, 设置成 1000 表示满10元可用
	}
	Discount struct { // 可用的满减活动券
		DiscountId    int64
		DiscountName  string
		DiscountMoney int
		Threshold     int
	}
	VipOffRate int // VIP的折扣  8 折  = 20% off

	handler cartBillCheckHandler
}

func NewCartBillChecker(items []*do.ShoppingCartItem, userId int64) *CartBillChecker {
	checker := new(CartBillChecker)
	checker.UserId = userId
	checker.checkingItems = items
	checker.handler = &checkerStarter{}
	// 通过责任链设置 要检查的各种优惠项
	checker.handler.SetNext(&couponChecker{}).
		SetNext(&discountChecker{}).
		SetNext(&vipChecker{})
	return checker
}

// GetBill 获取账单信息
func (cbc *CartBillChecker) GetBill() (*do.CartBillInfo, error) {
	err := cbc.handler.RunChecker(cbc)
	if err != nil {
		return nil, errcode.Wrap("CartBillCheckerError", err)
	}

	// 计算商品使用减免前的总价
	originalTotalPrice := lo.Reduce(cbc.checkingItems, func(agg int, item *do.ShoppingCartItem, index int) int {
		return agg + item.CommoditySellingPrice*item.CommodityNum
	}, 0)

	// VIP 能减免的金额
	vipDiscountMoney := int(math.Round(float64(originalTotalPrice) * float64(cbc.VipOffRate) / 100.0))
	totalPrice := originalTotalPrice - vipDiscountMoney

	// 满足优惠卷使用条件
	if cbc.Coupon.Threshold != 0 && originalTotalPrice > cbc.Coupon.Threshold {
		// 优惠券能减免的金额
		totalPrice -= cbc.Coupon.DiscountMoney
	}

	// 满足满减活动使用条件
	if cbc.Discount.Threshold != 0 && totalPrice > cbc.Discount.Threshold {
		// 满减活动能减免的金额
		totalPrice -= cbc.Discount.DiscountMoney
	}

	billInfo := new(do.CartBillInfo)
	billInfo.Coupon = cbc.Coupon
	billInfo.Discount = cbc.Discount
	billInfo.VipDiscountMoney = vipDiscountMoney
	billInfo.TotalPrice = totalPrice
	billInfo.OriginalTotalPrice = originalTotalPrice
	return billInfo, nil
}

type cartBillCheckHandler interface {
	RunChecker(*CartBillChecker) error
	SetNext(cartBillCheckHandler) cartBillCheckHandler
	Check(*CartBillChecker) error
}

// 充当抽象类型，实现公共方法，抽象方法留给实现类自己实现
type cartCommonChecker struct {
	nextHandler cartBillCheckHandler
}

func (n *cartCommonChecker) SetNext(handler cartBillCheckHandler) cartBillCheckHandler {
	n.nextHandler = handler
	return handler
}

// RunChecker 启动责任链，并传递给 nextHandler
// 执行 nextHandler，若无误，则调用 RunChecker 进行传递
func (n *cartCommonChecker) RunChecker(billChecker *CartBillChecker) error {
	if n.nextHandler != nil {
		if err := n.nextHandler.Check(billChecker); err != nil {
			err = errcode.Wrap("cartCommonChecker", err)
			return err
		}
		return n.nextHandler.RunChecker(billChecker)
	}
	return nil
}

//func (n *cartCommonChecker) Check(cbc *CartBillChecker) error {
//	fmt.Println("1111111111111111111111")
//	return nil
//}

type checkerStarter struct {
	cartCommonChecker
}

// Check 空方法，什么也不做，目的是让抽象类的 RunChecker 能启动调用链
func (cs *checkerStarter) Check(cbc *CartBillChecker) error {
	// 这里是因为，RunChecker 是从 nextHandler 开始的，所以设置一个 checkerStarter
	// 可简单理解为 dummy head，链表中的虚拟头节点
	return nil
}

// couponChecker 优惠券 checker
type couponChecker struct {
	cartCommonChecker
}

// Check 检查优惠券，如果用户有可用优惠券，则设置到 CartBillChecker 中，在本项目中可理解为执行 execute
func (cc *couponChecker) Check(cbc *CartBillChecker) error {
	// TODO: 查询用户是否有可用优惠券
	// 这里是 Mock 逻辑
	// dao.GetUserCoupon(cbc.UserId)
	cbc.Coupon = struct {
		CouponId      int64
		CouponName    string
		DiscountMoney int
		Threshold     int
	}{
		CouponId:      1,
		DiscountMoney: 100, // 假设可用优惠券 ID 为 1， 减免 100
		Threshold:     100,
	}
	return nil
}

// discountChecker 折扣减免 checker
type discountChecker struct {
	cartCommonChecker
}

// Check 检查折扣减免，如果用户有可用折扣减免，则设置到 CartBillChecker 中，在本项目中可理解为执行 execute
func (dc *discountChecker) Check(cbc *CartBillChecker) error {
	// TODO: 查用户是否有可用的减免活动
	// 这里是Mock逻辑
	// dao.GetApplicableDiscount(cbc.UserId)
	cbc.Discount = struct {
		DiscountId    int64
		DiscountName  string
		DiscountMoney int
		Threshold     int
	}{
		DiscountId:    1,
		DiscountMoney: 100, // 假设可用优惠券ID为1， 减免100
		Threshold:     1000,
	}
	return nil
}

// vipChecker VIP checker
type vipChecker struct {
	cartCommonChecker
}

// Check 检查 VIP，如果用户是 VIP，则设置到 CartBillChecker 中，在本项目中可理解为执行 execute
func (vc *vipChecker) Check(cbc *CartBillChecker) error {
	// TODO: 判断用户是不是会员, 有没有会员折扣
	//if isVip(userId) {
	//  cbc.VipOffRate = 12 // 比如vip减免12%
	//  return nil
	//}
	cbc.VipOffRate = 0 // 不是vip不减免
	return nil
}
