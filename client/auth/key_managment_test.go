package auth

import (
	"github.com/gynshu-one/goph-keeper/client/config"
	"testing"
)

func TestSetAndGetPass(t *testing.T) {
	// Set a test password
	testPass := "test_password"
	SetPass(testPass)

	// Get the password from the keyring
	retrievedPass := GetPass()

	// Check if the retrieved password matches the set password
	if retrievedPass != testPass {
		t.Errorf("Expected %s, got %s", testPass, retrievedPass)
	}
}

func TestSetAndGetSecret(t *testing.T) {
	// Set a test secret
	testSecret := "test_secret"
	SetSecret(testSecret)

	// Get the secret from the keyring
	retrievedSecret := GetSecret()

	// Check if the retrieved secret matches the set secret
	if retrievedSecret != testSecret {
		t.Errorf("Expected %s, got %s", testSecret, retrievedSecret)
	}
}

func TestGetNonExistentPass(t *testing.T) {
	config.CurrentUser.Username = "non_existent_user"
	config.CurrentUser.SessionID = "non_existent_session_id"
	// Attempt to get a password that doesn't exist in the keyring
	nonExistentPass := GetPass()

	// Check if the result is an empty string (no error)
	if nonExistentPass != "" {
		t.Errorf("Expected an empty string, got %s", nonExistentPass)
	}
}

func TestGetNonExistentSecret(t *testing.T) {
	config.CurrentUser.Username = "non_existent_user"
	config.CurrentUser.SessionID = "non_existent_session_id"
	// Attempt to get a secret that doesn't exist in the keyring
	nonExistentSecret := GetSecret()

	// Check if the result is an empty string (no error)
	if nonExistentSecret != "" {
		t.Errorf("Expected an empty string, got %s", nonExistentSecret)
	}
}
