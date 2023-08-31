package models

import (
	"github.com/gynshu-one/goph-keeper/shared/utils"
	"time"

	"github.com/google/uuid"
)

type Binary struct {
	// ID is the primary key
	ID string `json:"id" bson:"_id"`
	// OwnerID is the user who owns this binary
	OwnerID string `json:"owner_id" bson:"owner_id"`
	// Name is the name of the binary
	Name string `json:"name" bson:"name"`
	// Info is the additional info about the binary
	Info string `json:"info" bson:"info"`
	// Binary is the binary data
	Binary []byte `json:"binary" bson:"binary"`
	// CreatedAt is the time when this binary was created
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	// UpdatedAt is the time when this binary was last updated
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
}

// EncryptAll encrypts all sensitive data
func (data *Binary) EncryptAll(passphrase string) error {
	encryptedBinary, err := utils.EncryptData(data.Binary, passphrase)
	if err != nil {
		return err
	}
	data.Binary = encryptedBinary
	encryptedInfo, err := utils.EncryptData([]byte(data.Info), passphrase)
	if err != nil {
		return err
	}
	data.Info = string(encryptedInfo)
	data.UpdatedAt = time.Now().Unix()
	return nil
}

// DecryptAll decrypts all sensitive data
func (data *Binary) DecryptAll(passphrase string) error {
	decryptedBinary, err := utils.DecryptData(data.Binary, passphrase)
	if err != nil {
		return err
	}
	data.Binary = decryptedBinary
	decryptedInfo, err := utils.DecryptData([]byte(data.Info), passphrase)
	if err != nil {
		return err
	}
	data.Info = string(decryptedInfo)
	return nil
}

// GetOwnerID returns the owner id
func (data *Binary) GetOwnerID(id *string) string {
	if id != nil {
		data.OwnerID = *id
	}
	return data.OwnerID
}

func (data *Binary) GetDataID() UserDataID {
	return UserDataID(data.ID)
}

func (data *Binary) SetCreatedAt() {
	data.CreatedAt = time.Now().Unix()
}

func (data *Binary) SetUpdatedAt() {
	data.UpdatedAt = time.Now().Unix()
}

func (data *Binary) MakeID() {
	data.ID = uuid.New().String()
}

func (data *Binary) GetType() UserDataType {
	return BinaryType
}
