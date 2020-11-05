package cmd

import (
	"context"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/aaabhilash97/op/pkg/config"
	"github.com/aaabhilash97/op/pkg/db"
	"github.com/aaabhilash97/op/pkg/logger"
	"github.com/aaabhilash97/op/pkg/protocol/grpc"
	"github.com/aaabhilash97/op/pkg/protocol/rest"
	v1 "github.com/aaabhilash97/op/pkg/service/v1"
)

// RunServer runs gRPC server and HTTP gateway
func RunServer(opt config.Config) error {

	// Initialize logger, on any error fn will panic
	logger.Init(opt.Log.LogLevel, opt.Log.LogOutputPaths)

	ctx := context.Background()

	// get configuration

	if len(opt.Server.GrpcPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", opt.Server.GrpcPort)
	}

	if len(opt.Server.HttpPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", opt.Server.HttpPort)
	}

	if err := opt.Mongo.Init(); err != nil {
		return err
	}

	DB := db.Connect(db.Options{
		Client: opt.Mongo.Client,
		DB:     opt.Mongo.DB,
	})

	v1API := v1.NewOpServiceServer(v1.Options{
		Config: opt,
		DB:     DB,
	})

	// initialize access logger
	accessLog, err := logger.CreateLogger(
		logger.LogLevels[opt.Log.AccessLogLevel], opt.Log.AccessLogOutputPaths,
	)
	if err != nil {
		return err
	}

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(rest.RunServerOptions{
			Ctx:       ctx,
			GrpcPort:  opt.Server.GrpcPort,
			HttpPort:  opt.Server.HttpPort,
			AccessLog: accessLog,
		})
	}()

	return grpc.RunServer(grpc.RunServerOptions{
		Ctx:       ctx,
		V1API:     v1API,
		Port:      opt.Server.GrpcPort,
		AccessLog: accessLog,
	})
}
