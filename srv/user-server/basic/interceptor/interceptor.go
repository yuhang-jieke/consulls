package interceptor

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

var (
	logger *log.Logger
)

func init() {
	file, err := os.OpenFile("grpc_server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		logger.Printf("[入参] method=%s, request=%+v", info.FullMethod, req)

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		if err != nil {
			logger.Printf("[响应错误] method=%s, duration=%v, error=%v", info.FullMethod, duration, err)
		} else {
			logger.Printf("[响应成功] method=%s, duration=%v, response=%+v", info.FullMethod, duration, resp)
		}

		return resp, err
	}
}

func FormatLog(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
