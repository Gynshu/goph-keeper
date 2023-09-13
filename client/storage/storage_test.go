package storage

import (
	"testing"

	"github.com/gynshu-one/goph-keeper/client/auth"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/zalando/go-keyring"
)

func TestAddEncryptAndFindDecrypt(t *testing.T) {
	keyring.MockInit()
	// Set a secret for encryption
	auth.SetSecret("test_secret")

	// Create a new storage instance
	s := NewStorage()

	// Create some test data
	testData := models.Login{
		Username: "test_username",
		Password: "test_password",
	}

	// Create a data wrapper for the test data
	testWrapper := models.DataWrapper{
		ID:   "test_id",
		Type: models.LoginType,
	}

	// Add the test data to the storage
	err := s.AddEncrypt(&testData, testWrapper)
	if err != nil {
		t.Errorf("AddEncrypt returned an error: %v", err)
	}

	// Find and decrypt the test data from the storage
	decryptedData, wrapper, err := s.FindDecrypt("test_id")
	if err != nil {
		t.Errorf("FindDecrypt returned an error: %v", err)
	}

	// Check if the decrypted data matches the original test data
	if decryptedData.(models.Login).Username != testData.Username || decryptedData.(models.Login).Password != testData.Password {
		t.Errorf("Decrypted data does not match the original test data")
	}

	// Check if the wrapper matches the test wrapper
	if wrapper.ID != testWrapper.ID || wrapper.Type != testWrapper.Type {
		t.Errorf("Wrapper does not match the test wrapper")
	}
}

func TestDelete(t *testing.T) {
	keyring.MockInit()
	// Set a secret for encryption
	auth.SetSecret("test_secret")

	// Create a new storage instance
	s := NewStorage()

	// Create some test data
	testData := models.Login{
		Username: "test_username",
		Password: "test_password",
	}

	// Create a data wrapper for the test data
	testWrapper := models.DataWrapper{
		ID:   "test_id",
		Type: models.LoginType,
	}

	// Add the test data to the storage
	err := s.AddEncrypt(&testData, testWrapper)
	if err != nil {
		t.Errorf("AddEncrypt returned an error: %v", err)
	}

	// Delete the test data from the storage
	err = s.Delete("test_id")
	if err != nil {
		t.Errorf("Delete returned an error: %v", err)
	}

	// Try to find the deleted test data
	_, wrapper, err := s.FindDecrypt("test_id")

	// Check if the error is the expected "item deleted" error
	if err != models.ErrDeleted {
		t.Errorf("Expected error: %v, got: %v", models.ErrDeleted, err)
	}

	// Check if the wrapper has the DeletedAt field set
	if wrapper.DeletedAt == 0 {
		t.Errorf("Wrapper's DeletedAt field is not set")
	}
}

func TestSwapAndGetData(t *testing.T) {
	keyring.MockInit()
	// Set a secret for encryption
	auth.SetSecret("test_secret")

	// Create a new storage instance
	s := NewStorage()

	// Create some test data
	testData := []models.DataWrapper{
		{
			ID:   "test_id_1",
			Type: models.LoginType,
		},
		{
			ID:   "test_id_2",
			Type: models.BankCardType,
		},
	}

	// Swap the test data into the storage
	err := s.Swap(testData)
	if err != nil {
		t.Errorf("Swap returned an error: %v", err)
	}

	// Get the data from the storage
	data := s.Get()

	// Check if the length of the data matches the length of the test data
	if len(data) != len(testData) {
		t.Errorf("Length of data does not match length of test data")
	}

	// Check if the IDs and types match for each data wrapper
	for i, wrapper := range data {
		if wrapper.ID != testData[i].ID || wrapper.Type != testData[i].Type {
			t.Errorf("Data wrapper does not match test data wrapper")
		}
	}
}

func TestGetNonExistentData(t *testing.T) {
	keyring.MockInit()
	// Create a new storage instance
	s := NewStorage()

	// Get data from the storage
	data := s.Get()

	// Check if the data is an empty slice (no data in storage)
	if len(data) != 0 {
		t.Errorf("Expected an empty slice, got: %v", data)
	}
}

func TestSwap(t *testing.T) {
	keyring.MockInit()
	// Create a new storage instance
	s := NewStorage()

	// Define test data
	testData := []models.DataWrapper{
		{ID: "1", Type: "foo", Data: []byte("bar")},
		{ID: "2", Type: "baz", Data: []byte("qux")},
	}

	// Swap the data in the storage
	err := s.Swap(testData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Get the data from the storage
	data := s.Get()

	// Check if the length of the data matches the length of the test data
	if len(data) != len(testData) {
		t.Errorf("Length of data does not match length of test data")
	}

	for _, d := range data {
		if !(d.ID == "1" || d.ID == "2") {
			t.Errorf("Unexpected data: %v", d)
		}
	}
}
