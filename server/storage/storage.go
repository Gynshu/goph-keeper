package storage

import (
	"context"
	"errors"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var workers = workerpool.New(6)

// Storage is a struct that holds a sync.Map to store all models.
type storage struct {
	db    *mongo.Database
	mu    *sync.RWMutex
	cache map[models.UserDataID]models.UserData
}

type Storage interface {
	// Get returns the model with the given id.
	Get(ctx context.Context, id models.UserDataID) (models.UserData, error)
	GetUserData(ctx context.Context, userID string) ([]models.UserData, error)
	Set(ctx context.Context, id models.UserDataID, data models.UserData) error
	Delete(ctx context.Context, id models.UserDataID) error
}

// NewStorage returns a new Storage.
func NewStorage(db *mongo.Database) *storage {
	return &storage{
		db:    db,
		mu:    &sync.RWMutex{},
		cache: make(map[models.UserDataID]models.UserData),
	}
}

// Get returns the model with the given id.
func (s *storage) Get(ctx context.Context, id models.UserDataID) (models.UserData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var err error
	data, ok := s.cache[id]
	if !ok {
		log.Debug().Msgf("Cache miss for %s", id)
		workers.SubmitWait(func() {
			res := s.db.Collection("goph-keeper").FindOne(ctx, bson.M{"_id": id})
			if errors.Is(res.Err(), mongo.ErrNoDocuments) {
				err = ErrObjectMiss
				return
			}
			err = res.Decode(&data)
			if err != nil {
				return
			}
		})
	}

	return data, err
}

// Set sets the model with the given id.
func (s *storage) Set(ctx context.Context, id models.UserDataID, data models.UserData) error {

	s.mu.Lock()
	_, ok := s.cache[id]
	if !ok {
		data.SetCreatedAt()
		data.SetUpdatedAt()
		s.cache[id] = data
	} else {
		data.SetUpdatedAt()
		s.cache[id] = data
	}
	s.mu.Unlock()

	var err error
	workers.SubmitWait(func() {
		// create a new document in mongo
		_, err = s.db.Collection("goph-keeper").InsertOne(ctx, data)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				_, err = s.db.Collection("goph-keeper").ReplaceOne(ctx, bson.M{"_id": id}, data)
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

func (s *storage) GetUserData(ctx context.Context, userID string) ([]models.UserData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var err error
	var data []models.UserData
	workers.SubmitWait(func() {
		res, err := s.db.Collection("goph-keeper").Find(ctx, bson.M{"owner_id": userID})
		if err != nil {
			return
		}
		err = res.All(ctx, &data)
		if err != nil {
			return
		}
	})
	return data, err
}
