package utils

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
)

func TestEncryptAndDecryptData(t *testing.T) {
	// Generate a random master key for testing.
	masterKey := genRandomString(32)

	// Test data to encrypt and decrypt.
	testData := []byte("This is a test data for encryption and decryption.")

	// Encrypt the test data.
	encryptedData, err := EncryptData(testData, masterKey)
	if err != nil {
		t.Fatalf("Error encrypting data: %v", err)
	}

	// Decrypt the encrypted data.
	decryptedData, err := DecryptData(encryptedData, masterKey)
	if err != nil {
		t.Fatalf("Error decrypting data: %v", err)
	}

	// Convert the test data and decrypted data to hexadecimal strings for comparison.
	testDataHex := hex.EncodeToString(testData)
	decryptedDataHex := hex.EncodeToString(decryptedData)

	// Check if the decrypted data matches the original test data.
	if testDataHex != decryptedDataHex {
		t.Errorf("Decrypted data does not match original data.\nExpected: %s\nGot: %s", testDataHex, decryptedDataHex)
	}
}

func genRandomString(length int) string {
	// Generate a random string of given length.
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(randomBytes)
}
