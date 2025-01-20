package domainservice

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/dal/dao"
	"github.com/hd2yao/go-mall/logic/do"
	"github.com/hd2yao/go-mall/resources"
)

type CommodityDomainSvc struct {
	ctx          context.Context
	commodityDao *dao.CommodityDao
}

func NewCommodityDomainSvc(ctx context.Context) *CommodityDomainSvc {
	return &CommodityDomainSvc{
		ctx:          ctx,
		commodityDao: dao.NewCommodityDao(ctx),
	}
}

func (cds *CommodityDomainSvc) InitCategoryData() error {
	categories, err := cds.commodityDao.GetAllCategories()
	if err != nil {
		return errcode.Wrap("初始化商品分类错误", err)
	}
	if len(categories) > 1 {
		// 避免重复初始化
		return errcode.Wrap("重复初始化商品分类", errors.New("不能重复初始化商品分类"))
	}

	cateInitFileHandler, _ := resources.LoadResourceFile("category_init_data.json")

	categoryDos := make([]*do.CommodityCategory, 0, len(categories))
	decoder := json.NewDecoder(cateInitFileHandler)
	decoder.Decode(&categoryDos)

	err = cds.commodityDao.InitCategoryData(categoryDos)
	if err != nil {
		return errcode.Wrap("初始化商品分类错误", err)
	}

	return nil
}
