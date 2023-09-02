package sync

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"github.com/zalando/go-keyring"
	"net/http"
	"os"
)

const (
	RegisterEndpoint = "/user/create"
	LoginEndpoint    = "/user/login"
	SyncEndpoint     = "/user/sync"
)

type Mediator interface {
	Sync(ctx context.Context)
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) error
}

type mediator struct {
	client  *resty.Client
	storage storage.Storage
}

func NewMediator(storage storage.Storage) Mediator {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	md := &mediator{
		client:  resty.NewWithClient(&http.Client{Transport: tr}),
		storage: storage,
	}
	return md
}
func (m *mediator) SignUp(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username or password is empty")
	}
	get, err := m.client.NewRequest().SetContext(ctx).
		SetQueryParam("email", username).
		SetQueryParam("password", password).Get("https://" + config.GetConfig().ServerIP + RegisterEndpoint)
	if err != nil {
		return err
	}
	if get.StatusCode() != 200 {
		return fmt.Errorf("failed to register, status code: %d, and error %s", get.StatusCode(), get.Body())
	}
	// read cookie from response
	cookie := get.Cookies()
	if len(cookie) == 0 {
		return fmt.Errorf("failed to get cookie")
	}
	for _, c := range cookie {
		if c.Name == "session_id" {
			if c.Value == "" {
				return fmt.Errorf("failed to get cookie")
			}
			config.CurrentUser.Username = username
			config.CurrentUser.SessionID = c.Value
			m.createUserSessionFiles()
			m.Sync(ctx)
			return nil
		}
	}
	return fmt.Errorf("failed to get cookie")
}

func (m *mediator) SignIn(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username or password is empty")
	}
	get, err := m.client.NewRequest().SetContext(ctx).
		SetQueryParam("email", username).
		SetQueryParam("password", password).Get("https://" + config.GetConfig().ServerIP + LoginEndpoint)
	if err != nil {
		return err
	}
	if get.StatusCode() != 200 {
		return fmt.Errorf("failed to login, status code: %d and response %s", get.StatusCode(), get.Body())
	}
	// read cookie from response
	cookie := get.Cookies()
	if len(cookie) == 0 {
		return fmt.Errorf("failed to get cookie")
	}
	for _, c := range cookie {
		if c.Name == "session_id" {
			config.CurrentUser.Username = username
			config.CurrentUser.SessionID = c.Value
			m.createUserSessionFiles()
			m.Sync(ctx)
			return nil
		}
	}
	return nil
}

func (m *mediator) Sync(ctx context.Context) {
	var req models.SyncRequest
	req.ToDelete = []models.UserDataID{}
	req.ToUpdate = m.storage.Get()
	marshaledRequest, err := json.Marshal(req)
	if err != nil {
		config.ErrChan <- err
		return
	}
	// send data to server
	response, err := m.client.NewRequest().SetContext(ctx).
		SetBody(marshaledRequest).SetCookie(&http.Cookie{
		Name:  "session_id",
		Value: config.CurrentUser.SessionID,
	}).Post("https://" + config.GetConfig().ServerIP + SyncEndpoint)
	if err != nil {
		config.ErrChan <- err
		return
	}
	if response.StatusCode() == http.StatusUnauthorized {
		pass, err_ := keyring.Get(config.ServiceName, config.CurrentUser.Username)
		if err_ != nil {
			return
		}
		err_ = m.SignIn(ctx, config.CurrentUser.Username, pass)
		if err_ != nil {
			return
		}
	}

	if response.StatusCode() == http.StatusNoContent {
		return
	}

	var serverData = make(models.PackedUserData)
	// check fi body is empty or not
	if err = json.Unmarshal(response.Body(), &serverData); err != nil {
		if err.Error() == "EOF" {
			config.ErrChan <- err
			serverData = nil
		} else {
			config.ErrChan <- err
		}
	}
	err = m.storage.Put(serverData)
	if err != nil {
		config.ErrChan <- err
		return
	}
}

func (m *mediator) createUserSessionFiles() {
	// create session_id file
	file, err := os.Create(config.TempDir + "/" + config.SessionFile)
	if err != nil {
		config.ErrChan <- err
	}
	defer file.Close()

	// write session_id to file
	_, err = file.WriteString(config.CurrentUser.SessionID + "\n" + config.CurrentUser.Username)
	if err != nil {
		config.ErrChan <- err
	}
}
