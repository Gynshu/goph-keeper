package storage

import (
	"github.com/gynshu-one/goph-keeper/shared/models"
	"sync"
)

// Storage is a struct that holds a sync.Map to store all models.
type Storage interface {
	Add(data models.UserDataModel)
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

func (s *storage) Add(data models.UserDataModel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo[data.GetType()] = append(s.repo[data.GetType()], data)
}

func (s *storage) Put(data models.PackedUserData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo = data
	return nil
}

func (s *storage) Get() models.PackedUserData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := make(models.PackedUserData)
	for k, v := range s.repo {
		cp[k] = v
	}
	return cp
}
