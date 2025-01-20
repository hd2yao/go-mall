package domainservice

import (
	"context"
	"encoding/json"
	"errors"
	"sort"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/common/util"
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

// InitCategoryData 初始化分类信息测试数据
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

// GetHierarchicCategories 返回按层级划分的商品分类
func (cds *CommodityDomainSvc) GetHierarchicCategories() []*do.HierarchicCommodityCategory {
	categoryModels, _ := cds.commodityDao.GetAllCategories()
	FlatCategories := make([]*do.HierarchicCommodityCategory, 0, len(categoryModels))
	err := util.CopyProperties(&FlatCategories, categoryModels)
	if err != nil {
		logger.New(cds.ctx).Error("转换成 HierarchicCommodityCategory 失败")
		return nil
	}

	// 按照 level ASC, rank DESC, id ASC 排序
	sort.SliceStable(FlatCategories, func(i, j int) bool {
		if FlatCategories[i].Level != FlatCategories[j].Level {
			return FlatCategories[i].Level < FlatCategories[j].Level
		}
		if FlatCategories[i].Rank != FlatCategories[j].Rank {
			return FlatCategories[i].Rank > FlatCategories[j].Rank
		}
		return FlatCategories[i].ID < FlatCategories[j].ID
	})

	// 构造一个分类的临时 Map，key 为一二层级的分类 ID，value 为分类信息(包含其子分类)
	categoryTempMap := make(map[int64]*do.HierarchicCommodityCategory)
	for _, category := range FlatCategories {
		if category.ParentId == 0 {
			categoryTempMap[category.ID] = category
		} else if category.ParentId != 0 && category.Level == 2 {
			categoryTempMap[category.ID] = category
			// 把分类添加到父分类的 subCategories 中
			categoryTempMap[category.ParentId].SubCategories = append(categoryTempMap[category.ParentId].SubCategories, category)
		} else if category.ParentId != 0 && category.Level == 3 {
			// 把分类添加到父分类的 subCategories 中
			categoryTempMap[category.ParentId].SubCategories = append(categoryTempMap[category.ParentId].SubCategories, category)
		}
	}

	// 组装按层级划分的商品分类
	var hierarchyCategories []*do.HierarchicCommodityCategory
	for _, category := range FlatCategories {
		if category.ParentId != 0 {
			continue
		}
		category.SubCategories = categoryTempMap[category.ID].SubCategories
		for _, subCategory := range category.SubCategories {
			subCategory.SubCategories = categoryTempMap[subCategory.ID].SubCategories
		}
		hierarchyCategories = append(hierarchyCategories, category)
	}

	return hierarchyCategories
}
