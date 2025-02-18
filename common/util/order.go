package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hd2yao/go-mall/common/enum"
)

// GenOrderNo 生成订单号
func GenOrderNo(userId int64) string {
	day := time.Now().Format(enum.TimeFormatYMD)

	rand.Seed(time.Now().UnixNano())
	seqStr := fmt.Sprintf("%014d", rand.Intn(99999999999999))

	subId := fmt.Sprintf("%04d", userId)
	if len(subId) > 4 {
		subId = subId[len(subId)-5 : len(subId)-1]
	}

	return day + seqStr + subId
}
