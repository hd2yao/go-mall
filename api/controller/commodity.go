package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/logic/appservice"
)

// GetCategoryHierarchy 获取按层级划分后的所有分类
func GetCategoryHierarchy(c *gin.Context) {
	svc := appservice.NewCommodityAppSvc(c)
	replyData := svc.GetCategoryHierarchy()

	app.NewResponse(c).Success(replyData)
}
