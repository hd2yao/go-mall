package main

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/hd2yao/go-mall/config"
)

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    r.GET("/config-read", func(c *gin.Context) {
        database := config.Database
        c.JSON(http.StatusOK, gin.H{
            "type":     database.Type,
            "max_life": database.MaxLifeTime,
        })
    })

    r.Run()
}
