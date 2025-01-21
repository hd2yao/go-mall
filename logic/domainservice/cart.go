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
func (cds *CartDomainSvc) CartAddItem(cartItem *do.ShoppingCartItem) error {
	cartItemModel, err := cds.cartDao.GetUserCartItemWithCommodityId(cartItem.UserId, cartItem.CommodityId)
	if err != nil {
		return errcode.Wrap("CartAddItemError", err)
	}

	// 购物车中已存在该商品
	if cartItemModel != nil && cartItemModel.CartItemId != 0 {
		cartItemModel.CommodityNum += cartItem.CommodityNum
		return cds.cartDao.UpdateCartItem(cartItemModel)
	}

	err = util.CopyProperties(cartItemModel, cartItem)
	if err != nil {
		return errcode.ErrCoverData
	}

	return cds.cartDao.AddCartItem(cartItemModel)
}

// GetCheckedCartItems 获取选中的购物项
func (cds *CartDomainSvc) GetCheckedCartItems(cartItemIds []int64, userId int64) ([]*do.ShoppingCartItem, error) {
	cartItemModels, err := cds.cartDao.FindCartItems(cartItemIds)
	if err != nil {
		return nil, errcode.Wrap("GetCheckedCartItemsError", err)
	}

	// 确保购物项归属用户与请求用户一致

}
