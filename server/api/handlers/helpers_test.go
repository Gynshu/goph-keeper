package handlers

import (
	"net/http"
	"testing"

	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
)

func TestFindSession(t *testing.T) {
	// Create a new request with a test session ID cookie
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	// Define a test session
	session, err := auth.Sessions.CreateSession("test-user-id")
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}
	cookie := &http.Cookie{Name: "session_id", Value: session.ID}
	req.AddCookie(cookie)

	// Call the FindSession function
	foundSession, err := FindSession(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Check if the found session matches the test session
	if foundSession.GetUserID() != session.GetUserID() {
		t.Errorf("Found session does not match test session")
	}

	// Check if the found session is of the correct type
	if foundSession.ID != session.ID {
		t.Errorf("Found session is of the wrong type")
	}
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "wrong-id"})

	// Call the FindSession function with a wrong session ID
	_, err = FindSession(req)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
