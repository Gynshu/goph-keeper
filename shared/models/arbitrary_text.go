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
}

func (data *ArbitraryText) GetName() string {
	return data.Name
}

func (data *ArbitraryText) EncryptAll(passphrase string) error {
	encryptedText, err := utils.EncryptData([]byte(data.Text), passphrase)
	if err != nil {
		return err
	}
	data.Text = string(encryptedText)

	data.UpdatedAt = time.Now().Unix()
	return nil
}

func (data *ArbitraryText) DecryptAll(passphrase string) error {
	decryptedText, err := utils.DecryptData([]byte(data.Text), passphrase)
	if err != nil {
		return err
	}
	data.Text = string(decryptedText)

	return nil
}

func (data *ArbitraryText) GetOrSetOwnerID(id *string) string {
	if id != nil {
		data.OwnerID = *id
	}
	return data.OwnerID
}

func (data *ArbitraryText) GetDataID() UserDataID {
	return UserDataID(data.ID)
}

func (data *ArbitraryText) SetCreatedAt() {
	data.CreatedAt = time.Now().Unix()
}

func (data *ArbitraryText) SetUpdatedAt() {
	data.UpdatedAt = time.Now().Unix()
}

func (data *ArbitraryText) GetUpdatedAt() int64 {
	return data.UpdatedAt
}

func (data *ArbitraryText) MakeID() {
	data.ID = uuid.New().String()
}

func (data *ArbitraryText) GetType() UserDataType {
	return ArbitraryTextType
}
