package middlewares

import (
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"net/http"
)

func SessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("Session")
		if sessionID == "" {
			http.Error(w, "Session is empty", http.StatusUnauthorized)
			return
		}
		err := auth.Sessions.CheckSession(sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
