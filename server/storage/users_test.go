package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"

	"github.com/gynshu-one/goph-keeper/common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateUser(t *testing.T) {
	// Mock mongo
	opts := mtest.NewOptions().ClientType(mtest.Mock)
	mt := mtest.New(t, opts)
	defer mt.Close()
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{"ok", 1}})
		db := mt.Client.Database("test")
		s := NewStorage(db.Collection("user-data"), db.Collection("users"))
		// Define a test user
		user := models.User{
			Email:      "test1@example.com",
			Passphrase: "test-password",
			CreatedAt:  1234567890,
		}

		// Create the test user
		err := s.CreateUser(context.Background(), user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

}

func TestGetUser(t *testing.T) {
	// Mock mongo
	opts := mtest.NewOptions().ClientType(mtest.Mock)
	mt := mtest.New(t, opts)
	defer mt.Close()
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1,
			"DBName.CollectionName",
			mtest.FirstBatch, bson.D{
				{"_id", "123456"},
				{"ok", 1},
				{"key", "value"}}))
		db := mt.Client.Database("test")
		s := NewStorage(db.Collection("user-data"), db.Collection("users"))

		// Create the test user
		_, err := s.GetUser(context.Background(), "test1@example.com")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}
