package handlers

import (
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"net/http"
)

func FindSession(r *http.Request) (*auth.Session, error) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	session, err := auth.Sessions.GetSession(sessionID.Value)
	if err != nil {
		return nil, err
	}
	return session, nil
}
