package controller

import (
	"strconv"

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

// GetCategoriesWithParentId 按parentId查询分类列表
func GetCategoriesWithParentId(c *gin.Context) {
	parentId, _ := strconv.ParseInt(c.Query("parent_id"), 10, 64)
	svc := appservice.NewCommodityAppSvc(c)
	replyData := svc.GetSubCategories(parentId)

	app.NewResponse(c).Success(replyData)
}
