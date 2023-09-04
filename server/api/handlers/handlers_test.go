package handlers

import (
	"bytes"
	"encoding/json"
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gynshu-one/goph-keeper/common/models"
)

func TestSyncUserData(t *testing.T) {
	// Create a mock storage implementation.
	mock := &mockStorage{
		data: make(map[string][]models.DataWrapper),
	}

	// Create a test handler with the mock storage.
	hndlr := NewHandlers(mock)

	// Create a test user session.
	userID := "testUserID"
	auth.Sessions = auth.NewSessionManager()
	session, _ := auth.Sessions.CreateSession(userID)

	// Create test data to send in the request body.
	testData := []models.DataWrapper{
		{
			ID:        "1",
			OwnerID:   userID,
			Name:      "Test Data 1",
			Data:      []byte("Data 1"),
			UpdatedAt: time.Now().Unix(),
		},
		{
			ID:        "2",
			OwnerID:   userID,
			Name:      "Test Data 2",
			Data:      []byte("Data 2"),
			UpdatedAt: time.Now().Unix(),
		},
	}

	// Marshal test data to JSON.
	requestData, err := json.Marshal(testData)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test HTTP request with the session cookie and JSON data.
	request := httptest.NewRequest(http.MethodPost, "/sync", bytes.NewReader(requestData))
	request.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: session.ID,
	})
	response := httptest.NewRecorder()

	// Call the SyncUserData handler function.
	hndlr.SyncUserData(response, request)

	// Check the HTTP response status code.
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	// Unmarshal the response body to check if data was saved correctly.
	var responseData []models.DataWrapper
	err = json.Unmarshal(response.Body.Bytes(), &responseData)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the response data matches the test data.
	if len(responseData) != len(testData) {
		t.Errorf("Expected %d items in response, got %d", len(testData), len(responseData))
	}

	// Verify that the data was saved in the mock storage.
	storedData, err := mock.GetData(nil, userID)
	if err != nil {
		t.Errorf("Failed to retrieve data from storage: %v", err)
	}

	if len(storedData) != len(testData) {
		t.Errorf("Expected %d items in storage, got %d", len(testData), len(storedData))
	}
}
