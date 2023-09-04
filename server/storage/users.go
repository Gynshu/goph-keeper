package storage

import (
	"context"
	"github.com/gynshu-one/goph-keeper/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *storage) CreateUser(ctx context.Context, user models.User) error {
	// Try to create a new user
	_, err := s.userCollection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return err
		}
		return err
	}
	return nil
}

func (s *storage) GetUser(ctx context.Context, email string) (models.User, error) {
	var user models.User
	filter := bson.D{{"_id", email}}
	err := s.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
