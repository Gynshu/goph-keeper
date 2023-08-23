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
	ID string `json:"id" bson:"id"`
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
	// DeletedAt is the time when this text was deleted
	DeletedAt int64 `json:"deleted_at" bson:"deletedAt"`
}

func (b *BankCard) EncryptAll(passphrase string) error {
	encryptedCardNum, err := utils.EncryptData([]byte(string(b.CardNum)), passphrase)
	if err != nil {
		return err
	}
	b.CardNum = string(encryptedCardNum)

	encryptedCardName, err := utils.EncryptData([]byte(b.CardName), passphrase)
	if err != nil {
		return err
	}
	b.CardName = string(encryptedCardName)

	encryptedInfo, err := utils.EncryptData([]byte(b.Info), passphrase)
	if err != nil {
		return err
	}
	b.Info = string(encryptedInfo)

	encryptedCardCvv, err := utils.EncryptData([]byte(b.CardCvv), passphrase)
	if err != nil {
		return err
	}
	b.CardCvv = string(encryptedCardCvv)

	encryptedCardExp, err := utils.EncryptData([]byte(b.CardExp), passphrase)
	if err != nil {
		return err
	}
	b.CardExp = string(encryptedCardExp)

	b.UpdatedAt = time.Now().Unix()
	return nil
}

func (b *BankCard) DecryptAll(passphrase string) error {
	decryptedCardNum, err := utils.DecryptData([]byte(b.CardNum), passphrase)
	if err != nil {
		return err
	}
	b.CardNum = string(decryptedCardNum)

	decryptedCardName, err := utils.DecryptData([]byte(b.CardName), passphrase)
	if err != nil {
		return err
	}
	b.CardName = string(decryptedCardName)

	decryptedInfo, err := utils.DecryptData([]byte(b.Info), passphrase)
	if err != nil {
		return err
	}
	b.Info = string(decryptedInfo)

	decryptedCardCvv, err := utils.DecryptData([]byte(b.CardCvv), passphrase)
	if err != nil {
		return err
	}
	b.CardCvv = string(decryptedCardCvv)

	decryptedCardExp, err := utils.DecryptData([]byte(b.CardExp), passphrase)
	if err != nil {
		return err
	}
	b.CardExp = string(decryptedCardExp)

	return nil
}

func (b *BankCard) GetOwnerID() string {
	return b.OwnerID
}

func (b *BankCard) GetDataID() UserDataID {
	return UserDataID(b.ID)
}

func (b *BankCard) SetCreatedAt() {
	b.CreatedAt = time.Now().Unix()
}

func (b *BankCard) SetUpdatedAt() {
	b.UpdatedAt = time.Now().Unix()
}

func (b *BankCard) SetDeletedAt() {
	b.DeletedAt = time.Now().Unix()
}

func (b *BankCard) MakeID() {
	b.ID = uuid.New().String()
}

func (b *BankCard) GetType() UserDataType {
	return BankCardType
}
