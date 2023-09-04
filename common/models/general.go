package models

// BasicData is an interface for all data types
// It provides methods to encrypt and decrypt data
type BasicData interface {
	EncryptAll(passphrase string) (encryptedData []byte, err error)
	DecryptAll(passphrase string, encrypteData []byte) error
}

// DataWrapper is a struct that wraps BasicData and provides additional information about the data
// such as owner id, type, name, updated_at, created_at, deleted_at
// it makes easier to store data in the database that shouldn't know anything about the data
type DataWrapper struct {
	// ID is the unique identifier of the data
	ID string `json:"id" bson:"_id"`
	// OwnerID is the user who owns this data
	OwnerID string `json:"owner_id" bson:"owner_id"`
	// Type is the type of the data such as ArbitraryTextType, BankCardType, BinaryType, LoginType
	Type      string `json:"type" bson:"type"`
	Name      string `json:"name" bson:"name"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	DeletedAt int64  `json:"deleted_at" bson:"deleted_at"`
	Data      []byte `json:"data" bson:"data"`
}

const (
	ArbitraryTextType = "arbitrary_text"
	BankCardType      = "bank_card"
	BinaryType        = "binary"
	LoginType         = "login"
)
