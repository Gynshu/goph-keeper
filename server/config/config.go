package config

import (
	"github.com/halorium/env"
	"github.com/rs/zerolog/log"
	"sync"
)

var instance *config
var once = sync.Once{}

type config struct {
	// Server is the server configuration
	MongoURI       string `env:"MONGO_URI"`
	HttpServerPort string `env:"HTTP_SERVER_PORT"`
	CertFilePath   string `env:"CERT_FILE_PATH"`
	KeyFilePath    string `env:"KEY_FILE_PATH"`
}

// NewConfig creates a new configuration struct
func newConfig() error {
	instance = &config{}

	err := env.Unmarshal(instance)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
	}
	return nil
}

func GetConfig() *config {
	once.Do(func() {
		err := newConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get config")
		}
	})
	return instance
}
