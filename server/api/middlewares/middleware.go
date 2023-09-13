package middlewares

import (
	"github.com/gynshu-one/goph-keeper/server/api/auth"
	"net/http"
)

// SessionCheck checks if the session is valid and user authenticated
func SessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session id from header
		sessionID := r.Header.Get("Authorization")
		if sessionID == "" {
			// Get cookie if no session id in header
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			sessionID = cookie.Value
		}

		err := auth.Sessions.CheckSession(sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// go next
		next.ServeHTTP(w, r)
	})
}
