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
	// DeletedAt is the time when this binary was deleted
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at"`
}

// EncryptAll encrypts all sensitive data
func (b *Binary) EncryptAll(passphrase string) error {
	encryptedBinary, err := utils.EncryptData(b.Binary, passphrase)
	if err != nil {
		return err
	}
	b.Binary = encryptedBinary
	encryptedInfo, err := utils.EncryptData([]byte(b.Info), passphrase)
	if err != nil {
		return err
	}
	b.Info = string(encryptedInfo)
	b.UpdatedAt = time.Now().Unix()
	return nil
}

// DecryptAll decrypts all sensitive data
func (b *Binary) DecryptAll(passphrase string) error {
	decryptedBinary, err := utils.DecryptData(b.Binary, passphrase)
	if err != nil {
		return err
	}
	b.Binary = decryptedBinary
	decryptedInfo, err := utils.DecryptData([]byte(b.Info), passphrase)
	if err != nil {
		return err
	}
	b.Info = string(decryptedInfo)
	return nil
}

// GetOwnerID returns the owner id
func (b *Binary) GetOwnerID() string {
	return b.OwnerID
}

func (b *Binary) GetDataID() UserDataID {
	return UserDataID(b.ID)
}

func (b *Binary) SetCreatedAt() {
	b.CreatedAt = time.Now().Unix()
}

func (b *Binary) SetUpdatedAt() {
	b.UpdatedAt = time.Now().Unix()
}

func (b *Binary) SetDeletedAt() {
	b.DeletedAt = time.Now().Unix()
}

func (b *Binary) MakeID() {
	b.ID = uuid.New().String()
}

func (b *Binary) GetType() UserDataType {
	return BinaryType
}
