package controller

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/logic/appservice"
)

// AddCartItem 添加购物车
func AddCartItem(c *gin.Context) {
	requestData := new(request.AddCartItem)
	if err := c.ShouldBindJSON(requestData); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	cartAppSvc := appservice.NewCartAppSvc(c)
	err := cartAppSvc.AddCartItem(requestData, c.GetInt64("user_id"))
	if err != nil {
		if errors.Is(err, errcode.ErrCommodityNotExists) {
			app.NewResponse(c).Error(errcode.ErrCommodityNotExists)
		} else if errors.Is(err, errcode.ErrCommodityStockOut) {
			app.NewResponse(c).Error(errcode.ErrCommodityStockOut)
		} else {
			// WithCause 记得加, 不然请求的错误日志里记不到错误原因
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).SuccessOk()
}
