package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"path"
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
	instance = &config{}

	gp, err := os.Getwd()
	if err != nil {
		log.Fatal().Msg("GOPATH is not set")
	}

	// find "goph-keeper" directory
	for path.Base(gp) != "goph-keeper" {
		gp = path.Dir(gp)
		if gp == "/" {
			log.Fatal().Msg("Failed to find goph-keeper directory")
		}
	}

	// open json file
	file, err := os.Open(gp + "/server/config.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open config file")
	}
	defer file.Close()

	// decode json into struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(instance)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decode config")
	}

	instance.CertFilePath = gp + "/" + instance.CertFilePath
	instance.KeyFilePath = gp + "/" + instance.KeyFilePath
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
