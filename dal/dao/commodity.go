package dao

import (
	"context"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/model"
	"github.com/hd2yao/go-mall/logic/do"
)

type CommodityDao struct {
	ctx context.Context
}

func NewCommodityDao(ctx context.Context) *CommodityDao {
	return &CommodityDao{ctx: ctx}
}

// InitCategoryData 初始化商品分类数据
func (cd *CommodityDao) InitCategoryData(categoryDos []*do.CommodityCategory) error {
	categoryModels := make([]*model.CommodityCategory, 0, len(categoryDos))
	err := util.CopyProperties(&categoryModels, categoryDos)
	if err != nil {
		return errcode.ErrCoverData
	}
	return cd.BulkCreateCommodityCategories(categoryModels)
}

// BulkCreateCommodityCategories 批量创建商品分类
func (cd *CommodityDao) BulkCreateCommodityCategories(categories []*model.CommodityCategory) error {
	return DBMaster().WithContext(cd.ctx).Create(categories).Error
}

// GetAllCategories 获取所有商品分类
func (cd *CommodityDao) GetAllCategories() ([]*model.CommodityCategory, error) {
	categories := make([]*model.CommodityCategory, 0)
	err := DB().WithContext(cd.ctx).Find(&categories).Error
	return categories, err
}

// GetSubCategories 查询指定 ID 下的商品分类
func (cd *CommodityDao) GetSubCategories(parentId int64) ([]*model.CommodityCategory, error) {
	categories := make([]*model.CommodityCategory, 0)
	err := DB().WithContext(cd.ctx).
		Where("parent_id = ?", parentId).
		Order("rank DESC").
		Find(&categories).Error
	return categories, err
}

// InitCommodityData 初始化商品数据
func (cd *CommodityDao) InitCommodityData(commodityDos []*do.Commodity) error {
	commodityModels := make([]*model.Commodity, 0, len(commodityDos))
	err := util.CopyProperties(&commodityModels, commodityDos)
	if err != nil {
		return errcode.ErrCoverData
	}
	return cd.BulkCreateCommodities(commodityModels)
}

// BulkCreateCommodities 批量创建商品
func (cd *CommodityDao) BulkCreateCommodities(commodities []*model.Commodity) error {
	return DBMaster().WithContext(cd.ctx).Create(commodities).Error
}

// GetOneCommodity 无查询条件，返回一条数据
func (cd *CommodityDao) GetOneCommodity() (*model.Commodity, error) {
	commodity := new(model.Commodity)
	err := DB().WithContext(cd.ctx).Find(commodity).Error
	return commodity, err
}

// GetCategoryById 获取Id对应的分类信息
func (cd *CommodityDao) GetCategoryById(categoryId int64) (*model.CommodityCategory, error) {
	category := new(model.CommodityCategory)
	err := DB().WithContext(cd.ctx).Where("id = ?", categoryId).Find(category).Error
	return category, err
}

// GetThirdLevelCategories 获取指定分类下的所有三级分类 ID
func (cd *CommodityDao) GetThirdLevelCategories(categoryInfo *do.CommodityCategory) (categoryIds []int64, err error) {
	if categoryInfo.Level == 3 {
		return []int64{categoryInfo.ID}, nil
	} else if categoryInfo.Level == 2 {
		categoryIds, err = cd.getSubCategoryIdList([]int64{categoryInfo.ID})
		return
	} else if categoryInfo.Level == 1 {
		var secondCategoryId []int64
		secondCategoryId, err = cd.getSubCategoryIdList([]int64{categoryInfo.ID})
		if err != nil {
			return
		}
		categoryIds, err = cd.getSubCategoryIdList(secondCategoryId)
		return
	}
	return
}

// getSubCategoryIdList 查询分类的子分类ID
func (cd *CommodityDao) getSubCategoryIdList(parentCategoryIds []int64) (categoryIds []int64, err error) {
	err = DB().WithContext(cd.ctx).Model(&model.CommodityCategory{}).
		Where("parent_id IN (?)", parentCategoryIds).
		Order("rank DESC").Pluck("id", &categoryIds).Error
	return
}

// GetCommoditiesInCategory 查询分类下的商品列表
func (cd *CommodityDao) GetCommoditiesInCategory(categoryIds []int64, offset, returnSize int) (commodityList []*model.Commodity, totalRows int64, err error) {
	// 查询满足条件的商品
	err = DB().WithContext(cd.ctx).Omit("detail_content"). // 忽略商品详情 detail_content 字段
								Where("category_id IN (?)", categoryIds).
								Offset(offset).Limit(returnSize).
								Find(&commodityList).Error
	// 查询满足条件的商品总数
	DB().WithContext(cd.ctx).Model(model.Commodity{}).
		Where("category_id IN (?)", categoryIds).Count(&totalRows)
	return
}

// FindCommodityWithNameKeyword 按名称LIKE查询商品列表
func (cd *CommodityDao) FindCommodityWithNameKeyword(keyword string, offset, returnSize int) (commodityList []*model.Commodity, totalRows int64, err error) {
	err = DB().WithContext(cd.ctx).Omit("detail_content").
		Where("name LIKE ?", "%"+keyword+"%").
		Offset(offset).Limit(returnSize).
		Find(&commodityList).Error
	DB().WithContext(cd.ctx).Model(model.Commodity{}).Where("name LIKE ?", "%"+keyword+"%").Count(&totalRows)
	return
}

// FindCommodityById 通过ID查商品信息
func (cd *CommodityDao) FindCommodityById(commodityId int64) (*model.Commodity, error) {
	commodity := new(model.Commodity)
	err := DB().WithContext(cd.ctx).Omit("detail_content").
		Where("id = ?", commodityId).Find(commodity).Error
	return commodity, err
}

// FindCommodities 查询主键 id IN commodityIdList 的 商品
func (cd *CommodityDao) FindCommodities(commodityIdList []int64) ([]*model.Commodity, error) {
	commodities := make([]*model.Commodity, 0)
	err := DB().WithContext(cd.ctx).Find(&commodities, commodityIdList).Error
	return commodities, err
}
