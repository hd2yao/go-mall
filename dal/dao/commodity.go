package dao

import (
	"context"

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
		return err
	}
	return cd.BulkCreateCommodityCategories(categoryModels)
}

// GetAllCategories 获取所有商品分类
func (cd *CommodityDao) GetAllCategories() ([]*model.CommodityCategory, error) {
	categories := make([]*model.CommodityCategory, 0)
	err := DB().WithContext(cd.ctx).Find(&categories).Error
	return categories, err
}

// BulkCreateCommodityCategories 批量创建商品分类
func (cd *CommodityDao) BulkCreateCommodityCategories(categories []*model.CommodityCategory) error {
	return DBMaster().WithContext(cd.ctx).Create(categories).Error
}
