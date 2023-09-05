package config

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var instance *config
var TempDir = os.TempDir()
var CurrentUser = struct {
	Username  string
	SessionID string
}{
	Username:  "",
	SessionID: "",
}

const (
	ServiceName = "goph-keeper"
	SessionFile = "goph-keeper/session_id"
)

type config struct {
	// Server is the server configuration
	ServerIP  string        `json:"SERVER_IP" envDefault:"localhost:8080"`
	PollTimer time.Duration `json:"POLL_TIMER" envDefault:"5s"`
	DumpTimer time.Duration `json:"DUMP_TIMER" envDefault:"10s"`
}

// NewConfig creates a new configuration struct
func NewConfig() error {
	// find source directory
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "goph-keeper") {
		wd = filepath.Dir(wd)
	}

	configPath := wd + "/client/config.json"
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open config file")
	}
	instance = &config{}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close config file")
		}
	}()
	instance = &config{}
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

	// check if session_id file exists
	// if it does, read it and set CurrentUser
	_, err = os.Stat(TempDir + "/" + SessionFile)
	if err == nil {
		// open read text
		f, err_ := os.Open(TempDir + "/" + SessionFile)
		if err_ != nil {
			log.Fatal().Msg("Failed to open session_id file")
		}
		defer func(f *os.File) {
			err = f.Close()
			if err != nil {
				log.Fatal().Msg("Failed to close session_id file")
			}
		}(f)
		text, err_ := io.ReadAll(f)
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
	return nil
}

func GetConfig() *config {
	return instance
}
