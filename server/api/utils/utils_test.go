package utils

import "testing"

func TestCheckMasterKey(t *testing.T) {
	// Test case 1: Matching master key and hash.
	masterKey := "password123"
	hashedMasterKey := HashMasterKey(masterKey)
	matched := CheckMasterKey(hashedMasterKey, masterKey)
	if !matched {
		t.Error("Expected master key and hash to match, but they didn't.")
	}

	// Test case 2: Mismatched master key and hash.
	invalidMasterKey := "invalid123"
	mismatched := CheckMasterKey(hashedMasterKey, invalidMasterKey)
	if mismatched {
		t.Error("Expected master key and hash to not match, but they did.")
	}
}

func TestHashMasterKey(t *testing.T) {
	// Test hashing a master key.
	masterKey := "password123"
	hashedMasterKey := HashMasterKey(masterKey)
	if hashedMasterKey == "" {
		t.Error("Expected a non-empty hashed master key, but got an empty string.")
	}
}

func TestValidateEmail(t *testing.T) {
	// Test case 1: Valid email.
	validEmail := "test@example.com"
	valid := ValidateEmail(validEmail)
	if !valid {
		t.Error("Expected a valid email, but it was considered invalid.")
	}

	// Test case 2: Invalid email.
	invalidEmail := "invalid-email"
	invalid := ValidateEmail(invalidEmail)
	if invalid {
		t.Error("Expected an invalid email, but it was considered valid.")
	}
}

func TestGenRandomString(t *testing.T) {
	// Test generating a random string.
	length := 10
	randomString := GenRandomString(length)
	if len(randomString) != length {
		t.Errorf("Expected a random string of length %d, but got length %d.", length, len(randomString))
	}
}
