package db

import (
	"context"
	"time"

	"github.com/aaabhilash97/op/pkg/config"
	"github.com/aaabhilash97/op/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate go run gen/cmd/gen.go

// Collections _ connected collections from model
var Users *UserModel

// ConnectAll -- connect to mongo tables
func init() {
	if config.Mongo.Client == nil {
		logger.Debug("Mongo client is not initialized")
	}
	clientDB := config.Mongo.Client.Database(config.Mongo.DB)

	collectionName := "users"
	Users = &UserModel{
		collection: clientDB.Collection(collectionName),
	}
}

// MongoTransaction -
type MongoTransaction struct {
	context        context.Context
	Cancel         context.CancelFunc
	Session        mongo.Session
	SessionContext mongo.SessionContext
	err            error
}

// EndSession -
func (t *MongoTransaction) EndSession() {
	defer t.Cancel()
	defer t.Session.EndSession(t.context)
}

// AbortTransaction -
func (t *MongoTransaction) AbortTransaction() error {
	err := t.Session.AbortTransaction(t.SessionContext)
	defer t.EndSession()
	return err
}

// CommitTransaction -
func (t *MongoTransaction) CommitTransaction() error {
	err := t.Session.CommitTransaction(t.SessionContext)
	defer t.EndSession()
	return err
}

// SetError -
func (t *MongoTransaction) SetErr(err error) {
	t.err = err
}

// CommitOrAbort -
func (t *MongoTransaction) CommitOrAbort() error {
	if t.err != nil {
		return t.AbortTransaction()
	}
	return t.CommitTransaction()
}

// NewTransaction -  Start mongo transaction
func NewTransaction(client *mongo.Client) (*MongoTransaction, error) {
	t := &MongoTransaction{}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	t.Cancel = cancel
	t.context = ctx
	session, err := client.StartSession()
	if err != nil {
		cancel()
		return nil, err
	}
	if err = session.StartTransaction(); err != nil {
		cancel()
		session.EndSession(ctx)
		return nil, err
	}
	t.Session = session

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		t.SessionContext = sc
		return nil
	})
	if err != nil {
		cancel()
		session.EndSession(ctx)
		return nil, err
	}

	return t, nil
}
