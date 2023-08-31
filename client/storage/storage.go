package storage

import (
	"github.com/gynshu-one/goph-keeper/shared/models"
	"sync"
)

// Storage is a struct that holds a sync.Map to store all models.
type Storage interface {
	PutData(data models.PackedUserData) error
	GetData() models.PackedUserData
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

func (s *storage) PutData(data models.PackedUserData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo = data
	return nil
}

func (s *storage) GetData() models.PackedUserData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := make(models.PackedUserData)
	for k, v := range s.repo {
		cp[k] = v
	}
	return cp
}
