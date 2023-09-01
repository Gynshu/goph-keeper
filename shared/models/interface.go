package models

type UserDataModel interface {
	EncryptAll(passphrase string) error
	DecryptAll(passphrase string) error

	// GetOwnerID Returns the ownerID of the data
	// if id is not nil !!, it will be set to the ownerID
	GetOwnerID(id *string) string
	GetDataID() UserDataID

	MakeID()
	SetUpdatedAt()
	SetCreatedAt()

	GetType() UserDataType
}

type UserDataType string
type UserDataID string

type PackedUserData map[UserDataType][]UserDataModel

var UserDataTypes = []UserDataType{
	ArbitraryTextType,
	BankCardType,
	BinaryType,
	LoginType,
}

const (
	ArbitraryTextType = UserDataType("text")
	BankCardType      = UserDataType("bank_card")
	BinaryType        = UserDataType("binary")
	LoginType         = UserDataType("login")
)
