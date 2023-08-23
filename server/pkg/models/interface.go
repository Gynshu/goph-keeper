package models

type UserData interface {
	EncryptAll(passphrase string) error
	DecryptAll(passphrase string) error

	GetOwnerID() string
	GetDataID() string

	MakeID()
	SetUpdatedAt()
	SetDeletedAt()
	SetCreatedAt()

	GetType() string
}

const (
	TextType     = "text"
	BankCardType = "bank_card"
	BinaryType   = "binary"
	LoginType    = "login"
)
