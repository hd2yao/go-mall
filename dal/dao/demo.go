package dao

import (
    "context"

    "github.com/hd2yao/go-mall/common/util"
    "github.com/hd2yao/go-mall/dal/model"
    "github.com/hd2yao/go-mall/logic/do"
)

type DemoDao struct {
    ctx context.Context
}

func NewDemoDao(ctx context.Context) *DemoDao {
    return &DemoDao{ctx: ctx}
}

func (demo *DemoDao) GetAllDemos() (demos []*model.DemoOrder, err error) {
    err = DB().WithContext(demo.ctx).Find(&demos).Error
    if err != nil {
        return nil, err
    }
    return demos, err
}

// CreateDemoOrder 创建订单
func (demo *DemoDao) CreateDemoOrder(demoOrder *do.DemoOrder) (*model.DemoOrder, error) {
    modelData := new(model.DemoOrder)
    err := util.CopyProperties(modelData, demoOrder)
    if err != nil {
        return nil, err
    }
    err = DB().WithContext(demo.ctx).Create(modelData).Error
    return modelData, err
}
