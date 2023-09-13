package sync

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/gynshu-one/goph-keeper/client/auth"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/zalando/go-keyring"
)

const (
	RegisterEndpoint = "/user/create"
	LoginEndpoint    = "/user/login"
	Endpoint         = "/user/sync"
)

// Mediator is a mediator between client and server
// It is responsible for sending data to server and receiving data from server
// as well as signing up and signing in
// uses resty as http client
type Mediator interface {
	Sync(ctx context.Context) error
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) error
}

type mediator struct {
	client  *resty.Client
	storage storage.Storage
}

// NewMediator creates new mediator
func NewMediator(storage storage.Storage) *mediator {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	md := &mediator{
		client:  resty.NewWithClient(&http.Client{Transport: tr}),
		storage: storage,
	}
	return md
}

// SignUp sends request to server to create new user
// if request is successful, it will create session_id file
// and store session_id and username in it
// if request is not successful, it will return error
func (m *mediator) SignUp(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username or password is empty")
	}

	// Make request to server
	get, err := m.client.NewRequest().SetContext(ctx).
		SetQueryParam("email", username).
		SetQueryParam("password", password).Get("https://" + config.GetConfig().ServerIP + RegisterEndpoint)
	if err != nil {
		return err
	}
	if get.StatusCode() != 200 {
		return fmt.Errorf("failed to register, status code: %d, and error %s", get.StatusCode(), get.Body())
	}
	return setCookies(get.Cookies(), username)
}

// SignIn sends request to server to login user
// if request is successful, it will create session_id file
// and store session_id and username in it
func (m *mediator) SignIn(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username or password is empty")
	}

	// Make request to server
	get, err := m.client.NewRequest().SetContext(ctx).
		SetQueryParam("email", username).
		SetQueryParam("password", password).Get("https://" + config.GetConfig().ServerIP + LoginEndpoint)
	if err != nil {
		return err
	}
	if get.StatusCode() != 200 {
		return fmt.Errorf("failed to login, status code: %d and response %s", get.StatusCode(), get.Body())
	}
	return setCookies(get.Cookies(), username)
}

// Sync sends request to server to get data then swap it with local data
func (m *mediator) Sync(ctx context.Context) error {
	// Make request to server to get data don't forget to set cookie
	response, err := m.client.NewRequest().SetContext(ctx).
		SetBody(m.storage.Get()).SetCookie(&http.Cookie{
		Name:  "session_id",
		Value: auth.CurrentUser.SessionID,
	}).Post("https://" + config.GetConfig().ServerIP + Endpoint)
	if err != nil {
		return err
	}

	// Check if user is unauthorized
	if response.StatusCode() == http.StatusUnauthorized {
		// If so, try to get pass from keyring
		pass, err_ := keyring.Get(config.ServiceName, auth.CurrentUser.Username)
		if err_ != nil {
			return err_
		}

		// And sign in again
		err_ = m.SignIn(ctx, auth.CurrentUser.Username, pass)
		if err_ != nil {
			return err_
		}
	}

	// Check if response is empty
	if response.StatusCode() == http.StatusNoContent {
		return nil
	}

	// If not, unmarshal and swap data
	var serverData []models.DataWrapper
	body := response.Body()
	if len(body) == 0 {
		return nil
	}
	if err = json.Unmarshal(body, &serverData); err != nil {
		if err.Error() == "EOF" {
			serverData = nil
			return nil
		} else {
			return err
		}
	}
	err = m.storage.Swap(serverData)
	if err != nil {
		return err
	}
	return nil
}

func setCookies(cookie []*http.Cookie, username string) error {
	// Read cookie from response
	if len(cookie) == 0 {
		return fmt.Errorf("failed to get cookie")
	}

	// Loop through cookies and the one with session_id
	for _, c := range cookie {
		if c.Name == "session_id" {
			if c.Value == "" {
				return fmt.Errorf("failed to get cookie")
			}
			auth.CurrentUser.Username = username
			auth.CurrentUser.SessionID = c.Value
			return nil
		}
	}
	return fmt.Errorf("failed to get cookie")
}
