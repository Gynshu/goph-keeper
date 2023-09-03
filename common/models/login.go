package models

import (
	"encoding/json"
	"github.com/gynshu-one/goph-keeper/common/utils"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog/log"
)

// Login is the model for a login
// All changes should be done through methods to ensure data consistency and update time
type Login struct {
	// Name is the name of the login
	Name string `json:"name" bson:"name"`
	// Info is the additional info about the login
	Info string `json:"info" bson:"info"`
	// Username is the username
	Username string `json:"username" bson:"username"`
	// Password is the password
	Password string `json:"password" bson:"password"`
	// OneTimeOrigin is the origin of the one-time password
	OneTimeOrigin string `json:"one_time_origin" bson:"one_time_origin"`
	// RecoveryCodes is the recovery code
	RecoveryCodes string `json:"recovery_codes" bson:"recovery_codes"`
}

// EncryptAll encrypts all sensitive data
func (data *Login) EncryptAll(passphrase string) (encryptedData []byte, err error) {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return utils.EncryptData(marshaled, passphrase)
}

// DecryptAll decrypts all sensitive data
func (data *Login) DecryptAll(passphrase string, encrypteData []byte) error {
	decrypted, err := utils.DecryptData(encrypteData, passphrase)
	if err != nil {
		return err
	}
	return json.Unmarshal(decrypted, data)
}

// RegisterOneTime registers a new one-time password
func (data *Login) RegisterOneTime(secret string) (oneTime string, genTime time.Time, err error) {
	// Replace "your-secret-key" with your actual secret key
	secretKey := []byte(secret)

	// Generate a new TOTP configuration
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "goph-keeper",
		AccountName: data.Username,
		Secret:      secretKey,
	})
	if err != nil {
		log.Err(err).Msg("Error generating OTP")
		return "", time.Time{}, err
	}
	valid := totp.Validate(key.Digits().String(), secret)
	if !valid {
		log.Err(err).Msg("Error validating OTP")
		return "", time.Time{}, err
	}
	data.OneTimeOrigin = key.Secret()
	return data.GenerateOneTimePassword()
}

// GenerateOneTimePassword generates a new one-time password
func (data *Login) GenerateOneTimePassword() (oneTime string, genTime time.Time, err error) {
	start := time.Now()
	key, err := totp.GenerateCode(data.OneTimeOrigin, start)
	if err != nil {
		log.Err(err).Msg("Error generating OTP")
		return "", time.Time{}, err
	}
	return key, start, nil
}

// RegisterRecoveryCodes registers a new recovery code
func (data *Login) RegisterRecoveryCodes(recoveryCodes string) {
	data.RecoveryCodes = recoveryCodes
}
