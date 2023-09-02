package UI

import (
	"context"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
	"time"
)

func (u *ui) register() *tview.Form {
	var err error
	var secret string

	if config.CurrentUser.Username != "" {
		secret, err = keyring.Get(config.ServiceName, config.CurrentUser.Username)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get password from keyring")
		}
		u.mediator.Sync(context.Background())
	}

	form := tview.NewForm().
		AddTextArea("Email", config.CurrentUser.Username, 30, 1, 100, func(text string) {
			config.CurrentUser.Username = text
		}).
		AddPasswordField("Password", secret, 30, '*', func(text string) {
			secret = text
		}).
		AddButton("SignUp", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = u.mediator.SignUp(ctx, config.CurrentUser.Username, secret)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "register"), true)
				return
			}
			err = keyring.Set(config.ServiceName, config.CurrentUser.Username, secret)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "register"), true)
				return
			}
			u.pages.SwitchToPage("menu")
			return
		}).AddButton("SignIn", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = u.mediator.SignIn(ctx, config.CurrentUser.Username, secret)
		if err != nil {
			u.pages.AddAndSwitchToPage("error", u.errorModal(err, "register"), true)
			return
		}
		err = keyring.Set(config.ServiceName, config.CurrentUser.Username, secret)
		if err != nil {
			u.pages.AddAndSwitchToPage("error", u.errorModal(err, "register"), true)
			return
		}
		u.pages.SwitchToPage("menu")
		return
	}).AddButton("Quit", func() {
		u.app.Stop()
	})
	form.SetBorder(true).SetTitle("SignUp or login").SetTitleAlign(tview.AlignLeft)
	return form
}
