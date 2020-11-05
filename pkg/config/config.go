package config

import (
	"context"
	"fmt"
	"time"

	"github.com/aaabhilash97/op/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type ServerConfigurations struct {
	Host     string
	GrpcPort string
	HttpPort string
}

// MongoConfigurations - responsible for mongo configuration
type MongoConfigurations struct {
	Url    string
	DB     string
	Client *mongo.Client
	Cancel context.CancelFunc
	Ctx    context.Context
}

func (m *MongoConfigurations) Init() error {
	if m.Url == "" {
		logger.Error("Mongo config not found")
		return fmt.Errorf("Mongo url missing")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	mongoOpt := options.Client().ApplyURI(m.Url)

	client, err := mongo.Connect(ctx, mongoOpt)
	if err != nil {
		cancel()
		logger.Error("MongoDB connect error", zap.Error(err))
		return err
	}
	m.Client = client
	m.Cancel = cancel
	m.Ctx = ctx
	return nil
}

// Close -
func (m *MongoConfigurations) Close() {
	err := m.Client.Disconnect(m.Ctx)
	if err != nil {
		logger.Error("Error in closing mongo connection", zap.Error(err))
	}
	defer m.Cancel()
}

type LogConfigurations struct {
	LogLevel             string
	LogOutputPaths       string
	AccessLogLevel       string
	AccessLogOutputPaths string
}

type Config struct {
	Server ServerConfigurations
	Mongo  MongoConfigurations
	Log    LogConfigurations
}
