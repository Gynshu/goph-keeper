package handlers

import (
	"errors"
	"github.com/gynshu-one/goph-keeper/common/models"
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"github.com/gynshu-one/goph-keeper/server/api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// CreateUser creates a new user
// hashes it and stores to database
// user must pass email as "email" and password as "password" url parameters (GET request)
// {"email": "email", "password": "password"}
// This would register a new user and create a 24-hour session
// https request example:
// https://localhost:8080/user/create?email=tig.arsenyan@gmail.com&password=password
//
//	in response you will get a session_id cookie and Authorization header with session id
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" {
		http.Error(w, "email is empty", http.StatusBadRequest)
		return
	}
	if !utils.ValidateEmail(email) {
		http.Error(w, "email is invalid", http.StatusBadRequest)
		return
	}
	user.Email = email

	if password == "" {
		http.Error(w, "password is empty", http.StatusBadRequest)
		return
	}

	// Hash master key (no salt and pepper for now)
	user.Passphrase = utils.HashMasterKey(password)

	// clean mem
	password = utils.GenRandomString(len(password) + 1)

	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	// Try to create a new user
	_, err := h.db.Collection("users").InsertOne(r.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "user with this email already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Don't forget to create a session for the user
	session, err := auth.Sessions.CreateSession(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// header
	w.Header().Set("Authorization", "Bearer "+session.ID)
	// cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: session.ID,
	})
	w.Write([]byte(session.ID))
}

// LoginUser logs in a user
// returns a session ID and an error if something went wrong
// user must pass email as "email" and password as "password" url parameters (GET request)
// in response you will get a session_id cookie and Authorization header with session id
func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "email is empty", http.StatusBadRequest)
		return
	}
	if !utils.ValidateEmail(email) {
		http.Error(w, "email is invalid", http.StatusBadRequest)
		return
	}

	err := h.db.Collection("users").FindOne(r.Context(), bson.M{"_id": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "user not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		http.Error(w, "master key is empty", http.StatusBadRequest)
		return
	}
	if !utils.CheckMasterKey(user.Passphrase, password) {
		http.Error(w, "invalid master key", http.StatusBadRequest)
		return
	}

	// clean mem
	password = utils.GenRandomString(len(password) + 1)

	session, err := auth.Sessions.CreateSession(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// header
	w.Header().Set("Authorization", "Bearer "+session.ID)
	// cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: session.ID,
	})
	w.Write([]byte(session.ID))
}

// LogoutUser logs out a user
// user must pass session_id as "session_id" url parameter (GET request)
func (h *handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = auth.Sessions.DeleteSession(sessionID.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
