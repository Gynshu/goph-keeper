package storage

import (
	"context"
	"errors"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"github.com/rs/zerolog/log"
	"sync"

	"github.com/gammazero/workerpool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDecoder interface {
	Decode(interface{}) error
}

var workers = workerpool.New(6)

// Storage is a struct that holds a sync.Map to store all models.
type storage struct {
	db    *mongo.Database
	mu    *sync.RWMutex
	cache map[models.UserDataID]models.UserDataModel
}

type Storage interface {
	// Get returns the model with the given id.
	Get(ctx context.Context, userDataType models.UserDataType, ID models.UserDataID) (models.UserDataModel, error)
	GetUserData(ctx context.Context, userID string) ([]models.UserDataModel, error)
	Set(ctx context.Context, userDataType models.UserDataType, data models.UserDataModel) error
	Delete(ctx context.Context, id models.UserDataID) error
}

// NewStorage returns a new Storage.
func NewStorage(db *mongo.Database) *storage {
	return &storage{
		db:    db,
		mu:    &sync.RWMutex{},
		cache: make(map[models.UserDataID]models.UserDataModel),
	}
}

// Get returns the model with the given id.
func (s *storage) Get(ctx context.Context, userDataType models.UserDataType, ID models.UserDataID) (data models.UserDataModel, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.cache[ID]
	if ok {
		return data, nil
	}
	log.Debug().Msgf("Cache miss for %s", ID)

	workers.SubmitWait(func() {
		res := s.db.Collection(string(userDataType)).FindOne(ctx, bson.M{"_id": ID})
		if res.Err() != nil {
			if errors.Is(res.Err(), mongo.ErrNoDocuments) {
				return
			}
		}
		data, err = decode(res, userDataType)
		if err != nil {
			return
		}
	})
	return
}

// Set sets the model with the given id.
func (s *storage) Set(ctx context.Context, userDataType models.UserDataType, data models.UserDataModel) error {
	if data.GetDataID() == "" {
		data.MakeID()
	}

	s.mu.Lock()
	_, ok := s.cache[data.GetDataID()]
	if !ok {
		data.SetCreatedAt()
		data.SetUpdatedAt()
		s.cache[data.GetDataID()] = data
	} else {
		data.SetUpdatedAt()
		s.cache[data.GetDataID()] = data
	}
	s.mu.Unlock()

	var err error
	workers.SubmitWait(func() {
		// create a new document in mongo
		_, err = s.db.Collection(string(userDataType)).InsertOne(ctx, data)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				_, err = s.db.Collection("goph-keeper").ReplaceOne(ctx, bson.M{"_id": data.GetDataID()}, data)
			}
		}
	})
	return err
}

// Delete deletes the model with the given id.
func (s *storage) Delete(ctx context.Context, id models.UserDataID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.cache[id]

	var err error
	if !ok {
		workers.SubmitWait(func() {
			s.db.Collection("goph-keeper").DeleteOne(ctx, bson.M{"_id": id})
			err = ErrObjectMiss
		})
	}
	delete(s.cache, id)
	return err
}

func (s *storage) GetUserData(ctx context.Context, userID string) (result []models.UserDataModel, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workers.SubmitWait(func() {
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
				result = append(result, decoded)

				if res.TryNext(ctx) == false {
					break
				}
			}
		}
	})

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
