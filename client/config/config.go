package config

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
	"time"
)

var instance *config
var TempDir = ""
var CurrentUser = struct {
	Username  string
	SessionID string
}{
	Username:  "",
	SessionID: "",
}

var ErrChan = make(chan error)

const (
	ServiceName = "goph-keeper"
	cacheFolder = "goph-keeper/cache"
	CfgFile     = "goph-keeper/config.json"
	SessionFile = "goph-keeper/session_id"
)

func init() {
	TempDir = os.TempDir()
	// create if not exists folder user-data
	_, err := os.Stat(TempDir + "/" + ServiceName)
	if os.IsNotExist(err) {
		// only owner can read and write data
		err = os.Mkdir(TempDir+"/"+ServiceName, 744)
		if err != nil {
			log.Fatal().Msg("Failed to create folder user-data")
		}
	}

	_, err = os.Stat(TempDir + "/" + cacheFolder)
	if os.IsNotExist(err) {
		err = os.Mkdir(TempDir+"/"+cacheFolder, 744)
		if err != nil {
			log.Fatal().Msg("Failed to create folder cache")
		}
	}

	// check if config file exists
	_, err = os.Stat(TempDir + "/" + CfgFile)
	if os.IsNotExist(err) {
		// create config file
		_, err = os.Create(TempDir + "/" + CfgFile)
		if err != nil {
			log.Fatal().Msg("Failed to create config file")
		}
	}

	// check if session_id file exists
	_, err = os.Stat(TempDir + "/" + SessionFile)
	if err == nil {
		// open read text
		file, err_ := os.Open(TempDir + "/" + SessionFile)
		if err_ != nil {
			log.Fatal().Msg("Failed to open session_id file")
		}
		defer file.Close()
		text, err_ := io.ReadAll(file)
		if err_ != nil {
			log.Fatal().Msg("Failed to read session_id file")
		}
		if len(text) == 0 {
			log.Info().Msg("Session id file exists but is empty")
		} else {
			CurrentUser.SessionID = strings.Split(string(text), "\n")[0]
			log.Info().Msg("Session id file exists and is not empty")
			CurrentUser.Username = strings.Split(string(text), "\n")[1]
		}
	} else if os.IsNotExist(err) {
		log.Info().Msg("Session id file doesn't exist")
	} else {
		log.Fatal().Msg("Failed to check session_id file")
	}
}

type config struct {
	// Server is the server configuration
	ServerIP    string        `json:"SERVER_IP" envDefault:"localhost:8080"`
	PollTimer   time.Duration `json:"POLL_TIMER" envDefault:"5s"`
	DumpTimer   time.Duration `json:"DUMP_TIMER" envDefault:"10s"`
	CacheFolder string        `json:"CACHE_FOLDER" envDefault:"goph-keeper-cache"`
}

// NewConfig creates a new configuration struct
func NewConfig(path string) error {
	instance = &config{}

	// open json file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open config file")
	}
	defer file.Close()

	// decode json into struct
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(instance)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decode config")
	}
	if instance.ServerIP == "" {
		return errors.New("server ip is empty")
	}
	if instance.PollTimer == 0 {
		return errors.New("poll timer is empty")
	}
	if instance.DumpTimer == 0 {
		return errors.New("dump timer is empty")
	}

	return nil
}

func GetConfig() *config {
	return instance
}
