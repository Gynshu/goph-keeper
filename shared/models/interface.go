package models

import "encoding/json"

type UserDataModel interface {
	GetName() string
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

type SyncRequest struct {
	ToDelete []UserDataID   `json:"to_delete"`
	ToUpdate PackedUserData `json:"to_update"`
}

func NewSyncRequest() SyncRequest {
	return SyncRequest{
		ToUpdate: make(PackedUserData),
	}
}

func (p PackedUserData) UnmarshalJSON(data []byte) error {
	var m map[string][]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	for k, values := range m {
		switch k {
		case string(ArbitraryTextType):
			for _, value := range values {
				var text ArbitraryText
				text.ID = value.(map[string]interface{})["id"].(string)
				text.OwnerID = value.(map[string]interface{})["owner_id"].(string)
				text.Name = value.(map[string]interface{})["name"].(string)
				text.Text = value.(map[string]interface{})["text"].(string)
				text.CreatedAt = int64(value.(map[string]interface{})["created_at"].(float64))
				text.UpdatedAt = int64(value.(map[string]interface{})["updated_at"].(float64))
				p[ArbitraryTextType] = append(p[ArbitraryTextType], &text)
			}
		case string(BankCardType):
			for _, value := range values {
				var cards BankCard
				cards.ID = value.(map[string]interface{})["id"].(string)
				cards.OwnerID = value.(map[string]interface{})["owner_id"].(string)
				cards.Name = value.(map[string]interface{})["name"].(string)
				cards.Info = value.(map[string]interface{})["info"].(string)
				cards.CardType = CardType(value.(map[string]interface{})["card_type"].(string))
				cards.CardNum = value.(map[string]interface{})["card_num"].(string)
				cards.CardName = value.(map[string]interface{})["card_name"].(string)
				cards.CardCvv = value.(map[string]interface{})["card_cvv"].(string)
				cards.CardExp = value.(map[string]interface{})["card_exp"].(string)
				cards.CreatedAt = int64(value.(map[string]interface{})["created_at"].(float64))
				cards.UpdatedAt = int64(value.(map[string]interface{})["updated_at"].(float64))
				p[BankCardType] = append(p[BankCardType], &cards)
			}
		case string(BinaryType):
			for _, value := range values {
				var binary Binary
				binary.ID = value.(map[string]interface{})["id"].(string)
				binary.OwnerID = value.(map[string]interface{})["owner_id"].(string)
				binary.Name = value.(map[string]interface{})["name"].(string)
				binary.Info = value.(map[string]interface{})["info"].(string)
				binary.Binary = value.(map[string]interface{})["binary"].([]byte)
				binary.CreatedAt = int64(value.(map[string]interface{})["created_at"].(float64))
				binary.UpdatedAt = int64(value.(map[string]interface{})["updated_at"].(float64))
				p[BinaryType] = append(p[BinaryType], &binary)
			}
		case string(LoginType):
			for _, value := range values {
				var login Login
				login.ID = value.(map[string]interface{})["id"].(string)
				login.OwnerID = value.(map[string]interface{})["owner_id"].(string)
				login.Name = value.(map[string]interface{})["name"].(string)
				login.Info = value.(map[string]interface{})["info"].(string)
				login.Username = value.(map[string]interface{})["username"].(string)
				login.Password = value.(map[string]interface{})["password"].(string)
				login.OneTimeOrigin = value.(map[string]interface{})["one_time_origin"].(string)
				login.CreatedAt = int64(value.(map[string]interface{})["created_at"].(float64))
				login.UpdatedAt = int64(value.(map[string]interface{})["updated_at"].(float64))
				p[LoginType] = append(p[LoginType], &login)
			}
		}
	}
	return nil
}

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
