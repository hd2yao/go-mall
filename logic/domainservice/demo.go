package domainservice

import (
    "context"

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
        demoOrder := new(do.DemoOrder)
        err = util.CopyProperties(demoOrder, demo)
        if err != nil {
            err = errcode.Wrap("copy properties error", err)
            return nil, err
        }
        demoOrders = append(demoOrders, demoOrder)
    }
    return demoOrders, nil
}
