package models

type UserDataModel interface {
	EncryptAll(passphrase string) error
	DecryptAll(passphrase string) error

	// GetOrSetOwnerID Returns the ownerID of the data
	// if id is not nil !!, it will be set to the ownerID
	GetOrSetOwnerID(id *string) string
	GetDataID() UserDataID

	MakeID()
	SetUpdatedAt()
	SetCreatedAt()

	GetUpdatedAt() int64

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
