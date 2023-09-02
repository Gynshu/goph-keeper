package handlers

import (
	"encoding/json"
	"github.com/gynshu-one/goph-keeper/server/api/utils"
	"github.com/gynshu-one/goph-keeper/server/storage"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
)

type Handlers interface {
	CreateUser(w http.ResponseWriter, r *http.Request)

	LoginUser(w http.ResponseWriter, r *http.Request)

	LogoutUser(w http.ResponseWriter, r *http.Request)

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

// SyncUserData syncs the data for a user
// If client didn't send any data, all data from server is returned
// All new data is added to the db all existing data is updated by the newest one
// If some data is missing from the client, it will be deleted from the db
// Data's sensitive fields should be encrypted and passed as json models.SyncRequest
func (h *handler) SyncUserData(w http.ResponseWriter, r *http.Request) {
	session, err := FindSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := session.GetUserID()

	defer r.Body.Close()

	var newSyncRequest = models.NewSyncRequest()

	rAll, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(rAll, &newSyncRequest)
	if err != nil {
		if err.Error() == "EOF" {

		} else {
			log.Debug().Err(err).Msg("failed to decode data")
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	sessUserID := session.GetUserID()

	if newSyncRequest.ToUpdate != nil {
		for _, slice := range newSyncRequest.ToUpdate {
			for _, item := range slice {
				err = h.storage.SetData(r.Context(), item)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// If client sent data, we need to check if it's newer than the one we have
	if len(newSyncRequest.ToDelete) != 0 {
		// Get all data from db
		var storedData models.PackedUserData
		storedData, err = h.storage.GetData(r.Context(), sessUserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, slice := range storedData {
			for _, item := range slice {
				if !utils.Contains(newSyncRequest.ToDelete, item.GetDataID()) {
					err = h.storage.Delete(r.Context(), item.GetDataID())
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
		}
	}

	storedData, err := h.storage.GetData(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(storedData) == 0 {
		http.Error(w, ErrNoDataFound.Error(), http.StatusNoContent)
	}

	marshalledData, err := json.Marshal(storedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(marshalledData)
}
