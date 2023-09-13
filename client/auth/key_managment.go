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
	// Determine user home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get user home dir")
		return
	}

	// Create .goph-keeper dir in user home dir
	err = os.MkdirAll(homeDir+string(os.PathSeparator)+"."+config.ServiceName, 0700)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create .goph-keeper dir")
	}

	// Construct user file path
	userFile = homeDir +
		string(os.PathSeparator) +
		"." + config.ServiceName +
		string(os.PathSeparator) + "user.txt"
	file, err := os.ReadFile(userFile)

	// Something went wrong?
	if errors.Is(err, os.ErrNotExist) {
		log.Info().Err(err).Msg("No user file found")
		return
	} else if err != nil {
		log.Err(err).Msg("Failed to read user file")
		return
	}

	// Set current user name
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

	// Write username to user file. Next time user starts the client,
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
