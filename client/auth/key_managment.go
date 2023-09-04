package auth

import (
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
)

// SetPass sets pass to local os keyring for simplicity
func SetPass(pass string) {
	err := keyring.Set(config.ServiceName, config.CurrentUser.Username, pass)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set pass")
	}
}

// SetSecret sets secret to local os keyring for simplicity
func SetSecret(secret string) {
	err := keyring.Set(config.ServiceName, config.CurrentUser.Username+"s", secret)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set secret")
	}
}

// GetPass gets pass from local os keyring for simplicity
func GetPass() string {
	pass, err := keyring.Get(config.ServiceName, config.CurrentUser.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get pass")
		return ""
	}
	return pass
}

// GetSecret gets secret from local os keyring for simplicity
func GetSecret() string {
	secret, err := keyring.Get(config.ServiceName, config.CurrentUser.Username+"s")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get secret")
		return ""
	}
	return secret
}
