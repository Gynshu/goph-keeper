package models

import (
	"github.com/gynshu-one/goph-keeper/shared/utils"
	"time"

	"github.com/google/uuid"
)

type ArbitraryText struct {
	// ID is the primary key
	ID string `json:"id" bson:"_id"`
	// OwnerID is the user who owns this text
	OwnerID string `json:"owner_id" bson:"owner_id"`
	// Name is the name of the text
	Name string `json:"name" bson:"name"`
	// ArbitraryText is the text
	Text string `json:"text" bson:"text"`

	// CreatedAt is the time when this text was created
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	// UpdatedAt is the time when this text was last updated
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	// DeletedAt is the time when this text was deleted
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at"`
}

func (a *ArbitraryText) EncryptAll(passphrase string) error {
	encryptedText, err := utils.EncryptData([]byte(a.Text), passphrase)
	if err != nil {
		return err
	}
	a.Text = string(encryptedText)

	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *ArbitraryText) DecryptAll(passphrase string) error {
	decryptedText, err := utils.DecryptData([]byte(a.Text), passphrase)
	if err != nil {
		return err
	}
	a.Text = string(decryptedText)

	return nil
}

func (a *ArbitraryText) GetOwnerID() string {
	return a.OwnerID
}

func (a *ArbitraryText) GetDataID() UserDataID {
	return UserDataID(a.ID)
}

func (a *ArbitraryText) SetCreatedAt() {
	a.CreatedAt = time.Now().Unix()
}

func (a *ArbitraryText) SetUpdatedAt() {
	a.UpdatedAt = time.Now().Unix()
}

func (a *ArbitraryText) SetDeletedAt() {
	a.DeletedAt = time.Now().Unix()
}

func (a *ArbitraryText) MakeID() {
	a.ID = uuid.New().String()
}

func (a *ArbitraryText) GetType() UserDataType {
	return TextType
}
