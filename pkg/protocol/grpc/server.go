package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	v1 "github.com/aaabhilash97/op/pkg/api/v1"
	"github.com/aaabhilash97/op/pkg/logger"
	"github.com/aaabhilash97/op/pkg/protocol/grpc/middleware"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// accessLog := logger.InitLogger()
	// gRPC server statup options
	opts := []grpc.ServerOption{}

	logger.InitAccessLogger(logger.LogLevels[logger.Config.AccessLogLevel],
		logger.Config.AccessLogOutputPaths)
	// add middleware
	opts = middleware.AddLogging(logger.AccessLog, opts)

	// register service
	server := grpc.NewServer(opts...)
	v1.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Info("starting gRPC server...")
	return server.Serve(listen)
}
