package domainservice

import (
	"context"

	"github.com/samber/lo"

	"github.com/hd2yao/go-mall/api/request"
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

// GetUserCartItems 获取用户购物车里的购物项
func (cds *CartDomainSvc) GetUserCartItems(userId int64) ([]*do.ShoppingCartItem, error) {
	cartItemModels, err := cds.cartDao.GetUserCartItems(userId)

	if err != nil {
		return nil, errcode.Wrap("GetUserCartItemsError", err)
	}

	userCartItems := make([]*do.ShoppingCartItem, 0, len(cartItemModels))
	err = util.CopyProperties(&userCartItems, cartItemModels)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	if len(userCartItems) == 0 {
		return userCartItems, nil
	}

	err = cds.fillInCommodityInfo(userCartItems)
	if err != nil {
		return nil, err
	}
	return userCartItems, nil
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

// CartUpdateItem 更改购物项
func (cds *CartDomainSvc) CartUpdateItem(request *request.CartItemUpdate, userId int64) error {
	// 查询购物项信息
	cartItemModel, err := cds.cartDao.GetCartItemById(request.ItemId)
	if err != nil {
		return errcode.Wrap("CartUpdateItemError", err)
	}

	// 用户不匹配，商品不匹配
	if cartItemModel == nil || cartItemModel.UserId != userId {
		logger.New(cds.ctx).Error("DataMatchError", "cartItem", cartItemModel, "request", request, "requestUserId", userId)
		return errcode.ErrCartWrongUser
	}

	// 更新购物项信息
	cartItemModel.CommodityNum = request.CommodityNum
	err = cds.cartDao.UpdateCartItem(cartItemModel)
	if err != nil {
		return errcode.Wrap("CartUpdateItemError", err)
	}

	return nil
}

// DeleteUserCartItem 删除购物项
func (cds *CartDomainSvc) DeleteUserCartItem(cartItemId, userId int64) error {
	// 查询购物项信息
	cartItemModel, err := cds.cartDao.GetCartItemById(cartItemId)
	if err != nil {
		return errcode.Wrap("CartDeleteItemError", err)
	}

	// 用户不匹配，商品不匹配
	if cartItemModel == nil || cartItemModel.UserId != userId {
		logger.New(cds.ctx).Error("DataMatchError", "cartItem", cartItemModel, "cartItemId", cartItemId, "requestUserId", userId)
		return errcode.ErrCartWrongUser
	}

	err = cds.cartDao.DeleteCartItem(cartItemModel)
	if err != nil {
		return errcode.Wrap("CartDeleteItemError", err)
	}

	return nil
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
