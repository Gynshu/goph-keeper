package storage

import (
	"context"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/gynshu-one/goph-keeper/server/config"
	"testing"
	"time"
)

func TestSetData(t *testing.T) {
	ctx := context.Background()
	// I don't have time to create mocks for mongoDB client and responses
	db := config.NewDb()
	s := NewStorage(db.Collection("user-data"))
	data := models.DataWrapper{
		ID:        "123456",
		OwnerID:   "user1",
		Data:      []byte("some data"),
		Name:      "My data",
		UpdatedAt: time.Now().Unix(),
	}

	err := s.SetData(ctx, data)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetData(t *testing.T) {
	ctx := context.Background()
	// I don't have time to create mocks for mongoDB client and responses
	db := config.NewDb()
	s := NewStorage(db.Collection("user-data"))
	data := models.DataWrapper{
		ID:        "123456",
		OwnerID:   "user1",
		Data:      []byte("some data"),
		Name:      "My data",
		UpdatedAt: time.Now().Unix(),
	}

	err := s.SetData(ctx, data)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	result, err := s.GetData(ctx, "user1")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
	if result[0].ID != data.ID {
		t.Errorf("expected ID to be %s, got %s", data.ID, result[0].ID)
	}
	if string(data.Data) != string(result[0].Data) {
		t.Errorf("expected Data to be %s, got %s", data.Data, result[0].Data)
	}
}
