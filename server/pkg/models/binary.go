package models

import (
	"github.com/gynshu-one/goph-keeper/server/pkg/utils"
	"time"
)

type Binary struct {
	// ID is the primary key
	ID string `json:"id" bson:"_id"`
	// OwnerID is the user who owns this binary
	OwnerID string `json:"owner_id" bson:"owner_id"`
	// Name is the name of the binary
	Name string `json:"name" bson:"name"`
	// Binary is the binary data
	Binary []byte `json:"binary" bson:"binary"`
	// CreatedAt is the time when this binary was created
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	// UpdatedAt is the time when this binary was last updated
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	// DeletedAt is the time when this binary was deleted
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at"`
}

func (b *Binary) EncryptAll(passphrase string) error {
	encryptedBinary, err := utils.EncryptData(b.Binary, passphrase)
	if err != nil {
		return err
	}
	b.Binary = encryptedBinary

	b.UpdatedAt = time.Now().Unix()
	return nil
}

func (b *Binary) DecryptAll(passphrase string) error {
	decryptedBinary, err := utils.DecryptData(b.Binary, passphrase)
	if err != nil {
		return err
	}
	b.Binary = decryptedBinary

	return nil
}

func (b *Binary) GetOwnerID() string {
	return b.OwnerID
}
