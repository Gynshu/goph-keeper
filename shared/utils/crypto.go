package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"io"
)

// EncryptData encrypts data (Any length from 1 to ~) using a user's master key
func EncryptData(data []byte, key string) ([]byte, error) {
	// Generate a new AES cipher using the master key
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// Generate a new GCM cipher using the AES cipher
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// Generate a nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// DecryptData decrypts data (any length from 1 to ~) using a user's master key
func DecryptData(ciphertext []byte, key string) ([]byte, error) {
	// Generate a new AES cipher using the master key
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// Generate a new GCM cipher using the AES cipher
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// Get the nonce size
	nonceSize := gcm.NonceSize()

	// Get the nonce
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func PackData(allData []models.UserData) models.PackedUserData {
	var out = make(models.PackedUserData)
	for i := range allData {
		out[allData[i].GetType()][allData[i].GetDataID()] = allData[i]
	}
	return out
}
