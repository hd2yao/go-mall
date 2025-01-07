package app

import (
    "strconv"

    "github.com/gin-gonic/gin"

    "github.com/hd2yao/go-mall/config"
)

type Pagination struct {
    Page      int `json:"page"`
    PageSize  int `json:"page_size"`
    TotalRows int `json:"total_rows"`
}

func NewPagination(c *gin.Context) *Pagination {
    page, _ := strconv.Atoi(c.Query("page"))
    if page < 1 {
        page = 1
    }
    pageSize, _ := strconv.Atoi(c.Query("page_size"))
    if pageSize < 1 {
        pageSize = config.App.Pagination.DefaultSize
    }
    if pageSize > config.App.Pagination.MaxSize {
        pageSize = config.App.Pagination.MaxSize
    }
    return &Pagination{Page: page, PageSize: pageSize}
}
