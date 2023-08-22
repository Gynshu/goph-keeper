package models

type UserData interface {
	EncryptAll(passphrase string) error
	DecryptAll(passphrase string) error
	GetOwnerID() string
}
