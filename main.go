package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/common/middleware"
	"github.com/hd2yao/go-mall/config"
)

func main() {
	g := gin.New()

	// 有了AccessLog 后, 就不需要gin.Logger这个中间件啦
	// g.Use(gin.Logger(), middleware.StartTrace())
	g.Use(middleware.StartTrace(), middleware.LogAccess())

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

	g.POST("/access-log-test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	g.Run()
}
