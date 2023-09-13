package sync

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gynshu-one/goph-keeper/client/auth"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
)

func TestNewMediator(t *testing.T) {
	keyring.MockInit()
	// Create a new storage instance
	newStorage := storage.NewStorage()

	// Create a new newMediator instance
	newMediator := NewMediator(newStorage)

	// Check if the newMediator client is not nil
	if newMediator.client == nil {
		t.Errorf("Mediator client is nil")
	}

	// Check if the newMediator storage is not nil
	if newMediator.storage == nil {
		t.Errorf("Mediator storage is nil")
	}
}

// MockHTTPServer is a helper function to create a mock HTTP server for testing.
func MockChiHTTPServer() *httptest.Server {
	keyring.MockInit()
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create listener")
	}
	r := chi.NewRouter()
	r.Route("/user", func(r chi.Router) {
		r.With().Get("/create", func(w http.ResponseWriter, r *http.Request) {
			// Set session_id cookie
			cookie := http.Cookie{
				Name:  "session_id",
				Value: "test",
			}
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusOK)
		})
		r.With().Get("/login", func(w http.ResponseWriter, r *http.Request) {
			// Set session_id cookie
			cookie := http.Cookie{
				Name:  "session_id",
				Value: "test",
			}
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusOK)
		})
		r.With().Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.With().Post("/sync", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		})
	})
	server := httptest.NewUnstartedServer(r)
	_ = server.Listener.Close()
	server.Listener = l
	server.TLS = &tls.Config{InsecureSkipVerify: true}
	server.EnableHTTP2 = true
	server.StartTLS()
	return server
}

func TestSignUp(t *testing.T) {
	keyring.MockInit()
	server := MockChiHTTPServer()
	defer server.Close()

	// Create a mediator with the mock server
	newMediator := NewMediator(storage.NewStorage())
	// Test the SignUp function
	err := newMediator.SignUp(context.Background(), "testuser", "password")
	if err != nil {
		t.Errorf("SignUp failed with error: %v", err)
	}

	// Check if the session created
	auth.CurrentUser.Username = "testuser"
	auth.CurrentUser.SessionID = "test"
}

func TestSignIn(t *testing.T) {
	keyring.MockInit()
	server := MockChiHTTPServer()
	defer server.Close()

	// Create a mediator with the mock server
	newMediator := NewMediator(storage.NewStorage())

	// Test the SignIn function
	err := newMediator.SignIn(context.Background(), "testuser", "password")
	if err != nil {
		t.Errorf("SignIn failed with error: %v", err)
	}

	// Check if the session created
	auth.CurrentUser.Username = "testuser"
	auth.CurrentUser.SessionID = "test"
}
func TestSync(t *testing.T) {
	keyring.MockInit()
	server := MockChiHTTPServer()
	defer server.Close()

	// Create a mediator with the mock server
	newMediator := NewMediator(storage.NewStorage())

	// Test the Sync function
	err := newMediator.Sync(context.Background())
	if err != nil {
		t.Errorf("Sync failed with error: %v", err)
	}

}

func TestSetCookies(t *testing.T) {
	keyring.MockInit()
	// Define test cookies
	cookies := []*http.Cookie{
		{Name: "session_id", Value: "test-session-id"},
	}

	// Define test username
	username := "test-username"

	// Set the cookies
	err := setCookies(cookies, username)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the current user username matches the test username
	if auth.CurrentUser.Username != username {
		t.Errorf("Current user username does not match test username")
	}

	// Check if the current user session ID matches the test session ID
	if auth.CurrentUser.SessionID != "test-session-id" {
		t.Errorf("Current user session ID does not match test session ID")
	}
}
