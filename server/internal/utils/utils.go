package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/mail"
)

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

// ValidateEmail validates an email
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
