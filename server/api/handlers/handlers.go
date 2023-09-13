package handlers

import (
	"encoding/json"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/gynshu-one/goph-keeper/server/storage"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

// Handlers is an interface for all handlers at once
type Handlers interface {
	CreateUser(w http.ResponseWriter, r *http.Request)

	LoginUser(w http.ResponseWriter, r *http.Request)

	LogoutUser(w http.ResponseWriter, r *http.Request)

	SyncUserData(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	storage storage.Storage
}

// NewHandlers creates a new handlers instance
func NewHandlers(storage storage.Storage) *handler {
	return &handler{
		storage: storage,
	}
}

// SyncUserData syncs the data for a user
// If client didn't send any data, all data from server is returned
// All new data is added to the db all existing data is updated by the newest one
// If some data is missing from the client, it will be deleted from the db
// Data's sensitive fields should be encrypted into binary
// data should be sent in the []models.DataWrapper format:
func (h *handler) SyncUserData(w http.ResponseWriter, r *http.Request) {
	// Fist we need to get user id from session
	session, err := FindSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := session.GetUserID()

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Err(err).Msg("failed to close body")
		}
	}(r.Body)

	var fromClient []models.DataWrapper

	// Decode data from client
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&fromClient)
	if err != nil {
		if err.Error() != "EOF" {
			log.Debug().Err(err).Msg("failed to decode data")
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			err = nil
		}
	}

	sessUserID := session.GetUserID()

	// Upsert data to db if needed
	for _, data := range fromClient {
		if data.OwnerID != sessUserID {
			log.Info().Msg("user tried to sync data that doesn't belong to him")
			continue
		}
		err = h.storage.SetData(r.Context(), data)
		if err != nil {
			log.Err(err).Msg("failed to upsert data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Get fresh data from db
	storedData, err := h.storage.GetData(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(storedData) == 0 {
		http.Error(w, ErrNoDataFound.Error(), http.StatusNoContent)
	}

	// Marshal and send
	marshalledData, err := json.Marshal(storedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshalledData)
	if err != nil {
		log.Err(err).Msg("failed to write response")
		return
	}
}
