package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var Sessions Manager

func init() {
	Sessions = &sessionManager{
		mu:      &sync.RWMutex{},
		storage: make(map[string]Session),
	}
}

// Manager is an interface for managing sessions
type Manager interface {
	// CreateSession creates a new session for a user
	// returns a session ID and an error if something went wrong
	CreateSession(userID string) (*Session, error)

	// GetSession returns a session for a given session ID
	// returns an error if something went wrong
	GetSession(sessionID string) (*Session, error)

	// CheckSession checks if a session is valid
	// returns an error if something went wrong
	CheckSession(sessionID string) error

	// DeleteSession deletes a session for a given session ID
	// returns an error if something went wrong
	DeleteSession(sessionID string) error

	// DeleteAllSessions deletes all sessions for a given user ID
	// returns an error if something went wrong
	DeleteAllSessions(userID string) error
}

type sessionManager struct {
	mu      *sync.RWMutex
	storage map[string]Session
}

// CreateSession creates a new session for a user
// returns a session ID and an error if something went wrong
func (s *sessionManager) CreateSession(userID string) (*Session, error) {
	if userID == "" {
		return nil, errors.New("user id is empty")
	}
	// Create instance of Session
	session := Session{
		ID:        uuid.New().String(),
		userID:    userID,
		createdAt: time.Now(),
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add session to storage
	s.storage[session.ID] = session

	return &session, nil
}

// GetSession returns a session for a given session ID
// returns an error if something went wrong
func (s *sessionManager) GetSession(sessionID string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if session exists
	session, ok := s.storage[sessionID]
	if !ok {
		return nil, errors.New("session not found")
	}

	// Check if session is expired
	if time.Now().Sub(session.createdAt) > 24*time.Hour {
		delete(s.storage, sessionID)
		return nil, errors.New("session is expired")
	}
	return &session, nil
}

// DeleteSession deletes a session for a given session ID
// returns an error if something went wrong
func (s *sessionManager) DeleteSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.storage, sessionID)
	return nil
}

// CheckSession checks if a session is valid
// returns an error if something went wrong
func (s *sessionManager) CheckSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Check if session exists
	session, ok := s.storage[sessionID]
	if !ok {
		return errors.New("session not found")
	}

	// Check if session is expired
	if time.Now().Sub(session.createdAt) > 24*time.Hour {
		delete(s.storage, sessionID)
		return errors.New("session is expired")
	}
	return nil
}

// DeleteAllSessions deletes all sessions for a given user ID
// returns an error if something went wrong
func (s *sessionManager) DeleteAllSessions(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.storage {
		if v.userID == userID {
			delete(s.storage, k)
		}
	}
	return nil
}

// Session is a struct for storing session data
type Session struct {
	ID        string `json:"id" bson:"_id"`
	userID    string
	createdAt time.Time
}

// GetUserID returns a user ID for a given session
func (s *Session) GetUserID() string {
	return s.userID
}
