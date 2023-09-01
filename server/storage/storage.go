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
	GetData(ctx context.Context, userID string) (map[models.UserDataID]models.UserDataModel, error)
	SetData(ctx context.Context, userDataType models.UserDataType, data models.UserDataModel) error
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
func (s *storage) SetData(ctx context.Context, userDataType models.UserDataType, data models.UserDataModel) error {
	if data.GetDataID() == "" {
		data.MakeID()
	}
	// create a new document in mongo
	_, err := s.db.Collection(string(userDataType)).InsertOne(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			_, err = s.db.Collection("goph-keeper").ReplaceOne(ctx, bson.M{"_id": data.GetDataID()}, data)
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

func (s *storage) GetData(ctx context.Context, userID string) (result map[models.UserDataID]models.UserDataModel, err error) {
	for _, userDataType := range models.UserDataTypes {
		var res *mongo.Cursor
		res, err = s.db.Collection(string(userDataType)).Find(ctx, bson.D{{"owner_id", userID}})
		if err != nil {
			return
		}
		for res.Next(ctx) {
			var decoded models.UserDataModel
			decoded, err = decode(res, userDataType)
			if err != nil {
				return
			}
			result[decoded.GetDataID()] = decoded

			if res.TryNext(ctx) == false {
				break
			}
		}
	}
	return result, err
}

func decode(decoder mongoDecoder, userDataType models.UserDataType) (data models.UserDataModel, err error) {
	switch userDataType {
	case models.BinaryType:
		var binary models.Binary
		err = decoder.Decode(&binary)
		return &binary, err
	case models.ArbitraryTextType:
		var arbitraryText models.ArbitraryText
		err = decoder.Decode(&arbitraryText)
		return &arbitraryText, err
	case models.BankCardType:
		var bankCard models.BankCard
		err = decoder.Decode(&bankCard)
		return &bankCard, err
	case models.LoginType:
		var login models.Login
		err = decoder.Decode(&login)
		return &login, err
	default:
		err = ErrObjectMiss
	}
	return nil, err
}
