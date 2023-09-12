package storage

import (
	"context"
	"github.com/gynshu-one/goph-keeper/common/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// Storage is a struct that holds a sync.Map to store all models.
type storage struct {
	dataCollection *mongo.Collection
	userCollection *mongo.Collection
}

// Storage is an interface for all storage types.
// It provides methods to get and set data.
// Users are stored in a separate collection.
type Storage interface {
	// GetData returns the all data from database associated with the given user id.
	GetData(ctx context.Context, userID string) ([]models.DataWrapper, error)
	// SetData sets the model with the given id
	SetData(ctx context.Context, data models.DataWrapper) error
	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user models.User) error
	// GetUser returns the user with the given email
	GetUser(ctx context.Context, email string) (models.User, error)
}

// NewStorage returns a new Storage.
func NewStorage(dataCollection, userCollection *mongo.Collection) *storage {
	return &storage{
		dataCollection: dataCollection,
		userCollection: userCollection,
	}
}
