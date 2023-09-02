package storage

import (
	"errors"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"sync"
)

// Storage is a struct that holds a sync.Map to store all models.
type Storage interface {
	// Add adds a new model to the storage.
	// Use only For NEW Data
	Add(data models.UserDataModel) error
	Put(data models.PackedUserData) error
	Get() models.PackedUserData
}

type storage struct {
	mu   *sync.RWMutex
	repo models.PackedUserData
}

func NewStorage() *storage {
	return &storage{
		mu:   &sync.RWMutex{},
		repo: make(models.PackedUserData),
	}
}

// Add adds a new model to the storage.
//
//	Use only For NEW Data
func (s *storage) Add(data models.UserDataModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data.GetOrSetOwnerID(&config.CurrentUser.Username)
	data.MakeID()
	data.SetCreatedAt()
	data.SetUpdatedAt()
	s.repo[data.GetType()] = append(s.repo[data.GetType()], data)
	for _, v := range s.repo[data.GetType()] {
		if v.GetDataID() == data.GetDataID() {
			return nil
		}
	}
	return errors.New("failed to add data")
}

func (s *storage) Put(data models.PackedUserData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data == nil {
		return nil
	}
	s.repo = data
	return nil
}

func (s *storage) Get() models.PackedUserData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.repo
}
