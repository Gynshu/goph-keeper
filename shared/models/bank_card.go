package models

import (
	"github.com/gynshu-one/goph-keeper/shared/utils"
	"time"

	"github.com/google/uuid"
)

type CardType string

const (
	Visa       CardType = "Visa"
	MasterCard CardType = "MasterCard"
	Amex       CardType = "American Express"
	Discover   CardType = "Discover"
)

type BankCard struct {
	// ID is the primary key
	ID string `json:"id" bson:"_id"`
	// OwnerID is the user who owns this text
	OwnerID string `json:"owner_id" bson:"ownerID"`
	// Name is the name of the bank
	Name string `json:"name" bson:"name"`
	// Info is the additional info about the card
	Info string `json:"info" bson:"info"`
	// CardType is the type of card such as Visa, MasterCard, etc.
	CardType CardType `json:"card_type" bson:"cardType"`
	// CardNum is the card number
	CardNum string `json:"card_num" bson:"cardNum"`
	// CardName is the name on the card
	CardName string `json:"card_name" bson:"cardName"`
	// CardCvv is the card's CVV
	CardCvv string `json:"card_cvv" bson:"cardCvv"`
	// CardExp is the card's expiration date
	CardExp string `json:"card_exp" bson:"cardExp"`
	// CreatedAt is the time when this text was created
	CreatedAt int64 `json:"created_at" bson:"createdAt"`
	// UpdatedAt is the time when this text was last updated
	UpdatedAt int64 `json:"updated_at" bson:"updatedAt"`
}

func (data *BankCard) GetName() string {
	return data.Name
}
func (data *BankCard) EncryptAll(passphrase string) error {
	encryptedCardNum, err := utils.EncryptData([]byte(string(data.CardNum)), passphrase)
	if err != nil {
		return err
	}
	data.CardNum = string(encryptedCardNum)

	encryptedCardName, err := utils.EncryptData([]byte(data.CardName), passphrase)
	if err != nil {
		return err
	}
	data.CardName = string(encryptedCardName)

	encryptedInfo, err := utils.EncryptData([]byte(data.Info), passphrase)
	if err != nil {
		return err
	}
	data.Info = string(encryptedInfo)

	encryptedCardCvv, err := utils.EncryptData([]byte(data.CardCvv), passphrase)
	if err != nil {
		return err
	}
	data.CardCvv = string(encryptedCardCvv)

	encryptedCardExp, err := utils.EncryptData([]byte(data.CardExp), passphrase)
	if err != nil {
		return err
	}
	data.CardExp = string(encryptedCardExp)

	data.UpdatedAt = time.Now().Unix()
	return nil
}

func (data *BankCard) DecryptAll(passphrase string) error {
	decryptedCardNum, err := utils.DecryptData([]byte(data.CardNum), passphrase)
	if err != nil {
		return err
	}
	data.CardNum = string(decryptedCardNum)

	decryptedCardName, err := utils.DecryptData([]byte(data.CardName), passphrase)
	if err != nil {
		return err
	}
	data.CardName = string(decryptedCardName)

	decryptedInfo, err := utils.DecryptData([]byte(data.Info), passphrase)
	if err != nil {
		return err
	}
	data.Info = string(decryptedInfo)

	decryptedCardCvv, err := utils.DecryptData([]byte(data.CardCvv), passphrase)
	if err != nil {
		return err
	}
	data.CardCvv = string(decryptedCardCvv)

	decryptedCardExp, err := utils.DecryptData([]byte(data.CardExp), passphrase)
	if err != nil {
		return err
	}
	data.CardExp = string(decryptedCardExp)

	return nil
}

func (data *BankCard) GetOrSetOwnerID(id *string) string {
	if id != nil {
		data.OwnerID = *id
	}
	return data.OwnerID
}

func (data *BankCard) GetDataID() UserDataID {
	return UserDataID(data.ID)
}

func (data *BankCard) SetCreatedAt() {
	data.CreatedAt = time.Now().Unix()
}

func (data *BankCard) SetUpdatedAt() {
	data.UpdatedAt = time.Now().Unix()
}

func (data *BankCard) GetUpdatedAt() int64 {
	return data.UpdatedAt
}

func (data *BankCard) MakeID() {
	data.ID = uuid.New().String()
}

func (data *BankCard) GetType() UserDataType {
	return BankCardType
}
