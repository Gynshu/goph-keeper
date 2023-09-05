package UI

import (
	"context"
	"fmt"
	"github.com/gynshu-one/goph-keeper/client/auth"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/common/utils"
	"github.com/rivo/tview"
	"time"
)

func (u *ui) register() *tview.Form {
	var err error
	var pass string
	var secret string

	if config.CurrentUser.Username != "" {
		pass = auth.GetPass()
		secret = auth.GetSecret()
	}

	form := tview.NewForm().
		AddTextArea("Email", config.CurrentUser.Username, 30, 1, 100, func(text string) {
			config.CurrentUser.Username = text
		}).
		AddPasswordField("Password (for server)", pass, 30, '*', func(text string) {
			pass = text
		}).AddPasswordField("Master Key (for local encryption)", secret, 30, '*', func(text string) {
		secret = text
	}).
		AddButton("SignUp", func() {
			if secret == "" || pass == "" || config.CurrentUser.Username == "" {
				u.throwModal(fmt.Errorf("please fill all data"), "register")
				return
			}
			// validate email
			if !utils.ValidateEmail(config.CurrentUser.Username) {
				u.throwModal(fmt.Errorf("please enter valid email"), "register")
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = u.mediator.SignUp(ctx, config.CurrentUser.Username, pass)
			if err != nil {
				u.throwModal(err, "register")
				return
			}
			auth.SetSecret(secret)
			auth.SetPass(pass)
			u.goToMenu()
			return
		}).AddButton("SignIn", func() {
		if secret == "" || pass == "" || config.CurrentUser.Username == "" {
			u.throwModal(fmt.Errorf("please fill all data"), "register")
			return
		}
		if !utils.ValidateEmail(config.CurrentUser.Username) {
			u.throwModal(fmt.Errorf("please enter valid email"), "register")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = u.mediator.SignIn(ctx, config.CurrentUser.Username, pass)
		if err != nil {
			u.throwModal(err, "register")
			return
		}
		auth.SetSecret(secret)
		auth.SetPass(pass)
		u.goToMenu()
		return
	})
	form.SetBorder(true).SetTitle(" SignUp or login (for simplicity your login info will be saved in OS keychain)").SetTitleAlign(tview.AlignLeft)
	return form
}
