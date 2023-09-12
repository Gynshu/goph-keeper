// Package auth is about storing and retrieving passwords and secrets from OS keyring
// for simplicity. It also stores current user name in a file in user's home directory.
// Every time user starts the client, it reads the file and sets CurrentUser.Username.
// If the file does not exist, In this case user must login again.
package auth

import (
	"errors"
	"os"

	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
)

var userFile string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get user home dir")
		return
	}

	err = os.MkdirAll(homeDir+string(os.PathSeparator)+"."+config.ServiceName, 0700)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create .goph-keeper dir")
	}

	userFile = homeDir +
		string(os.PathSeparator) +
		"." + config.ServiceName +
		string(os.PathSeparator) + "user.txt"
	file, err := os.ReadFile(userFile)
	if errors.Is(err, os.ErrNotExist) {
		log.Info().Err(err).Msg("No user file found")
		return
	} else if err != nil {
		log.Err(err).Msg("Failed to read user file")
		return
	}
	CurrentUser.Username = string(file)
}

var CurrentUser = struct {
	Username  string
	SessionID string
}{
	Username:  "",
	SessionID: "",
}

// SetPass sets pass to local os keyring for simplicity
func SetPass(pass string) {
	err := keyring.Set(config.ServiceName, CurrentUser.Username, pass)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to set pass")
	}

	// We must know username to get pass and secret from keyring
	f, err := os.Create(userFile)
	if err != nil {
		log.Err(err).Msg("Failed to create user file")
	}
	_, err = f.WriteString(CurrentUser.Username)
	if err != nil {
		log.Err(err).Msg("Failed to write to user file")
	}
	err = f.Close()
	if err != nil {
		log.Err(err).Msg("Failed to close user file")
	}
}

// SetSecret sets secret to local os keyring for simplicity
func SetSecret(secret string) {
	err := keyring.Set(config.ServiceName, CurrentUser.Username+"s", secret)
	if err != nil {
		log.Err(err).Msg("Failed to set secret")
	}
}

// GetPass gets pass from local os keyring
func GetPass() string {
	pass, err := keyring.Get(config.ServiceName, CurrentUser.Username)
	if err != nil {
		log.Err(err).Msg("Failed to get pass")
		return ""
	}
	return pass
}

// GetSecret gets secret from local os keyring
func GetSecret() string {
	secret, err := keyring.Get(config.ServiceName, CurrentUser.Username+"s")
	if err != nil {
		log.Err(err).Msg("Failed to get secret")
		return ""
	}
	return secret
}
