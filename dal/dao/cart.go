package dao

import (
	"context"

	"github.com/hd2yao/go-mall/dal/model"
)

type CartDao struct {
	ctx context.Context
}

func NewCartDao(ctx context.Context) *CartDao {
	return &CartDao{ctx: ctx}
}

// GetUserCartItemWithCommodityId 根据 userId 和 commodityId 查询购物车信息
func (cd *CartDao) GetUserCartItemWithCommodityId(userId, commodityId int64) (*model.ShoppingCartItem, error) {
	cartItem := new(model.ShoppingCartItem)
	err := DB().WithContext(cd.ctx).Where(
		model.ShoppingCartItem{UserId: userId, CommodityId: commodityId},
		"UserId", "CommodityId"). // 保证Struct中的UserId, CommodityId为零值时仍用他们构建查询条件
		Find(&cartItem).Error
	return cartItem, err
}

// UpdateCartItem 更新购物车购物项
func (cd *CartDao) UpdateCartItem(cartItem *model.ShoppingCartItem) error {
	return DBMaster().WithContext(cd.ctx).Model(cartItem).Updates(cartItem).Error
}

// AddCartItem 添加购物车购物项
func (cd *CartDao) AddCartItem(cartItem *model.ShoppingCartItem) error {
	return DBMaster().WithContext(cd.ctx).Create(cartItem).Error
}
