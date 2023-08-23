package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Mediator interface {
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
func (m *mediator) Sync(ctx context.Context) {
	pollTimer := time.NewTimer(5 * time.Minute)
	dumpTimer := time.NewTimer(10 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			err := m.dumpToFile()
			if err != nil {
				log.Err(err).Msg("failed to dump data to file")
				return
			}
			return
		case <-dumpTimer.C:
			err := m.dumpToFile()
			if err != nil {
				log.Err(err).Msg("failed to dump data to file")
				return
			}
		case <-pollTimer.C:
		}
	}
}

func (m *mediator) dumpToFile() error {
	data := m.storage.GetData()

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
	_, err = os.Stat(thisDir + "/data-keeper")
	if os.IsNotExist(err) {
		err = os.Mkdir(thisDir+"/data-keeper", os.ModePerm)
		if err != nil {
			log.Trace().Msg("failed to create data-keeper folder")
			return err
		}
	}

	year, month, day := time.Now().Date()

	minute, hour, sec := time.Now().Clock()

	var fileName = fmt.Sprintf("%d-%d-%d-%d-%d-%d.json", year, month, day, hour, minute, sec)

	file, err := os.Create(thisDir + "/data-keeper/" + fileName)
	if err != nil {
		log.Trace().Msg("failed to create file")
		// if file exists, try to open it
		file, err = os.OpenFile(thisDir+"/data-keeper/"+fileName, os.O_APPEND|os.O_WRONLY, os.ModePerm)
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

	_, err = os.Stat(thisDir + "/data-keeper")
	if os.IsNotExist(err) {
		log.Trace().Msg("data-keeper folder doesn't exist")
	}
	files, err := os.ReadDir(thisDir + "/data-keeper")
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

		fileInfo, err := os.Stat(thisDir + "/data-keeper/" + file.Name())
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
	file, err := os.Open(thisDir + "/data-keeper/" + latestFile)
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

	err = m.storage.PutData(data)
	if err != nil {
		log.Trace().Msg("failed to put data")
		return err
	}
	return nil
}
