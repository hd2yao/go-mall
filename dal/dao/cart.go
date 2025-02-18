package dao

import (
	"context"

	"gorm.io/gorm"

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
		// 使用 Struct 作为 Where 的参数时 最好指定要搜索的字段，否则字段值为零值时不会用来构建查询条件
		// 文档 https://gorm.io/docs/query.html#Specify-Struct-search-fields
		Find(&cartItem).Error
	return cartItem, err
}

// GetUserCartItems 获取用户购物车里的购物项
func (cd *CartDao) GetUserCartItems(userId int64) ([]*model.ShoppingCartItem, error) {
	cartItems := make([]*model.ShoppingCartItem, 0)
	err := DB().WithContext(cd.ctx).Where(model.ShoppingCartItem{UserId: userId}, "UserId").Find(&cartItems).Error
	return cartItems, err
}

// GetCartItemById 根据购物项 ID 获取信息
func (cd *CartDao) GetCartItemById(cartItemId int64) (*model.ShoppingCartItem, error) {
	cartItem := new(model.ShoppingCartItem)
	err := DB().WithContext(cd.ctx).Where(model.ShoppingCartItem{CartItemId: cartItemId}, "CartItemId").Find(&cartItem).Error
	return cartItem, err
}

// FindCartItems 获取多个ID指定的购物项
func (cd *CartDao) FindCartItems(cartItemIdList []int64) ([]*model.ShoppingCartItem, error) {
	items := make([]*model.ShoppingCartItem, 0)
	// 查询主键 id IN cartItemIdList 的购物项
	err := DB().WithContext(cd.ctx).Find(&items, cartItemIdList).Error
	return items, err
}

// UpdateCartItem 更新购物车购物项
func (cd *CartDao) UpdateCartItem(cartItem *model.ShoppingCartItem) error {
	return DBMaster().WithContext(cd.ctx).Model(cartItem).Updates(cartItem).Error
}

// AddCartItem 添加购物车购物项
func (cd *CartDao) AddCartItem(cartItem *model.ShoppingCartItem) error {
	return DBMaster().WithContext(cd.ctx).Create(cartItem).Error
}

// DeleteCartItem 删除购物车购物项
func (cd *CartDao) DeleteCartItem(cartItem *model.ShoppingCartItem) error {
	return DBMaster().WithContext(cd.ctx).Delete(cartItem).Error
}

// DeleteMultiCartItemInTx 创建订单后调用该方法删除购物项
func (cd *CartDao) DeleteMultiCartItemInTx(tx *gorm.DB, cartIdList []int64) error {
	return tx.WithContext(cd.ctx).Delete(&model.ShoppingCartItem{}, cartIdList).Error
}
