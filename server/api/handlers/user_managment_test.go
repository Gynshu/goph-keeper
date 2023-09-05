package handlers

import (
	"context"
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/gynshu-one/goph-keeper/server/api/utils"
)

func TestCreateUser(t *testing.T) {
	auth.Sessions = auth.NewSessionManager()
	// Create a mock database.
	stor := &mockStorage{
		user: models.User{},
	}

	// Create a test handlers with the mock database.
	handlers := NewHandlers(stor)

	// Create test data for the request.
	formData := url.Values{
		"email":    {"test@example.com"},
		"password": {"password123"},
	}

	// Create a test HTTP request.
	request := httptest.NewRequest(http.MethodGet, "/user/create?"+formData.Encode(), nil)
	response := httptest.NewRecorder()

	// Call the CreateUser handlers function.
	handlers.CreateUser(response, request)

	// Check the HTTP response status code.
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	// Check if the user was created in the mock database.
	user, err := stor.GetUser(context.Background(), formData.Get("email"))
	if err != nil {
		t.Error("User not created")
	}
	// Check if the password was hashed.
	if user.Passphrase == "password123" {
		t.Error("Password was not hashed")
	}
	// Check if the session was created.
	err = auth.Sessions.CheckSession(response.Result().Cookies()[0].Value)
	if err != nil {
		t.Error("Session not created")
	}
}

func TestLoginUser(t *testing.T) {
	auth.Sessions = auth.NewSessionManager()
	// Create a mock database.
	stor := &mockStorage{
		user: models.User{},
	}

	// Create a test handler with the mock database.
	handlers := NewHandlers(stor)

	// Create a test user and add it to the mock database.
	user := models.User{
		Email:      "test@example.com",
		Passphrase: utils.HashMasterKey("password123"),
	}
	err := stor.CreateUser(context.Background(), user)
	if err != nil {
		return
	}

	// Create test data for the request.
	formData := url.Values{
		"email":    {"test@example.com"},
		"password": {"password123"},
	}

	// Create a test HTTP request.
	request := httptest.NewRequest(http.MethodGet, "/user/login?"+formData.Encode(), nil)
	response := httptest.NewRecorder()

	// Call the LoginUser handler function.
	handlers.LoginUser(response, request)

	// Check the HTTP response status code.
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	// Check if the session was created.
	err = auth.Sessions.CheckSession(response.Result().Cookies()[0].Value)
	if err != nil {
		t.Error("Session not created")
	}
}

func TestLogoutUser(t *testing.T) {
	auth.Sessions = auth.NewSessionManager()
	// Create a mock database.
	stor := &mockStorage{
		user: models.User{},
	}

	// Create a test handler with the mock database.
	handlers := NewHandlers(stor)

	// Create a test user and add it to the mock database.
	user := models.User{
		Email:      "test@example.com",
		Passphrase: utils.HashMasterKey("password123"),
	}
	err := stor.CreateUser(context.Background(), user)
	if err != nil {
		return
	}

	// Create a test session and add it to the session manager.
	session, _ := auth.Sessions.CreateSession(user.Email)

	// Create a test HTTP request with a session cookie.
	request := httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	request.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: session.ID,
	})
	response := httptest.NewRecorder()

	// Call the LogoutUser handler function.
	handlers.LogoutUser(response, request)

	// Check the HTTP response status code.
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	// Check if the session was deleted.
	err = auth.Sessions.CheckSession(session.ID)
	if err == nil {
		t.Error("Session not deleted")
	}
}
