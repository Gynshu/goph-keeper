package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/gynshu-one/goph-keeper/server/pkg/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var workers = workerpool.New(6)

// Storage is a struct that holds a sync.Map to store all models.
type storage struct {
	db    mongo.Database
	mu    *sync.RWMutex
	cache map[string]models.UserData
}

type Storage interface {
	// Get returns the model with the given id.
	Get(ctx context.Context, id string) (models.UserData, error)
	Set(ctx context.Context, id string, data models.UserData) error
	Delete(ctx context.Context, id string) error
}

// NewStorage returns a new Storage.
func NewStorage(db mongo.Database) *storage {
	return &storage{
		db:    db,
		mu:    &sync.RWMutex{},
		cache: make(map[string]models.UserData),
	}
}

// Get returns the model with the given id.
func (s *storage) Get(ctx context.Context, id string) (models.UserData, error) {
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
			res.Decode(&data)
		})
	}

	return data, err
}

// Set sets the model with the given id.
func (s *storage) Set(ctx context.Context, id string, data models.UserData) error {

	s.mu.Lock()
	s.cache[id] = data
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
func (s *storage) Delete(ctx context.Context, id string) error {
	_, ok := s.cache[id]
	if !ok {
		s.db.Collection("goph-keeper").DeleteOne(ctx, bson.M{"_id": id})
		return ErrObjectMiss
	}
	delete(s.cache, id)
	return nil
}
