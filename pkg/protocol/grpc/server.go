package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/aaabhilash97/op/pkg/api/v1"
	"github.com/aaabhilash97/op/pkg/logger"
	"github.com/aaabhilash97/op/pkg/protocol/grpc/middleware"
)

type RunServerOptions struct {
	Ctx       context.Context
	V1API     v1.OpServiceServer
	Port      string
	AccessLog *zap.Logger
}

// RunServer runs gRPC service to publish Op service
func RunServer(option RunServerOptions) error {

	listen, err := net.Listen("tcp", ":"+option.Port)
	if err != nil {
		return err
	}

	// accessLog := logger.InitLogger()
	// gRPC server statup options
	opts := []grpc.ServerOption{}

	// add middleware
	opts = middleware.AddLogging(option.AccessLog, opts)

	// register service
	server := grpc.NewServer(opts...)
	v1.RegisterOpServiceServer(server, option.V1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Warn("shutting down gRPC server...")

			server.GracefulStop()

			<-option.Ctx.Done()
		}
	}()

	// start gRPC server
	logger.Info("starting gRPC server", zap.String("port", option.Port))
	return server.Serve(listen)
}
