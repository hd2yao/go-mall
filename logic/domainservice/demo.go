package domainservice

import (
    "context"
    "fmt"

    "github.com/hd2yao/go-mall/common/errcode"
    "github.com/hd2yao/go-mall/common/util"
    "github.com/hd2yao/go-mall/dal/dao"
    "github.com/hd2yao/go-mall/logic/do"
)

// 演示 DEMO，后期使用时删掉

type DemoDomainSvc struct {
    ctx     context.Context
    DemoDao *dao.DemoDao
}

func NewDemoDomainSvc(ctx context.Context) *DemoDomainSvc {
    return &DemoDomainSvc{
        ctx:     ctx,
        DemoDao: dao.NewDemoDao(ctx),
    }
}

func (dds *DemoDomainSvc) GetDemos() ([]*do.DemoOrder, error) {
    demos, err := dds.DemoDao.GetAllDemos()
    if err != nil {
        err = errcode.Wrap("query entity error", err)
        return nil, err
    }

    demoOrders := make([]*do.DemoOrder, 0, len(demos))
    for _, demo := range demos {
        fmt.Println(demo)
        demoOrder := new(do.DemoOrder)
        err = util.CopyProperties(demoOrder, demo)
        fmt.Println(demoOrder)
        if err != nil {
            err = errcode.Wrap("copy properties error", err)
            return nil, err
        }
        demoOrders = append(demoOrders, demoOrder)
    }
    return demoOrders, nil
}

func (dds *DemoDomainSvc) CreateDemoOrder(demoOrder *do.DemoOrder) (*do.DemoOrder, error) {
    // 生成订单号
    demoOrder.OrderNo = "20240627596615375920904456"
    demoOrderModel, err := dds.DemoDao.CreateDemoOrder(demoOrder)
    if err != nil {
        err = errcode.Wrap("create demo order error", err)
        return nil, err
    }
    // TODO: 写订单快照
    // 这里一般要在事务里写订单商品快照表
    err = util.CopyProperties(demoOrder, demoOrderModel)
    // 返回 domain 对象
    return demoOrder, err
}
