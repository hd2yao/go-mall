package domainservice

import (
	"context"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/dao"
	"github.com/hd2yao/go-mall/logic/do"
)

type CartDomainSvc struct {
	ctx     context.Context
	cartDao *dao.CartDao
}

func NewCartDomainSvc(ctx context.Context) *CartDomainSvc {
	return &CartDomainSvc{
		ctx:     ctx,
		cartDao: dao.NewCartDao(ctx),
	}
}

// CartAddItem 添加商品到购物车
func (c *CartDomainSvc) CartAddItem(cartItem *do.ShoppingCartItem) error {
	cartItemModel, err := c.cartDao.GetUserCartItemWithCommodityId(cartItem.UserId, cartItem.CommodityId)
	if err != nil {
		return errcode.Wrap("CartAddItemError", err)
	}

	// 购物车中已存在该商品
	if cartItemModel != nil && cartItemModel.CartItemId != 0 {
		cartItemModel.CommodityNum += cartItem.CommodityNum
		return c.cartDao.UpdateCartItem(cartItemModel)
	}

	err = util.CopyProperties(cartItemModel, cartItem)
	if err != nil {
		return errcode.ErrCoverData
	}

	return c.cartDao.AddCartItem(cartItemModel)
}
