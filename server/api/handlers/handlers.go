package handlers

import (
	"encoding/json"
	"errors"
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	"github.com/gynshu-one/goph-keeper/server/api/utils"
	"github.com/gynshu-one/goph-keeper/server/storage"
	"github.com/gynshu-one/goph-keeper/shared/models"
	utils2 "github.com/gynshu-one/goph-keeper/shared/utils"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers interface {
	// CreateUser creates a new user
	// hashes it and stores to database, returns an error if something went wrong
	CreateUser(w http.ResponseWriter, r *http.Request)

	// LoginUser logs in a user
	// returns a session ID and an error if something went wrong
	LoginUser(w http.ResponseWriter, r *http.Request)

	// LogoutUser logs out a user and deletes the session
	// returns an error if something went wrong
	LogoutUser(w http.ResponseWriter, r *http.Request)

	// SetUserData sets (creates or updates)
	// the data and its type must be passed to the request through
	// "data" and "type" parameters respectively  all data should be encrypted
	SetUserData(w http.ResponseWriter, r *http.Request)

	// GetUserData returns the data for a user
	// the data id must be passed to the request through "id" parameter
	GetUserData(w http.ResponseWriter, r *http.Request)

	// DeleteUserData deletes the data for a user
	DeleteUserData(w http.ResponseWriter, r *http.Request)

	// SyncUserData syncs the data for a user
	// server will return all data client have on server
	// mapped by type slice of structs e.g. map[string][]models.UserData
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
	_, err = auth.Sessions.CreateSession(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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

	err := h.db.Collection("users").FindOne(r.Context(), bson.M{"email": email}).Decode(&user)
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

	session, err := auth.Sessions.CreateSession(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(session.ID))
}

func (h *handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("sessionID")

	err := auth.Sessions.DeleteSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) SetUserData(w http.ResponseWriter, r *http.Request) {
	tp := models.UserDataType(r.FormValue("type"))
	if tp == "" {
		http.Error(w, "type is empty", http.StatusBadRequest)
	}
	inputData := r.FormValue("data")
	if inputData == "" {
		http.Error(w, "data is empty", http.StatusBadRequest)
	}

	var data models.UserData

	switch tp {
	case models.LoginType:
		var login models.Login
		err := json.Unmarshal([]byte(inputData), &login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		data = &login
	case models.BinaryType:
		var binary models.Binary
		err := json.Unmarshal([]byte(inputData), &binary)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		data = &binary
	case models.BankCardType:
		var bankCard models.BankCard
		err := json.Unmarshal([]byte(inputData), &bankCard)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		data = &bankCard
	case models.TextType:
		var text models.ArbitraryText
		err := json.Unmarshal([]byte(inputData), &text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		data = &text
	default:
		http.Error(w, "invalid type", http.StatusBadRequest)
	}

	err := h.storage.Set(r.Context(), data.GetDataID(), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
	return
}

func (h *handler) GetUserData(w http.ResponseWriter, r *http.Request) {
	id := models.UserDataID(r.FormValue("id"))
	if id == "" {
		http.Error(w, "id is empty", http.StatusBadRequest)
	}

	// no need to chec err here, bec of middleware
	session, _ := auth.Sessions.GetSession(r.Header.Get("sessionID"))

	data, err := h.storage.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data.GetOwnerID() != session.GetUserID() {
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

func (h *handler) DeleteUserData(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("sessionID")
	session, _ := auth.Sessions.GetSession(sessionID)

	userID := session.GetUserID()

	id := models.UserDataID(r.FormValue("id"))

	data, err := h.storage.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrObjectMiss) {
			http.Error(w, ErrNothingToDelete.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data.GetOwnerID() != userID {
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

func (h *handler) SyncUserData(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("sessionID")
	session, _ := auth.Sessions.GetSession(sessionID)

	userID := session.GetUserID()

	allData, err := h.storage.GetUserData(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(allData) == 0 {
		http.Error(w, ErrNoDataFound.Error(), http.StatusNotFound)
	}

	marshalledData, err := json.Marshal(utils2.PackData(allData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(marshalledData)
}
