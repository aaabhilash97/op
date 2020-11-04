package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserModel is the model for base table
type UserModel struct {
	collection *mongo.Collection

	ID     *primitive.ObjectID `bson:"_id,omitempty"`
	UserID *primitive.ObjectID `bson:"user_id,omitempty"`

	// fields
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
	CreatedAt *time.Time `bson:"created_at,omitempty"`
}
