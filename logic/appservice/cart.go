package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/logic/do"
	"github.com/hd2yao/go-mall/logic/domainservice"
)

type CartAppSvc struct {
	ctx           context.Context
	cartDomainSvc *domainservice.CartDomainSvc
}

func NewCartAppSvc(ctx context.Context) *CartAppSvc {
	return &CartAppSvc{
		ctx:           ctx,
		cartDomainSvc: domainservice.NewCartDomainSvc(ctx),
	}
}

// AddCartItem 添加商品到购物车
func (cas *CartAppSvc) AddCartItem(request *request.AddCartItem, userId int64) error {
	commodityDomainSvc := domainservice.NewCommodityDomainSvc(cas.ctx)
	commodityInfo := commodityDomainSvc.GetCommodityInfo(request.CommodityId)
	if commodityInfo == nil || commodityInfo.ID == 0 { // 商品不存在
		return errcode.ErrCommodityNotExists
	}
	if commodityInfo.StockNum < request.CommodityNum {
		// 先初步判断库存是否充足, 下单时需要重新用当前读判断库存
		return errcode.ErrCommodityStockOut
	}

	shoppingCartItem := new(do.ShoppingCartItem)
	err := util.CopyProperties(shoppingCartItem, request)
	if err != nil {
		return errcode.ErrCoverData
	}
	shoppingCartItem.UserId = userId

	return cas.cartDomainSvc.CartAddItem(shoppingCartItem)
}

// CheckCartItemBill 查看购物项账单
func (cas *CartAppSvc) CheckCartItemBill(cartItemIds []int64, userId int64) (*reply.CheckedCartItemBill, error) {

}
