package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yuhang-jieke/consulls/srv/pkg/consul"
	_ "github.com/yuhang-jieke/consulls/srv/user-server/basic/inits"
	"github.com/yuhang-jieke/consulls/srv/user-server/basic/interceptor"
	__ "github.com/yuhang-jieke/consulls/srv/user-server/handler/proto"
	"github.com/yuhang-jieke/consulls/srv/user-server/handler/server"
	"google.golang.org/grpc"
)

const (
	serviceName    = "user-service"
	serviceAddress = "115.190.57.118"
	consulAddr     = "115.190.57.118:8500"
)

var (
	port = flag.Int("port", 8081, "The server port")
)

func main() {
	flag.Parse()

	client, err := consul.NewClient(consulAddr)
	if err != nil {
		log.Fatalf("Failed to create consul client: %v", err)
	}

	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, serviceAddress, *port)
	tags := []string{"grpc", "v1", "user"}

	err = client.RegisterServiceWithTTL(serviceID, serviceName, serviceAddress, *port, tags, "10s")
	if err != nil {
		log.Fatalf("Failed to register service to consul: %v", err)
	}
	log.Printf("Service registered to consul: %s", serviceID)

	client.PassTTL(serviceID, "service is healthy")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := client.PassTTL(serviceID, "service is healthy"); err != nil {
					log.Printf("Failed to update TTL: %v", err)
				} else {
					log.Println("TTL heartbeat sent")
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingUnaryInterceptor()),
	)
	__.RegisterEcommerceServiceServer(s, &server.Server{})
	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("Failed to serve gRPC: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	cancel()
	s.GracefulStop()
	if err := client.DeregisterService(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}
	log.Println("Service deregistered from consul")
}
