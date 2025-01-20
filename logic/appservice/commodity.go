package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
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
		// 出错后 Depress Error, 不触发 500 Server Error
		// 这个错误判断经常用到, 使用一个统一的 ErrorMsg, 方便在生产环境做日志监控 在生产环境上监控系统监控日志中的 ConvertDataError 关键字, 来实现主动告警
		logger.New(cas.ctx).Error(errcode.ErrCoverData.Msg(), "err", err)
		return replyData
	}
	return replyData
}
