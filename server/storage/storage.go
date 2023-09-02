package storage

import (
	"context"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDecoder interface {
	Decode(interface{}) error
}

// Storage is a struct that holds a sync.Map to store all models.
type storage struct {
	db *mongo.Database
}

type Storage interface {
	GetData(ctx context.Context, userID string) (models.PackedUserData, error)
	SetData(ctx context.Context, data models.UserDataModel) error
	Delete(ctx context.Context, id models.UserDataID) error
}

// NewStorage returns a new Storage.
func NewStorage(db *mongo.Database) *storage {
	return &storage{
		db: db,
	}
}

// SetData sets the model with the given id
// if it exists it will be updated, otherwise it will be created.
func (s *storage) SetData(ctx context.Context, data models.UserDataModel) error {
	if data.GetDataID() == "" {
		data.MakeID()
	}

	// create a new document in mongo
	_, err := s.db.Collection(string(data.GetType())).InsertOne(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			filter := bson.D{
				{"_id", data.GetDataID()},
				{"owner_id", data.GetOrSetOwnerID(nil)},
				{"updated_at", bson.D{{"$lt", data.GetUpdatedAt()}}}}

			_, err = s.db.Collection("goph-keeper").ReplaceOne(ctx, filter, data)

		}
	}
	return err
}

// Delete deletes the model with the given id.
func (s *storage) Delete(ctx context.Context, id models.UserDataID) error {
	_, err := s.db.Collection("goph-keeper").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return err
}

func (s *storage) GetData(ctx context.Context, userID string) (result models.PackedUserData, err error) {
	result = make(models.PackedUserData)
	for _, userDataType := range models.UserDataTypes {
		var res *mongo.Cursor
		res, err = s.db.Collection(string(userDataType)).Find(ctx, bson.D{{"owner_id", userID}})
		if err != nil {
			return
		}
		switch userDataType {
		case models.BinaryType:
			var binary []models.Binary
			err = res.Decode(&binary)
			for _, bin := range binary {
				result[models.BinaryType] = append(result[models.BinaryType], &bin)
			}
		case models.ArbitraryTextType:
			var arbitraryText []models.ArbitraryText
			err = res.Decode(&arbitraryText)
			for _, text := range arbitraryText {
				result[models.ArbitraryTextType] = append(result[models.ArbitraryTextType], &text)
			}
		case models.BankCardType:
			var bankCard []models.BankCard
			err = res.Decode(&bankCard)
			for _, card := range bankCard {
				result[models.BankCardType] = append(result[models.BankCardType], &card)
			}
		case models.LoginType:
			var login []models.Login
			err = res.Decode(&login)
			for _, log := range login {
				result[models.LoginType] = append(result[models.LoginType], &log)
			}
		}
	}
	return result, err
}
