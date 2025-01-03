package main

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/hd2yao/go-mall/common/logger"
    "github.com/hd2yao/go-mall/common/middleware"
    "github.com/hd2yao/go-mall/config"
)

func main() {
    g := gin.Default()

    g.Use(gin.Logger(), middleware.StartTrace())
    g.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    g.GET("/config-read", func(c *gin.Context) {
        database := config.Database

        // 测试 Zap 初始化的临时代码
        logger.ZapLoggerTest(c)

        c.JSON(http.StatusOK, gin.H{
            "type":     database.Type,
            "max_life": database.MaxLifeTime,
        })
    })

    g.GET("/logger-test", func(c *gin.Context) {
        logger.New(c).Info("logger test", "key", "keyName", "val", 2)
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
        })
    })

    g.Run()
}
