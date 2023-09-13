package router

import (
	"github.com/gynshu-one/goph-keeper/common/models"
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"github.com/gynshu-one/goph-keeper/server/api/handlers"
	"github.com/gynshu-one/goph-keeper/server/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	// Create a mock storage implementation.
	mock := &storage.MockStorage{
		Data: make(map[string][]models.DataWrapper),
	}

	// Create a test handler with the mock storage.
	hndlr := handlers.NewHandlers(mock)

	// Create a new router.
	r := NewRouter(hndlr)

	// Create a new test HTTP server using the router
	ts := httptest.NewServer(r)
	defer ts.Close()

	session, err := auth.Sessions.CreateSession("test-user-id")
	if err != nil {
		return
	}

	cookie := &http.Cookie{Name: "session_id", Value: session.ID}

	req, err := http.NewRequest("GET", ts.URL+"/user/logout", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	// Send a GET request to the test server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	req.AddCookie(cookie)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
