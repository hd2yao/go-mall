package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/logic/do"
)

type DummyDemoOrder struct {
	OrderNo string `redis:"orderNo"`
	UserId  int64  `redis:"userId"`
}

// SetDemoOrderStruct 使用 HSET 的存储结构体数据
func SetDemoOrderStruct(ctx context.Context, demoOrder *do.DemoOrder) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	data := struct {
		OrderNo string `redis:"orderNo"`
		UserId  int64  `redis:"userId"`
	}{
		OrderNo: demoOrder.OrderNo,
		UserId:  demoOrder.UserId,
	}
	_, err := Redis().HSet(ctx, redisKey, data).Result()
	if err != nil {
		logger.New(ctx).Error("redis error", "err", err)
		return err
	}
	return nil
}

// GetDemoOrderStruct 使用 HGETALL 和 Scan 读取结构体数据
func GetDemoOrderStruct(ctx context.Context, orderNo string) (*DummyDemoOrder, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, orderNo)
	data := new(DummyDemoOrder)
	err := Redis().HGetAll(ctx, redisKey).Scan(&data)
	Redis().Get(ctx, redisKey).String()
	if err != nil {
		logger.New(ctx).Error("redis error", "err", err)
		return nil, err
	}
	logger.New(ctx).Info("scan data from redis", "data", &data)
	return data, nil
}

// SetDemoOrder 缓存设置
func SetDemoOrder(ctx context.Context, demoOrder *do.DemoOrder) error {
	jsonDataBytes, _ := json.Marshal(demoOrder)
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	_, err := Redis().Set(ctx, redisKey, jsonDataBytes, 0).Result()
	if err != nil {
		logger.New(ctx).Error("redis error", "err", err)
		return err
	}
	return nil
}

// GetDemoOrder 缓存获取
func GetDemoOrder(ctx context.Context, orderNo string) (*do.DemoOrder, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, orderNo)
	jsonBytes, err := Redis().Get(ctx, redisKey).Bytes()
	if err != nil {
		logger.New(ctx).Error("redis error", "err", err)
		return nil, err
	}
	data := new(do.DemoOrder)
	json.Unmarshal(jsonBytes, &data)
	return data, nil
}
