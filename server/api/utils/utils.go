package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"math/rand"
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

func PackData(allData []models.UserDataModel) models.PackedUserData {
	var out = make(models.PackedUserData)
	for i := range allData {
		out[allData[i].GetType()] = append(out[allData[i].GetType()],
			allData[i])
	}
	return out
}

// GenRandomString generates a random string of length n
func GenRandomString(n int) string {
	// Create a slice of runes to represent the possible characters in the string.
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// Create a byte slice to store the random string.
	b := make([]rune, n)

	// Generate a random number for each position in the byte slice.
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	// Return the random string as a string.
	return string(b)
}
