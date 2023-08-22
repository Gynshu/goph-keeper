package models

type User struct {
	// ID is the primary key
	ID string `json:"id" bson:"_id"`
	// Name is the name of the user
	Name string `json:"name" bson:"name"`
	// Email is the email of the user
	Email string `json:"email" bson:"email"`
	// EncryptedPassphrase is the encrypted passphrase of the user
	EncryptedPassphrase string `json:"encrypted_passphrase" bson:"encrypted_passphrase"`
	// CreatedAt is the time when this user was created
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	// UpdatedAt is the time when this user was last updated
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	// DeletedAt is the time when this user was deleted
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at"`
}
