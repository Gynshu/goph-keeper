package storage

import (
	"github.com/gynshu-one/goph-keeper/common/models"
	"sync"
)

// Storage is a struct that holds a sync.Map to store all models.
type Storage interface {
	// Add adds a new model to the storage.
	// Use only For NEW Data
	Add(data models.UserDataModel) error
	Put(data []models.UserDataModel) error
	Get() (data []models.UserDataModel)
}

type storage struct {
	mu *sync.RWMutex
	// repo is a map of models.UserDataModel key is ID field of data
	repo map[string]models.UserDataModel
}

func NewStorage() *storage {
	return &storage{
		mu:   &sync.RWMutex{},
		repo: make(map[string]models.UserDataModel),
	}
}

// Add adds a new model to the storage.
//
//	Use only For NEW Data
func (s *storage) Add(data models.UserDataModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo[data.ID] = data
	return nil
}

func (s *storage) Put(data []models.UserDataModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data == nil {
		return nil
	}
	for _, v := range data {
		s.repo[v.ID] = v
	}
	return nil
}

func (s *storage) Get() (data []models.UserDataModel) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.repo {
		data = append(data, v)
	}
	return
}
