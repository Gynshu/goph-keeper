package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// GenerateMasterKeyForUser generates a new master key for a user
func GenerateMasterKeyForUser() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	hashedKey := sha256.Sum256(key)
	return hex.EncodeToString(hashedKey[:]), nil
}

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

// CheckMasterKey checks if a user's master key matches the stored hash
func CheckMasterKey(masterKeyHash string, userMasterKey string) bool {
	hashedUserMasterKey := sha256.Sum256([]byte(userMasterKey))
	return hex.EncodeToString(hashedUserMasterKey[:]) == masterKeyHash[:]
}

// HashMasterKey hashes a user's master key
func HashMasterKey(masterKey string) string {
	hashedMasterKey := sha256.Sum256([]byte(masterKey))
	return hex.EncodeToString(hashedMasterKey[:])
}
