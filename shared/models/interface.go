package models

type UserData interface {
	EncryptAll(passphrase string) error
	DecryptAll(passphrase string) error

	GetOwnerID() string
	GetDataID() UserDataID

	MakeID()
	SetUpdatedAt()
	SetDeletedAt()
	SetCreatedAt()

	GetType() UserDataType
}

type UserDataType string

var UserDataTypes = []UserDataType{
	TextType,
	BankCardType,
	BinaryType,
	LoginType,
}

type UserDataID string

type PackedUserData map[UserDataType]map[UserDataID]UserData

const (
	TextType     = UserDataType("text")
	BankCardType = UserDataType("bank_card")
	BinaryType   = UserDataType("binary")
	LoginType    = UserDataType("login")
)
