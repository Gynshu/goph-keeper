package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var instance *config
var once = sync.Once{}

type config struct {
	// Server is the server configuration
	MongoURI       string `json:"MONGO_URI"`
	HttpServerPort string `json:"HTTP_SERVER_PORT"`
	CertFilePath   string `json:"CERT_FILE_PATH"`
	KeyFilePath    string `json:"KEY_FILE_PATH"`
}

// NewConfig creates a new configuration struct
func newConfig() error {
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "goph-keeper") {
		wd = filepath.Dir(wd)
	}
	configPath := wd + "/server/config.json"
	if configPath == "" {
		log.Fatal().Msg("Config path is empty")
	}

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

	// decode json into struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(instance)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decode config")
	}

	folder := path.Dir(configPath)

	instance.CertFilePath = folder + "/cert/" + instance.CertFilePath
	instance.KeyFilePath = folder + "/cert/" + instance.KeyFilePath
	return nil
}

// GetConfig returns the configuration initialized by newConfig
func GetConfig() *config {
	once.Do(func() {
		err := newConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get config")
		}
	})
	return instance
}
