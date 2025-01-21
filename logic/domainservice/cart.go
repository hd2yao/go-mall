package domainservice

import (
	"context"

	"github.com/samber/lo"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/dao"
	"github.com/hd2yao/go-mall/dal/model"
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
	userCartItemModels := lo.Filter(cartItemModels, func(item *model.ShoppingCartItem, index int) bool {
		return item.UserId == userId
	})
	if len(userCartItemModels) != len(cartItemModels) {
		return nil, errcode.ErrCartWrongUser
	}

	userCartItems := make([]*do.ShoppingCartItem, 0, len(userCartItemModels))
	err = util.CopyProperties(&userCartItems, userCartItemModels)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	// 填充购物项的商品信息
	err = cds.fillInCommodityInfo(userCartItems)
	if err != nil {
		return nil, errcode.Wrap("GetCheckedCartItemsError", err)
	}

	return userCartItems, nil
}

// fillInCommodityInfo 为购物项填充商品信息
func (cds *CartDomainSvc) fillInCommodityInfo(cartItems []*do.ShoppingCartItem) error {
	// 获取购物项中的商品 ID
	commodityIdList := lo.Map(cartItems, func(item *do.ShoppingCartItem, index int) int64 {
		return item.CommodityId
	})

	// 查询商品信息
	commodityDao := dao.NewCommodityDao(cds.ctx)
	commodities, err := commodityDao.FindCommodities(commodityIdList)
	if err != nil {
		return errcode.Wrap("CartItemFillInCommodityInfoError", err)
	}
	if len(commodities) != len(cartItems) {
		logger.New(cds.ctx).Error("fillInCommodityError", "err", "商品信息不匹配", "commodityIdList", commodityIdList,
			"fetchedCommodities", commodities)
		return errcode.ErrCartItemParam
	}

	// 转换成以 ID 为 Key 的商品 Map
	commodityMap := lo.SliceToMap(commodities, func(item *model.Commodity) (int64, *model.Commodity) {
		return item.ID, item
	})
	for _, cartItem := range cartItems {
		cartItem.CommodityName = commodityMap[cartItem.CommodityId].Name
		cartItem.CommodityImg = commodityMap[cartItem.CommodityId].CoverImg
		cartItem.CommoditySellingPrice = commodityMap[cartItem.CommodityId].SellingPrice
	}

	return nil
}
