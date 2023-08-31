package handlers

import (
	"encoding/json"
	"errors"
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"github.com/gynshu-one/goph-keeper/server/api/utils"
	"github.com/gynshu-one/goph-keeper/server/storage"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers interface {
	CreateUser(w http.ResponseWriter, r *http.Request)

	LoginUser(w http.ResponseWriter, r *http.Request)

	LogoutUser(w http.ResponseWriter, r *http.Request)

	SetUserData(w http.ResponseWriter, r *http.Request)

	GetUserData(w http.ResponseWriter, r *http.Request)

	DeleteUserData(w http.ResponseWriter, r *http.Request)

	SyncUserData(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db      *mongo.Database
	storage storage.Storage
}

func NewHandlers(db *mongo.Database, storage storage.Storage) *handler {
	return &handler{
		db:      db,
		storage: storage,
	}
}

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

// SetUserData sets (creates or updates)
// Data should be encrypted and passed as json
// {		"type: "login",
//
//		"data: {
//		"type": "login",
//		"owner_id": "email",
//		"name": "login item name",
//		"info": "login item info",
//		"username": "username",
//		"password": "password",
//		"one_time_origin": "one time origin if exists"
//	}
func (h *handler) SetUserData(w http.ResponseWriter, r *http.Request) {
	session, err := FindSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tp := models.UserDataType(r.FormValue("type"))

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	// Refuse to receive unknown type
	decoder.DisallowUnknownFields()
	var data models.UserData

	switch tp {
	case models.LoginType:
		var login models.Login
		err = decoder.Decode(&login)
		if err != nil {
			http.Error(w, ErrInvalidType.Error(), http.StatusBadRequest)
			return
		}
		data = &login
	case models.ArbitraryTextType:
		var text models.ArbitraryText
		err = decoder.Decode(&text)
		if err != nil {
			http.Error(w, ErrInvalidType.Error(), http.StatusBadRequest)
			return
		}
		data = &text
	case models.BankCardType:
		var bankCard models.BankCard
		err = decoder.Decode(&bankCard)
		if err != nil {
			http.Error(w, ErrInvalidType.Error(), http.StatusBadRequest)
			return
		}
		data = &bankCard
	case models.BinaryType:
		var binary models.Binary
		err = decoder.Decode(&binary)
		if err != nil {
			http.Error(w, ErrInvalidType.Error(), http.StatusBadRequest)
			return
		}
	}
	if data == nil {
		http.Error(w, "invalid data type", http.StatusBadRequest)
		return
	}

	sessUserID := session.GetUserID()
	data.GetOwnerID(&sessUserID)

	err = h.storage.Set(r.Context(), tp, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
	return
}

// GetUserData gets a user data
// user must pass id as "id" url parameter (GET request)
func (h *handler) GetUserData(w http.ResponseWriter, r *http.Request) {
	id := models.UserDataID(r.FormValue("id"))
	if id == "" {
		http.Error(w, "id is empty", http.StatusBadRequest)
	}
	userDataType := models.UserDataType(r.FormValue("type"))
	if userDataType == "" {
		http.Error(w, "type is empty", http.StatusBadRequest)
	}

	session, err := FindSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.storage.Get(r.Context(), userDataType, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data == nil {
		http.Error(w, ErrNoDataFound.Error(), http.StatusNotFound)
		return
	}

	if data.GetOwnerID(nil) != session.GetUserID() {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	encData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(encData)
}

// DeleteUserData deletes a user data
func (h *handler) DeleteUserData(w http.ResponseWriter, r *http.Request) {
	session, err := FindSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := session.GetUserID()

	id := models.UserDataID(r.FormValue("id"))
	if id == "" {
		http.Error(w, "id is empty", http.StatusBadRequest)
	}
	userDataType := models.UserDataType(r.FormValue("type"))
	if userDataType == "" {
		http.Error(w, "type is empty", http.StatusBadRequest)
	}

	data, err := h.storage.Get(r.Context(), userDataType, id)
	if err != nil {
		if errors.Is(err, storage.ErrObjectMiss) {
			http.Error(w, ErrNothingToDelete.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data.GetOwnerID(nil) != userID {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	err = h.storage.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}

// SyncUserData syncs the data for a user
// server will return all data client have on server
// mapped by type slice of structs e.g. map[string][]models.UserData
func (h *handler) SyncUserData(w http.ResponseWriter, r *http.Request) {
	session, err := FindSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := session.GetUserID()

	allData, err := h.storage.GetUserData(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(allData) == 0 {
		http.Error(w, ErrNoDataFound.Error(), http.StatusNotFound)
	}

	marshalledData, err := json.Marshal(utils.PackData(allData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(marshalledData)
}
