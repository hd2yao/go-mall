package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/logic/do"
)

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
