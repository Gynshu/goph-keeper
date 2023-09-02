package storage

import (
	"context"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Storage is a struct that holds a sync.Map to store all models.
type storage struct {
	collection *mongo.Collection
}

type Storage interface {
	GetData(ctx context.Context, userID string) ([]models.UserDataModel, error)
	SetData(ctx context.Context, data models.UserDataModel) error
	Delete(ctx context.Context, id string) error
}

// NewStorage returns a new Storage.
func NewStorage(collection *mongo.Collection) *storage {
	return &storage{
		collection: collection,
	}
}

// SetData sets the model with the given id
// if it exists it will be updated, otherwise it will be created.
func (s *storage) SetData(ctx context.Context, data models.UserDataModel) error {
	// create a new document in mongo
	_, err := s.collection.InsertOne(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			filter := bson.D{
				{"_id", data.ID},
				{"owner_id", data.OwnerID},
				{"updated_at", bson.D{{"$lt", data.UpdatedAt}}}}

			_, err = s.collection.ReplaceOne(ctx, filter, data)

		}
	}
	return err
}

// Delete deletes the model with the given id.
func (s *storage) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return err
}

func (s *storage) GetData(ctx context.Context, userID string) (result []models.UserDataModel, err error) {
	res, err := s.collection.Find(ctx, bson.D{{"owner_id", userID}})
	if err != nil {
		return
	}
	defer func(res *mongo.Cursor, ctx context.Context) {
		err = res.Close(ctx)
		if err != nil {
			return
		}
	}(res, ctx)
	err = res.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
