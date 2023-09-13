package storage

import (
	"context"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetData sets the model with the given id
// if {it exists and update date of new data is newer it will be updated}, {otherwise it will be created.}
func (s *storage) SetData(ctx context.Context, data models.DataWrapper) error {
	// Create a new document in mongo
	_, err := s.dataCollection.InsertOne(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// We have it already
			// So if the owner id and update time are valid, update it
			filter := bson.D{
				{"_id", data.ID},
				{"owner_id", data.OwnerID},
				{"updated_at", bson.D{{"$lt", data.UpdatedAt}}}}

			update := bson.D{
				{"$set", bson.D{
					{"data", data.Data},
					{"name", data.Name},
					{"updated_at", data.UpdatedAt},
					{"deleted_at", data.DeletedAt},
				}},
			}
			_, err = s.dataCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// GetData returns the all data from database associated with the given user id.
func (s *storage) GetData(ctx context.Context, userID string) (result []models.DataWrapper, err error) {
	res, err := s.dataCollection.Find(ctx, bson.D{{"owner_id", userID}})
	if err != nil {
		return
	}
	defer func(res *mongo.Cursor, ctx context.Context) {
		err = res.Close(ctx)
		if err != nil {
			log.Err(err).Msg("failed to close cursor")
			return
		}
	}(res, ctx)
	err = res.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
