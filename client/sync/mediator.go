package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

const (
	RegisterEndpoint = "/user/create"
	LoginEndpoint    = "/user/login"
	SetEndpoint      = "/user/set"
	GetEndpoint      = "/user/get"
	DeleteEndpoint   = "/user/delete"
	SyncEndpoint     = "/user/sync"
)

type Mediator interface {
	Sync(ctx context.Context)
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) error
}

type mediator struct {
	client  *resty.Client
	storage storage.Storage
}

func NewMediator(storage storage.Storage) Mediator {
	return &mediator{
		client:  resty.New(),
		storage: storage,
	}
}
func (m *mediator) SignUp(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username or password is empty")
	}
	get, err := m.client.NewRequest().SetContext(ctx).
		SetPathParam("email", username).
		SetPathParam("password", password).Get("https://" + config.GetConfig().ServerIP + RegisterEndpoint)
	if err != nil {
		return err
	}
	if get.StatusCode() != 200 {
		return fmt.Errorf("failed to register, status code: %d", get.StatusCode())
	}
	// read cookie from response
	cookie := get.Cookies()
	if len(cookie) == 0 {
		return fmt.Errorf("failed to get cookie")
	}
	for _, c := range cookie {
		if c.Name == "session_id" {
			if c.Value == "" {
				return fmt.Errorf("failed to get cookie")
			}
			config.CurrentUser.Username = username
			config.CurrentUser.SessionID = c.Value
			return nil
		}
	}
	return fmt.Errorf("failed to get cookie")
}

func (m *mediator) SignIn(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username or password is empty")
	}
	get, err := m.client.NewRequest().SetContext(ctx).
		SetPathParam("email", username).
		SetPathParam("password", password).Get("https://" + config.GetConfig().ServerIP + LoginEndpoint)
	if err != nil {
		return err
	}
	if get.StatusCode() != 200 {
		return fmt.Errorf("failed to login, status code: %d", get.StatusCode())
	}
	// read cookie from response
	cookie := get.Cookies()
	if len(cookie) == 0 {
		return fmt.Errorf("failed to get cookie")
	}
	for _, c := range cookie {
		if c.Name == "session_id" {
			config.CurrentUser.Username = username
			config.CurrentUser.SessionID = c.Value
			return nil
		}
	}
	return nil
}

func (m *mediator) Sync(ctx context.Context) {
	pollTimer := time.NewTimer(config.GetConfig().PollTimer)
	dumpTimer := time.NewTimer(config.GetConfig().DumpTimer)

	for {
		select {
		case <-ctx.Done():
			err := m.dumpToFile()
			if err != nil {
				log.Err(err).Msg("failed to dump data to file")
			}
			return
		case <-dumpTimer.C:
			err := m.dumpToFile()
			if err != nil {
				log.Err(err).Msg("failed to dump data to file")
			}
		case <-pollTimer.C:
		}
	}
}

func (m *mediator) dumpToFile() error {
	data := m.storage.Get()

	marshaled, err := json.Marshal(data)
	if err != nil {
		log.Trace().Msg("failed to marshal data")
		return err
	}

	thisDir, err := os.Getwd()
	if err != nil {
		log.Trace().Msg("failed to get working directory")
		return err
	}

	// create folder goph-keeper-data if it doesn't exist
	_, err = os.Stat(thisDir + config.GetConfig().CacheFolder)
	if os.IsNotExist(err) {
		err = os.Mkdir(thisDir+config.GetConfig().CacheFolder, os.ModePerm)
		if err != nil {
			log.Trace().Msg("failed to create data-keeper folder")
			return err
		}
	}

	year, month, day := time.Now().Date()

	minute, hour, sec := time.Now().Clock()

	var fileName = fmt.Sprintf("%d-%d-%d-%d-%d-%d.json", year, month, day, hour, minute, sec)

	file, err := os.Create(thisDir + config.GetConfig().CacheFolder + "/" + fileName)
	if err != nil {
		log.Trace().Msg("failed to create file")
		// if file exists, try to open it
		file, err = os.OpenFile(thisDir+config.GetConfig().CacheFolder+"/"+fileName, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Trace().Msg("failed to open file")
			return err
		}
		// if file is opened, write to it
		_, err = file.Write(marshaled)
		if err != nil {
			log.Trace().Msg("failed to write to file")
			return err
		}
		return err
	}
	return nil
}

func (m *mediator) loadFromFile() error {
	thisDir, err := os.Getwd()
	if err != nil {
		log.Trace().Msg("failed to get working directory")
		return err
	}

	_, err = os.Stat(thisDir + config.GetConfig().CacheFolder)
	if os.IsNotExist(err) {
		log.Trace().Msgf("%s folder doesn't exist", config.GetConfig().CacheFolder)
	}
	files, err := os.ReadDir(thisDir + config.GetConfig().CacheFolder)
	if err != nil {
		log.Trace().Msg("failed to read data-keeper folder")
		return err
	}
	var latestTime time.Time
	var latestFile string

	// loop files in data-keeper folder
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := os.Stat(thisDir + config.GetConfig().CacheFolder + "/" + file.Name())
		if err != nil {
			log.Trace().Msg("failed to get file info")
			continue
		}
		modTime := fileInfo.ModTime()
		if modTime.After(latestTime) {
			latestTime = modTime
			latestFile = file.Name()
		}
	}

	// open latest file
	file, err := os.Open(thisDir + config.GetConfig().CacheFolder + "/" + latestFile)
	if err != nil {
		log.Trace().Msg("failed to open latest file")
		return err
	}

	var data = make(models.PackedUserData)
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		log.Trace().Msg("failed to decode data")
		return err
	}

	err = m.storage.Put(data)
	if err != nil {
		log.Trace().Msg("failed to put data")
		return err
	}
	return nil
}
