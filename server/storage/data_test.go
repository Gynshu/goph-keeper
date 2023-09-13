package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"

	"github.com/gynshu-one/goph-keeper/common/models"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSetData(t *testing.T) {
	ctx := context.Background()
	// Mock mongo
	opts := mtest.NewOptions().ClientType(mtest.Mock)
	mt := mtest.New(t, opts)
	defer mt.Close()
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{"ok", 1}})
		db := mt.Client.Database("test")

		s := NewStorage(db.Collection("user-data"), db.Collection("users"))
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
	})
}

func TestGetData(t *testing.T) {
	ctx := context.Background()
	// Mock mongo
	opts := mtest.NewOptions().ClientType(mtest.Mock)
	mt := mtest.New(t, opts)
	defer mt.Close()
	mt.Run("test", func(mt *mtest.T) {
		db := mt.Client.Database("test")
		s := NewStorage(db.Collection("user-data"), db.Collection("users"))
		mt.AddMockResponses(mtest.CreateCursorResponse(1,
			"DBName.CollectionName",
			mtest.FirstBatch, bson.D{
				{"_id", "123456"},
				{"ok", 1},
				{"key", "value"}}))
		_, err := s.GetData(ctx, "user1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

}
