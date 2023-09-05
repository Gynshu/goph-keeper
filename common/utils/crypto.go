package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// EncryptData encrypts data (Any length from 1 to ~) using a user's master key
func EncryptData(data []byte, key string) ([]byte, error) {
	key, err := deriveAESKey(key)
	if err != nil {
		return nil, err
	}
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
	key, err := deriveAESKey(key)
	if err != nil {
		return nil, err
	}
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

// deriveAESKey derives a 256-bit AES key from a user's master key
func deriveAESKey(userKey string) (string, error) {
	// Hash the user-provided key using SHA-256 to generate a 256-bit key
	hasher := sha256.New()
	_, err := hasher.Write([]byte(userKey))
	if err != nil {
		return "", err
	}
	derivedKey := hasher.Sum(nil)
	return string(derivedKey), nil
}
