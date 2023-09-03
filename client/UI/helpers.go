package UI

import (
	"context"
	"github.com/google/uuid"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/zalando/go-keyring"
	"time"
)

// encryptSave encrypts data and saves it to the storage
// by creating models.UserDataModel struct and adding it to the storage
func (u *ui) encryptSave(data models.DataModeler, wrapper models.UserDataModel) error {
	pass, err := keyring.Get(config.ServiceName, config.CurrentUser.Username)
	if err != nil {
		return err
	}
	encrypted, err := data.EncryptAll(pass)
	if err != nil {
		return err
	}
	t := time.Now().Unix()
	wrapper.Data = encrypted
	wrapper.CreatedAt = t
	wrapper.UpdatedAt = t
	wrapper.ID = uuid.NewString()
	wrapper.OwnerID = config.CurrentUser.Username
	wrapper.DeletedAt = 0

	if err = u.storage.Add(wrapper); err != nil {
		return err
	}
	u.mediator.Sync(context.Background())
	return nil
}
