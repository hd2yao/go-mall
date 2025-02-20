package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/router"
	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/config"
)

func main() {
	if config.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.New()
	router.RegisterRoutes(g)
	server := http.Server{
		Addr:    ":8080",
		Handler: g,
	}

	log := logger.New(context.Background())

	// 创建系统信号接收器
	done := make(chan os.Signal)
	// 接收系统信号 os.Interrupt, syscall.SIGINT, syscall.SIGTERM，在收到信号后将信号发送到 done channel
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done // 等待信号
		// 当 done 通道接收到系统信号时，执行 server.Shutdown() 进行优雅关闭
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("ShutdownServerError", "err", err)
		}
	}()

	log.Info("Starting GO MALL HTTP server...")
	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// 服务正常收到关闭信号后 Close
			log.Info("Server closed under request")
		} else {
			// 服务异常关闭
			log.Error("Server closed unexpected", "err", err)
		}
	}
}
