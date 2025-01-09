package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/logic/do"
	"github.com/hd2yao/go-mall/logic/domainservice"
)

// 演示 DEMO，后期使用时删掉

type DemoAppSvc struct {
	ctx           context.Context
	demoDomainSvc *domainservice.DemoDomainSvc
}

func NewDemoAppSvc(ctx context.Context) *DemoAppSvc {
	return &DemoAppSvc{
		ctx:           ctx,
		demoDomainSvc: domainservice.NewDemoDomainSvc(ctx),
	}
}

func (das *DemoAppSvc) GetDemoIdentities() ([]int64, error) {
	demos, err := das.demoDomainSvc.GetDemos()
	if err != nil {
		return nil, err
	}
	identities := make([]int64, 0, len(demos))
	for _, demo := range demos {
		identities = append(identities, demo.Id)
	}
	return identities, nil
}

func (das *DemoAppSvc) CreateDemoOrders(orderRequest *request.DemoOrderCreate) (*reply.DemoOrder, error) {
	demoOrderDo := new(do.DemoOrder)
	err := util.CopyProperties(demoOrderDo, orderRequest)
	if err != nil {
		err = errcode.Wrap("请求转换成 demoOrderDo 失败", err)
		return nil, err
	}
	demoOrder, err := das.demoDomainSvc.CreateDemoOrder(demoOrderDo)
	if err != nil {
		return nil, err
	}

	// 做一些其他的创建订单成功后的外围逻辑
	// 比如异步发送创建订单创建通知

	replyDemoOrder := new(reply.DemoOrder)
	err = util.CopyProperties(replyDemoOrder, demoOrder)
	if err != nil {
		err = errcode.Wrap("demoOrder 转换成 replyDemoOrder 失败", err)
		return nil, err
	}

	return replyDemoOrder, err
}
