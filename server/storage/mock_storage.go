package storage

import (
	"context"
	"fmt"
	"github.com/gynshu-one/goph-keeper/common/models"
)

type MockStorage struct {
	User models.User
	Data map[string][]models.DataWrapper
}

func (m *MockStorage) SetData(ctx context.Context, data models.DataWrapper) error {
	userID := data.OwnerID
	if _, ok := m.Data[userID]; !ok {
		m.Data[userID] = []models.DataWrapper{}
	}
	m.Data[userID] = append(m.Data[userID], data)
	return nil
}

func (m *MockStorage) GetData(ctx context.Context, userID string) ([]models.DataWrapper, error) {
	data, ok := m.Data[userID]
	if !ok {
		return nil, fmt.Errorf("no data for user %s", userID)
	}
	return data, nil
}

func (m *MockStorage) CreateUser(ctx context.Context, user models.User) error {
	m.User = user
	return nil
}

func (m *MockStorage) GetUser(ctx context.Context, userID string) (models.User, error) {
	return m.User, nil
}
