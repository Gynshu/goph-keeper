package models

// User is a struct for user of client app
type User struct {
	// Email is the email of the user
	Email string `json:"email" bson:"_id"`
	// EncryptedPassphrase is the encrypted passphrase of the user
	Passphrase string `json:"encrypted_passphrase" bson:"encrypted_passphrase"`
	// CreatedAt is the time when this user was created
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	// UpdatedAt is the time when this user was last updated
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	// DeletedAt is the time when this user was deleted
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at"`
}
