package db

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// BaseModel is the model for base table
type BaseModel struct {
	collection *mongo.Collection
	// fields
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
	CreatedAt *time.Time `bson:"created_at,omitempty"`
}
