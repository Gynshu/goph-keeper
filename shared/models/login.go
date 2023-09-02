package models

import (
	"github.com/gynshu-one/goph-keeper/shared/utils"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog/log"
)

// Login is the model for a login
// All changes should be done through methods to ensure data consistency and update time
type Login struct {
	// string is the primary key
	ID string `json:"id" bson:"_id"`
	// OwnerID is the user who owns this text
	OwnerID string `json:"owner_id" bson:"owner_id"`
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
	// CreatedAt is the time when this login was created
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	// UpdatedAt is the time when this login was last updated
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
}

func (data *Login) GetName() string {
	return data.Name
}

// EncryptAll encrypts all sensitive data
func (data *Login) EncryptAll(passphrase string) error {
	encryptedPassword, err := utils.EncryptData([]byte(data.Password), passphrase)
	if err != nil {
		return err
	}
	data.Password = string(encryptedPassword)

	encryptedOneTimeOrigin, err := utils.EncryptData([]byte(data.OneTimeOrigin), passphrase)
	if err != nil {
		return err
	}
	data.OneTimeOrigin = string(encryptedOneTimeOrigin)
	encryptedInfo, err := utils.EncryptData([]byte(data.Info), passphrase)
	if err != nil {
		return err
	}
	data.Info = string(encryptedInfo)
	encryptedRecoveryCodes, err := utils.EncryptData([]byte(data.RecoveryCodes), passphrase)
	if err != nil {
		return err
	}
	data.RecoveryCodes = string(encryptedRecoveryCodes)
	return nil
}

// DecryptAll decrypts all the sensitive data
func (data *Login) DecryptAll(passphrase string) error {
	decryptedPassword, err := utils.DecryptData([]byte(data.Password), passphrase)
	if err != nil {
		return err
	}
	data.Password = string(decryptedPassword)

	decryptedOneTimeOrigin, err := utils.DecryptData([]byte(data.OneTimeOrigin), passphrase)
	if err != nil {
		return err
	}
	data.OneTimeOrigin = string(decryptedOneTimeOrigin)
	decryptedInfo, err := utils.DecryptData([]byte(data.Info), passphrase)
	if err != nil {
		return err
	}
	data.Info = string(decryptedInfo)
	decryptedRecoveryCodes, err := utils.DecryptData([]byte(data.RecoveryCodes), passphrase)
	if err != nil {
		return err
	}
	data.RecoveryCodes = string(decryptedRecoveryCodes)

	return nil
}

// GetOwnerID returns the owner id
func (data *Login) GetOrSetOwnerID(id *string) string {
	if id != nil {
		data.OwnerID = *id
	}
	return data.OwnerID
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
	data.UpdatedAt = time.Now().Unix()
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
	data.UpdatedAt = time.Now().Unix()
}

// GetDataID  returns the data id
func (data *Login) GetDataID() UserDataID {
	return UserDataID(data.ID)
}

// SetCreatedAt sets the created at time
func (data *Login) SetCreatedAt() {
	data.CreatedAt = time.Now().Unix()
}

// SetUpdatedAt sets the updated at time
func (data *Login) SetUpdatedAt() {
	data.UpdatedAt = time.Now().Unix()
}

func (data *Login) GetUpdatedAt() int64 {
	return data.UpdatedAt
}

// MakeID generates a new id
func (data *Login) MakeID() {
	data.ID = uuid.New().String()
}

// GetType returns the type of the data
func (data *Login) GetType() UserDataType {
	return LoginType
}
