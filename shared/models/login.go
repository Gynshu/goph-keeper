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
	// DeletedAt is the time when this login was deleted
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at"`
}

// EncryptAll encrypts all sensitive data
func (l *Login) EncryptAll(passphrase string) error {
	encryptedPassword, err := utils.EncryptData([]byte(l.Password), passphrase)
	if err != nil {
		return err
	}
	l.Password = string(encryptedPassword)

	encryptedOneTimeOrigin, err := utils.EncryptData([]byte(l.OneTimeOrigin), passphrase)
	if err != nil {
		return err
	}
	l.OneTimeOrigin = string(encryptedOneTimeOrigin)
	encryptedInfo, err := utils.EncryptData([]byte(l.Info), passphrase)
	if err != nil {
		return err
	}
	l.Info = string(encryptedInfo)
	encryptedRecoveryCodes, err := utils.EncryptData([]byte(l.RecoveryCodes), passphrase)
	if err != nil {
		return err
	}
	l.RecoveryCodes = string(encryptedRecoveryCodes)
	return nil
}

// DecryptAll decrypts all the sensitive data
func (l *Login) DecryptAll(passphrase string) error {
	decryptedPassword, err := utils.DecryptData([]byte(l.Password), passphrase)
	if err != nil {
		return err
	}
	l.Password = string(decryptedPassword)

	decryptedOneTimeOrigin, err := utils.DecryptData([]byte(l.OneTimeOrigin), passphrase)
	if err != nil {
		return err
	}
	l.OneTimeOrigin = string(decryptedOneTimeOrigin)
	decryptedInfo, err := utils.DecryptData([]byte(l.Info), passphrase)
	if err != nil {
		return err
	}
	l.Info = string(decryptedInfo)
	decryptedRecoveryCodes, err := utils.DecryptData([]byte(l.RecoveryCodes), passphrase)
	if err != nil {
		return err
	}
	l.RecoveryCodes = string(decryptedRecoveryCodes)

	return nil
}

// GetOwnerID returns the owner id
func (l *Login) GetOwnerID() string {
	return l.OwnerID
}

// RegisterOneTime registers a new one-time password
func (l *Login) RegisterOneTime(secret string) (oneTime string, genTime time.Time, err error) {
	// Replace "your-secret-key" with your actual secret key
	secretKey := []byte(secret)

	// Generate a new TOTP configuration
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "goph-keeper",
		AccountName: l.Username,
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
	l.OneTimeOrigin = key.Secret()
	l.UpdatedAt = time.Now().Unix()
	return l.GenerateOneTimePassword()
}

// GenerateOneTimePassword generates a new one-time password
func (l *Login) GenerateOneTimePassword() (oneTime string, genTime time.Time, err error) {
	start := time.Now()
	key, err := totp.GenerateCode(l.OneTimeOrigin, start)
	if err != nil {
		log.Err(err).Msg("Error generating OTP")
		return "", time.Time{}, err
	}
	return key, start, nil
}

// RegisterRecoveryCodes registers a new recovery code
func (l *Login) RegisterRecoveryCodes(recoveryCodes string) {
	l.RecoveryCodes = recoveryCodes
	l.UpdatedAt = time.Now().Unix()
}

// GetDataID  returns the data id
func (l *Login) GetDataID() UserDataID {
	return UserDataID(l.ID)
}

// SetCreatedAt sets the created at time
func (l *Login) SetCreatedAt() {
	l.CreatedAt = time.Now().Unix()
}

// SetUpdatedAt sets the updated at time
func (l *Login) SetUpdatedAt() {
	l.UpdatedAt = time.Now().Unix()
}

// SetDeletedAt sets the deleted at time
func (l *Login) SetDeletedAt() {
	l.DeletedAt = time.Now().Unix()
}

// MakeID generates a new id
func (l *Login) MakeID() {
	l.ID = uuid.New().String()
}

// GetType returns the type of the data
func (l *Login) GetType() UserDataType {
	return LoginType
}
