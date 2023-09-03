package models

import (
	"encoding/json"
	"github.com/gynshu-one/goph-keeper/common/utils"
)

type Binary struct {
	// Name is the name of the binary
	Name string `json:"name" bson:"name"`
	// Info is the additional info about the binary
	Info string `json:"info" bson:"info"`
	// Binary is the binary data
	Binary []byte `json:"binary" bson:"binary"`
}

// EncryptAll encrypts all sensitive data
func (data *Binary) EncryptAll(passphrase string) (encryptedData []byte, err error) {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return utils.EncryptData(marshaled, passphrase)
}

// DecryptAll decrypts all sensitive data
func (data *Binary) DecryptAll(passphrase string, encrypteData []byte) error {
	decrypted, err := utils.DecryptData(encrypteData, passphrase)
	if err != nil {
		return err
	}
	return json.Unmarshal(decrypted, data)
}
