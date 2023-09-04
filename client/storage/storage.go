package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gynshu-one/goph-keeper/client/auth"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

// Storage is a struct that holds a sync.Map to store all models.
type Storage interface {
	// AddEncrypt adds a new model to the storage.
	// Use only For NEW Data
	AddEncrypt(data models.BasicData, wrapper models.DataWrapper) error
	Swap(data []models.DataWrapper) error
	FindDecrypt(id string) (data any, wrapper models.DataWrapper, err error)
	Delete(id string) error
	Get() (data []models.DataWrapper)
}

type storage struct {
	mu *sync.RWMutex
	// repo is a map of models.DataWrapper key is ID field of data
	repo map[string]models.DataWrapper
}

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
	wrapper.OwnerID = config.CurrentUser.Username
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
	wrapper, ok := s.repo[id]
	if !ok {
		return nil, wrapper, models.ErrDeleted
	}
	secret := auth.GetSecret()
	if secret == "" {
		log.Fatal().Msg("secret is nil")
	}

	if wrapper.DeletedAt > 0 {
		return nil, wrapper, models.ErrDeleted
	}

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

	return nil, wrapper, models.UnknownType
}
func (s *storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.repo[id]
	if !ok {
		return fmt.Errorf("item with id %s not found", id)
	}
	item.DeletedAt = time.Now().Unix()
	item.UpdatedAt = time.Now().Unix()
	item.Data = nil
	s.repo[id] = item
	return nil
}
func (s *storage) Swap(data []models.DataWrapper) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data == nil {
		return nil
	}
	// full clear
	s.repo = make(map[string]models.DataWrapper)
	for i, _ := range data {
		s.repo[data[i].ID] = data[i]
	}
	return nil
}

func (s *storage) Get() (data []models.DataWrapper) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for v, _ := range s.repo {
		data = append(data, s.repo[v])
	}
	return
}
