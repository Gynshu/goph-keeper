package auth

import (
	"github.com/google/uuid"
	"testing"
)

func TestSessionManager_CreateSession(t *testing.T) {
	// Test creating a new session.
	userID := "testUserID"
	session, err := Sessions.CreateSession(userID)
	if err != nil {
		t.Errorf("CreateSession failed: %v", err)
	}

	// Check if the session was created with a valid ID.
	if session == nil || session.ID == "" {
		t.Error("Invalid session ID")
	}
}

func TestSessionManager_GetSession(t *testing.T) {
	// Create a test session.
	userID := "testUserID"
	session, _ := Sessions.CreateSession(userID)
	sessionID := session.ID

	// Test retrieving an existing session.
	retrievedSession, err := Sessions.GetSession(sessionID)
	if err != nil {
		t.Errorf("GetSession failed: %v", err)
	}

	// Check if the retrieved session matches the created one.
	if retrievedSession == nil || retrievedSession.ID != sessionID || retrievedSession.userID != userID {
		t.Error("Retrieved session does not match")
	}

	// Test retrieving a non-existent session.
	_, err = Sessions.GetSession(uuid.New().String())
	if err == nil {
		t.Error("GetSession should have returned an error for a non-existent session")
	}
}

func TestSessionManager_CheckSession(t *testing.T) {
	// Create a test session.
	userID := "testUserID"
	session, _ := Sessions.CreateSession(userID)

	// Test checking a valid session.
	err := Sessions.CheckSession(session.ID)
	if err != nil {
		t.Errorf("CheckSession failed for a valid session: %v", err)
	}

	// We can't test checking an expired session because the session manager doesn't share
	// field values with the session struct and returns only values

	// Test checking a non-existent session.
	err = Sessions.CheckSession(uuid.New().String())
	if err == nil {
		t.Error("CheckSession should have returned an error for a non-existent session")
	}
}

func TestSessionManager_DeleteSession(t *testing.T) {
	// Create a test session.
	userID := "testUserID"
	session, _ := Sessions.CreateSession(userID)
	sessionID := session.ID

	// Test deleting an existing session.
	err := Sessions.DeleteSession(sessionID)
	if err != nil {
		t.Errorf("DeleteSession failed: %v", err)
	}

	// Test deleting a non-existent session.
	err = Sessions.CheckSession(sessionID)
	if err == nil {
		t.Error("DeleteSession should have returned an error for a non-existent session")
	}
}

func TestSessionManager_DeleteAllSessions(t *testing.T) {
	// Create test sessions.
	userID := "testUserID"
	session1, _ := Sessions.CreateSession(userID)
	session2, _ := Sessions.CreateSession(userID)

	// Test deleting all sessions for a user.
	err := Sessions.DeleteAllSessions(userID)
	if err != nil {
		t.Errorf("DeleteAllSessions failed: %v", err)
	}

	// Test that the sessions were deleted.
	_, err = Sessions.GetSession(session1.ID)
	if err == nil {
		t.Error("Session 1 should have been deleted")
	}

	_, err = Sessions.GetSession(session2.ID)
	if err == nil {
		t.Error("Session 2 should have been deleted")
	}
}
