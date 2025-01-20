package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/logic/domainservice"
)

type CommodityAppSvc struct {
	ctx                context.Context
	commodityDomainSvc *domainservice.CommodityDomainSvc
}

func NewCommodityAppSvc(ctx context.Context) *CommodityAppSvc {
	return &CommodityAppSvc{
		ctx:                ctx,
		commodityDomainSvc: domainservice.NewCommodityDomainSvc(ctx),
	}
}

// GetCategoryHierarchy 获取按层级划分的商品分类
func (cas *CommodityAppSvc) GetCategoryHierarchy() []*reply.HierarchicCommodityCategory {
	categories := cas.commodityDomainSvc.GetHierarchicCategories()
	replyData := make([]*reply.HierarchicCommodityCategory, 0, len(categories))
	if len(categories) == 0 {
		return replyData
	}

	err := util.CopyProperties(&replyData, categories)
	if err != nil {
		return replyData
	}
	return replyData
}
