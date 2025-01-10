package library

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/common/util/httptool"
)

type DemoLib struct {
	ctx context.Context
}

// NewDemoLib 创建时上层通过ctx 把 gin.Ctx传递过来
func NewDemoLib(ctx context.Context) *DemoLib {
	return &DemoLib{ctx: ctx}
}

type OrderCreateResult struct {
	UserId    int64  `json:"userId"`
	BillMoney int64  `json:"billMoney"`
	OrderNo   string `json:"orderNo"`
	State     int8   `json:"state"`
	PaidAt    string `json:"paidAt"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// 用http调自己项目里的POST接口, 演示用, 实际使用时不要这么干

func (lib *DemoLib) TestPostCreateOrder() (*OrderCreateResult, error) {
	data := &request.DemoOrderCreate{
		UserId:       12345,
		BillMoney:    20,
		OrderGoodsId: 1111110,
	}
	jsonReq, _ := json.Marshal(data)
	httCode, respBody, err := httptool.Post(lib.ctx, "http://localhost:8080/building/demo-order-create", jsonReq)
	logger.New(lib.ctx).Info("create-demo-order api response ", "code", httCode, "data", string(respBody), "err", err)

	if err != nil {
		return nil, err
	}

	if !json.Valid([]byte(respBody)) {
		fmt.Println("Invalid JSON")
	}

	reply := &struct {
		Code int `json:"code"`
		Data *OrderCreateResult
	}{}
	json.Unmarshal(respBody, reply)
	fmt.Printf("%+v\n", reply.Data)
	return reply.Data, nil
}
