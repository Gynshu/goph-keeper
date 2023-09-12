package middlewares_test

import (
	"github.com/gynshu-one/goph-keeper/common/models"
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"github.com/gynshu-one/goph-keeper/server/api/handlers"
	"github.com/gynshu-one/goph-keeper/server/api/router"
	"github.com/gynshu-one/goph-keeper/server/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionCheck(t *testing.T) {
	//// Create a mock database.
	stor := &storage.MockStorage{
		User: models.User{},
	}

	// Create a test handler with the mock database.
	hndlr := handlers.NewHandlers(stor)

	// Create a new router using the NewRouter function
	r := router.NewRouter(hndlr)

	// Wrap the router with the SessionCheck middleware
	r.With()

	// Create a new test HTTP server using the router
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Create a new session for the test user
	session, err := auth.Sessions.CreateSession("testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Send a GET request to the /user endpoint with the session ID in the Authorization header
	req, err := http.NewRequest("GET", ts.URL+"/user/logout", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	cookie := &http.Cookie{Name: "session_id", Value: session.ID}
	req.AddCookie(cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	// Verify that the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
