package handlers

import (
	"context"
	"fmt"
	"github.com/gynshu-one/goph-keeper/common/models"
)

type mockStorage struct {
	user models.User
	data map[string][]models.DataWrapper
}

func (m *mockStorage) SetData(ctx context.Context, data models.DataWrapper) error {
	userID := data.OwnerID
	if _, ok := m.data[userID]; !ok {
		m.data[userID] = []models.DataWrapper{}
	}
	m.data[userID] = append(m.data[userID], data)
	return nil
}

func (m *mockStorage) GetData(ctx context.Context, userID string) ([]models.DataWrapper, error) {
	data, ok := m.data[userID]
	if !ok {
		return nil, fmt.Errorf("no data for user %s", userID)
	}
	return data, nil
}

func (m *mockStorage) CreateUser(ctx context.Context, user models.User) error {
	m.user = user
	return nil
}

func (m *mockStorage) GetUser(ctx context.Context, userID string) (models.User, error) {
	return m.user, nil
}
