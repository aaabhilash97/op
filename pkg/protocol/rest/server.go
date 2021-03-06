package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/aaabhilash97/op/pkg/api/v1"
	"github.com/aaabhilash97/op/pkg/logger"
	"github.com/aaabhilash97/op/pkg/protocol/rest/middleware"
)

type RunServerOptions struct {
	Ctx       context.Context
	GrpcPort  string
	HttpPort  string
	AccessLog *zap.Logger
}

// RunServer runs HTTP/REST gateway
func RunServer(option RunServerOptions) error {
	ctx, cancel := context.WithCancel(option.Ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterOpServiceHandlerFromEndpoint(ctx, mux, "localhost:"+option.GrpcPort, opts); err != nil {
		logger.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
	}

	srv := &http.Server{
		Addr: ":" + option.HttpPort,
		// add handler with middleware
		Handler: middleware.AddRequestID(
			middleware.AddLogger(option.AccessLog, mux)),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	logger.Info("starting HTTP/REST gateway", zap.String("port", option.HttpPort))
	return srv.ListenAndServe()
}
