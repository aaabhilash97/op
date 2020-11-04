package config

import (
	"context"
	"log"
	"time"

	"github.com/aaabhilash97/op/pkg/logger"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MongoConfigurations - responsible for mongo configuration
type MongoConfigurations struct {
	Url    string
	DB     string
	Client *mongo.Client
	Cancel context.CancelFunc
	Ctx    context.Context
}

// Close -
func (m *MongoConfigurations) Close() {
	err := m.Client.Disconnect(m.Ctx)
	if err != nil {
		logger.Error("Error in closing mongo connection", zap.Error(err))
	}
	defer m.Cancel()
}

var Mongo MongoConfigurations

// MongoInit -
func init() {
	err := viper.UnmarshalKey("mongo", &Mongo)
	if err != nil {
		log.Fatalf("Unable to decode into mongo config struct, %v", err)
	}
	if Mongo.Url == "" {
		logger.Debug("Mongo config not found")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	mongoOpt := options.Client().ApplyURI(Mongo.Url)

	client, err := mongo.Connect(ctx, mongoOpt)
	if err != nil {
		cancel()
		logger.Fatal("MongoDB connect error", zap.Error(err))
	}
	Mongo.Client = client
	Mongo.Cancel = cancel
	Mongo.Ctx = ctx
}
