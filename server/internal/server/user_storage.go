package server

import (
	"context"
	"errors"

	"github.com/gynshu-one/goph-keeper/server/pkg/models"
	"github.com/gynshu-one/goph-keeper/server/pkg/storage"
	"github.com/gynshu-one/goph-keeper/server/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserKeeper interface {
	CreateUser(ctx context.Context, user models.User) error

	SetUserData(ctx context.Context, dataID, userID, secret string, data models.UserData) error

	GetUserData(ctx context.Context, dataID, userID, secret string) (*models.UserData, error)
}

type userStorage struct {
	db   *mongo.Database
	data storage.Storage
}

func NewUserKeeper(data storage.Storage) *userStorage {
	return &userStorage{
		data: data,
	}
}

// CreateUser creates a new user
func (s *userStorage) CreateUser(ctx context.Context, user models.User) error {
	// TODO: implement
	return nil
}

// SetUserData sets (creates or updates) the data for a user,
// if the user does not exist, it returns ErrInvalidLogin
// if the secret is invalid, it returns ErrInvalidSecret
// data should contain the ownerID
func (s *userStorage) SetUserData(ctx context.Context, dataID, secret string, data models.UserData) error {
	var user models.User

	userFromDB := s.db.Collection("users").FindOne(ctx, bson.M{"_id": data.GetOwnerID()})
	if errors.Is(userFromDB.Err(), mongo.ErrNoDocuments) {
		return ErrInvalidLogin
	}

	userFromDB.Decode(&user)

	ok := utils.CheckMasterKey(user.EncryptedPassphrase, secret)
	if !ok {
		return ErrInvalidSecret
	}

	err := data.EncryptAll(user.EncryptedPassphrase)
	if err != nil {
		return err
	}

	err = s.data.Set(ctx, dataID, data)
	if err != nil {
		return err
	}

	return nil
}

// GetUserData returns the data for a user,
// if the user does not exist, it returns ErrInvalidLogin
// if the secret is invalid, it returns ErrInvalidSecret
func (s *userStorage) GetUserData(ctx context.Context, dataID, secret string) (*models.UserData, error) {
	var user models.User

	data, err := s.data.Get(ctx, dataID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrObjectMiss
	}

	userFromDB := s.db.Collection("users").FindOne(ctx, bson.M{"_id": data.GetOwnerID()})
	if errors.Is(userFromDB.Err(), mongo.ErrNoDocuments) {
		return nil, ErrInvalidLogin
	}

	userFromDB.Decode(&user)

	ok := utils.CheckMasterKey(user.EncryptedPassphrase, secret)
	if !ok {
		return nil, ErrInvalidSecret
	}

	return &data, nil
}
