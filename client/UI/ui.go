package UI

import (
	"context"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/client/sync"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
	"time"
)

type UI interface {
	RegisterForm() *tview.Form
	ErrorPrim(err error) *tview.Modal
	Pages() *tview.Pages
}

type ui struct {
	pages    *tview.Pages
	app      *tview.Application
	mediator sync.Mediator
	storage  storage.Storage
}

func NewUI(app *tview.Application) UI {
	newStorage := storage.NewStorage()
	u := &ui{
		pages:    tview.NewPages(),
		app:      app,
		mediator: sync.NewMediator(newStorage),
		storage:  newStorage,
	}
	u.pages.AddPage("login", u.RegisterForm(), true, true)
	return u
}

func (u *ui) Pages() *tview.Pages {
	return u.pages
}

func (u *ui) RegisterForm() *tview.Form {
	var err error
	var secret string

	if config.CurrentUser.Username != "" {
		secret, err = keyring.Get(config.ServiceName, config.CurrentUser.Username)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get password from keyring")
		}
	}

	form := tview.NewForm().
		AddTextArea("Email", config.CurrentUser.Username, 30, 1, 100, func(text string) {
			config.CurrentUser.Username = text
		}).
		AddPasswordField("Password (Autofilled If already logged in)", secret, 30, '*', func(text string) {
			secret = text
		}).
		AddButton("SignUp", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = u.mediator.SignUp(ctx, config.CurrentUser.Username, secret)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.ErrorPrim(err), true)
				return
			}
			err = keyring.Set(config.ServiceName, config.CurrentUser.Username, secret)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.ErrorPrim(err), true)
				return
			}
			return
		}).AddButton("SignIn", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = u.mediator.SignIn(ctx, config.CurrentUser.Username, secret)
		if err != nil {
			u.pages.AddAndSwitchToPage("error", u.ErrorPrim(err), true)
			return
		}
		err = keyring.Set(config.ServiceName, config.CurrentUser.Username, secret)
		if err != nil {
			u.pages.AddAndSwitchToPage("error", u.ErrorPrim(err), true)
			return
		}
		return
	}).AddButton("Quit", func() {
		u.app.Stop()
	})
	form.SetBorder(true).SetTitle("SignUp or login").SetTitleAlign(tview.AlignLeft)
	return form
}

func (u *ui) ErrorPrim(err error) *tview.Modal {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"Retry"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			u.pages.RemovePage("error")
			u.pages.SwitchToPage("login")
		})
	modal.SetBorder(true).SetTitle("Error").SetTitleAlign(tview.AlignLeft)
	return modal
}
