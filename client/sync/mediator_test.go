package sync

import (
	"context"
	"crypto/tls"
	"github.com/go-chi/chi/v5"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// MockHTTPServer is a helper function to create a mock HTTP server for testing.
func MockChiHTTPServer() *httptest.Server {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal()
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
func init() {
	pwd, _ := os.Getwd()
	// strip last path
	pwd = pwd[:len(pwd)-len("/sync")]
	err := config.NewConfig(pwd + "/config.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file please check if it exists and is valid" +
			"Config should be in json format and contain SERVER_IP, POLL_TIMER, DUMP_TIMER")
	}
}
func TestSignUp(t *testing.T) {
	tempDir := "/tmp"
	server := MockChiHTTPServer()
	defer server.Close()

	// Create a mediator with the mock server
	newMediator := NewMediator(storage.NewStorage())
	// Test the SignUp function
	err := newMediator.SignUp(context.Background(), "testuser", "password")
	if err != nil {
		t.Errorf("SignUp failed with error: %v", err)
	}

	// Check if the session file was created
	sessionFilePath := tempDir + "/" + config.SessionFile
	_, err = os.Stat(sessionFilePath)
	if err != nil {
		t.Errorf("Session file was not created: %v", err)
	}

	// Clean up: Remove the session file
	_ = os.Remove(sessionFilePath)
}

func TestSignIn(t *testing.T) {
	server := MockChiHTTPServer()
	defer server.Close()
	// Create a temporary directory for session files (change as needed)
	tempDir := "/tmp"

	// Create a mediator with the mock server
	newMediator := NewMediator(storage.NewStorage())

	// Test the SignIn function
	err := newMediator.SignIn(context.Background(), "testuser", "password")
	if err != nil {
		t.Errorf("SignIn failed with error: %v", err)
	}

	// Check if the session file was created
	sessionFilePath := tempDir + "/" + config.SessionFile
	_, err = os.Stat(sessionFilePath)
	if err != nil {
		t.Errorf("Session file was not created: %v", err)
	}

	// Clean up: Remove the session file
	_ = os.Remove(sessionFilePath)
}
func TestSync(t *testing.T) {
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

func TestCreateUserSessionFiles(t *testing.T) {
	tempDir := "/tmp"

	server := MockChiHTTPServer()
	defer server.Close()

	// Create a mediator instance
	newMediator := NewMediator(storage.NewStorage())

	// Test the createUserSessionFiles function
	err := newMediator.createUserSessionFiles()
	if err != nil {
		t.Errorf("createUserSessionFiles failed with error: %v", err)
	}

	// Check if the session file was created
	sessionFilePath := tempDir + "/" + config.SessionFile
	_, err = os.Stat(sessionFilePath)
	if err != nil {
		t.Errorf("Session file was not created: %v", err)
	}

	// Clean up: Remove the session file
	_ = os.Remove(sessionFilePath)
}
