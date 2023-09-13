package auth

import (
	"errors"
	"github.com/gynshu-one/goph-keeper/client/config"
	"os"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestInit(t *testing.T) {
	// check home .goph-keeper dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Failed to get user home dir, got %s", err)
		return
	}
	userFile = homeDir +
		string(os.PathSeparator) +
		"." + config.ServiceName +
		string(os.PathSeparator) + "user.txt"

	// check user file
	file, err := os.ReadFile(userFile)
	if errors.Is(err, os.ErrNotExist) {
		t.Log("No user file found")
		return
	}
	if err == nil {
		if string(file) != CurrentUser.Username {
			t.Errorf("Expected %s, got %s", CurrentUser.Username, string(file))
			return
		}
		return
	}

}

func TestSetAndGetPass(t *testing.T) {
	keyring.MockInit()
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
	keyring.MockInit()
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
	keyring.MockInit()
	CurrentUser.Username = "non_existent_user"
	CurrentUser.SessionID = "non_existent_session_id"
	// Attempt to get a password that doesn't exist in the keyring
	nonExistentPass := GetPass()

	// Check if the result is an empty string (no error)
	if nonExistentPass != "" {
		t.Errorf("Expected an empty string, got %s", nonExistentPass)
	}
}

func TestGetNonExistentSecret(t *testing.T) {
	keyring.MockInit()
	CurrentUser.Username = "non_existent_user"
	CurrentUser.SessionID = "non_existent_session_id"
	// Attempt to get a secret that doesn't exist in the keyring
	nonExistentSecret := GetSecret()

	// Check if the result is an empty string (no error)
	if nonExistentSecret != "" {
		t.Errorf("Expected an empty string, got %s", nonExistentSecret)
	}
}
