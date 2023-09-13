package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gynshu-one/goph-keeper/client/auth"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rs/zerolog/log"
)

// Storage is a struct that holds a sync.Map to store all models.
type Storage interface {
	// AddEncrypt adds a new model to the storage.
	// Use only For NEW Data
	// it encrypts data and saves it to the storage
	// by creating models.DataWrapper struct and adding it to the storage
	// Wrapper should be passed with Name and Type fields
	AddEncrypt(data models.BasicData, wrapper models.DataWrapper) error
	// Swap replaces all data in the storage with new data
	// this is server client exchange method
	Swap(data []models.DataWrapper) error
	// FindDecrypt finds a model in the storage by id and decrypts it
	// returns decrypted data and wrapper
	// if wrapper content (data) is deleted returns error and wrapper
	// for ui
	FindDecrypt(id string) (data any, wrapper models.DataWrapper, err error)
	// Delete sets deleted time and clears data field of
	Delete(id string) error

	// Get returns all data from storage
	// for server
	Get() (data []models.DataWrapper)
}

type storage struct {
	mu *sync.RWMutex
	// repo is a map of models.DataWrapper key is ID field of data
	repo map[string]models.DataWrapper
}

// NewStorage creates a new storage instance
func NewStorage() Storage {
	return &storage{
		mu:   &sync.RWMutex{},
		repo: make(map[string]models.DataWrapper),
	}
}

// AddEncrypt adds a new model to the storage.
// Use only For NEW Data
// it encrypts data and saves it to the storage
// by creating models.DataWrapper struct and adding it to the storage
// Wrapper should be passed with Name and Type fields
func (s *storage) AddEncrypt(data models.BasicData, wrapper models.DataWrapper) error {
	secret := auth.GetSecret()
	if secret == "" {
		log.Fatal().Msg("secret is nil")
	}
	encrypted, err := data.EncryptAll(secret)
	if err != nil {
		return err
	}
	t := time.Now().Unix()
	wrapper.Data = encrypted
	if wrapper.CreatedAt == 0 {
		wrapper.CreatedAt = t
	}
	wrapper.UpdatedAt = t
	if wrapper.ID == "" {
		wrapper.ID = uuid.NewString()
	}
	wrapper.OwnerID = auth.CurrentUser.Username
	wrapper.DeletedAt = 0
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo[wrapper.ID] = wrapper
	return nil
}

// FindDecrypt finds a model in the storage by id and decrypts it
// returns decrypted data and wrapper
// if wrapper content (data) is deleted returns error and wrapper
func (s *storage) FindDecrypt(id string) (data any, wrapper models.DataWrapper, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Find wrapper
	wrapper, ok := s.repo[id]
	if !ok {
		return nil, wrapper, models.ErrDeleted
	}

	// Decrypt data with os keyring secret
	secret := auth.GetSecret()
	if secret == "" {
		log.Fatal().Msg("secret is nil")
	}

	if wrapper.DeletedAt > 0 {
		return nil, wrapper, models.ErrDeleted
	}

	// Determine type and decrypt
	switch wrapper.Type {
	case models.LoginType:
		var login models.Login
		if err = login.DecryptAll(secret, wrapper.Data); err != nil {
			return nil, wrapper, err
		}
		return login, wrapper, nil
	case models.ArbitraryTextType:
		var text models.ArbitraryText
		if err = text.DecryptAll(secret, wrapper.Data); err != nil {
			return nil, wrapper, err
		}
		return text, wrapper, nil
	case models.BankCardType:
		var bankCard models.BankCard
		if err = bankCard.DecryptAll(secret, wrapper.Data); err != nil {
			return nil, wrapper, err
		}
		return bankCard, wrapper, nil
	case models.BinaryType:
		var binary models.Binary
		if err = binary.DecryptAll(secret, wrapper.Data); err != nil {
			return nil, wrapper, err
		}
		return binary, wrapper, nil
	}

	return nil, wrapper, models.ErrUnknownType
}

// Delete sets deleted time and clears data field of
func (s *storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find wrapper
	item, ok := s.repo[id]
	if !ok {
		return fmt.Errorf("item with id %s not found", id)
	}

	// Set basic fields
	item.DeletedAt = time.Now().Unix()
	item.UpdatedAt = time.Now().Unix()
	item.Data = nil
	s.repo[id] = item
	return nil
}

// Swap replaces all data in the storage with new data
// this is server client exchange method
func (s *storage) Swap(data []models.DataWrapper) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data == nil {
		return nil
	}

	// Completely replace the storage
	s.repo = make(map[string]models.DataWrapper)
	for i := range data {
		s.repo[data[i].ID] = data[i]
	}
	return nil
}

// Get returns all data from storage
func (s *storage) Get() (data []models.DataWrapper) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Copy all data from the storage to the slice
	for v := range s.repo {
		data = append(data, s.repo[v])
	}
	return
}
