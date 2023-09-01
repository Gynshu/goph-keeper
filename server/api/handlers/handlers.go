package handlers

import (
	"encoding/json"
	"github.com/gynshu-one/goph-keeper/server/api/utils"
	"github.com/gynshu-one/goph-keeper/server/storage"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
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
// Data's sensitive fields should be encrypted and passed as json
// {		"type: "login",
//
//		"data: {
//		"type": "login",
//		"owner_id": "email",
//		"name": *****",
//		"info": "*****",
//		"username": "*****",
//		"password": "*****",
//		"one_time_origin": "*****"
//		"created_at": unix timestamp,
//		"updated_at": unix timestamp
//	}
func (h *handler) SyncUserData(w http.ResponseWriter, r *http.Request) {
	session, err := FindSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := session.GetUserID()

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var data models.PackedUserData

	// check fi body is empty or not
	if err = decoder.Decode(&data); err != nil {
		if err.Error() == "EOF" {
			data = make(models.PackedUserData)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	sessUserID := session.GetUserID()

	// If client sent data, we need to check if it's newer than the one we have
	if len(data) != 0 {
		// Get all data from db
		var storedData map[models.UserDataID]models.UserDataModel
		storedData, err = h.storage.GetData(r.Context(), sessUserID)
		if err != nil {
			return
		}

		// Loop
		for _, v := range data {
			for _, item := range v {
				// If item is owned by the user
				if item.GetOrSetOwnerID(&sessUserID) == sessUserID {
					// If item is not in db, delete it
					if _, ok := storedData[item.GetDataID()]; !ok {
						err = h.storage.Delete(r.Context(), item.GetDataID())
						if err != nil {
							log.Err(err).Msg("failed to delete data")
						}
					} else {
						// If item is in db, check if it's newer
						if item.GetUpdatedAt() > storedData[item.GetDataID()].GetUpdatedAt() {
							err = h.storage.SetData(r.Context(), item.GetType(), item)
							if err != nil {
								log.Err(err).Msg("failed to set data")
							}
						}
					}
				}
			}
		}
	}

	allData, err := h.storage.GetData(r.Context(), userID)
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
