package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuhang-jieke/consulls/srv/getaway/basic/config"
	"github.com/yuhang-jieke/consulls/srv/getaway/client"
	"github.com/yuhang-jieke/consulls/srv/getaway/handler"
	"github.com/yuhang-jieke/consulls/srv/getaway/router"
)

func main() {
	// 加载配置
	cfg := config.DefaultConfig()

	// 设置 gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化 Order 客户端（通过 Consul 发现服务）
	orderClient, err := client.NewOrderClient(cfg.Consul.Address, cfg.Order.ServiceName)
	if err != nil {
		log.Fatalf("Failed to create order client: %v", err)
	}
	log.Printf("Order client initialized, service discovered: %s", cfg.Order.ServiceName)
	// 创建订单处理器
	orderHandler := handler.NewOrderHandler(orderClient)

	// 初始化 Gin 路由
	r := router.Router(orderHandler)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the client with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
