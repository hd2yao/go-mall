package appservice

import (
	"context"

	"github.com/hd2yao/go-mall/logic/domainservice"
)

type OrderAppSvc struct {
	ctx            context.Context
	orderDomainSvc *domainservice.OrderDomainSvc
}
