package UI

import (
	"context"
	"fmt"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/client/sync"
	"github.com/gynshu-one/goph-keeper/shared/models"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
	"io"
	"os"
	"time"
)

type UI interface {
	register() *tview.Form
	errorModal(err error, page string) *tview.Modal
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

	return u
}

func (u *ui) Pages() *tview.Pages {
	u.pages.AddPage("login", u.register(), true, true)
	u.pages.AddPage("menu", u.AddData(), true, false)
	u.pages.AddPage("text", u.text(), true, false)
	u.pages.AddPage("bank_card", u.bankCard(), true, false)
	u.pages.AddPage("binary", u.binary(), true, false)
	u.pages.AddPage("login", u.login(), true, false)
	u.pages.AddPage("ok", u.okModal(), true, false)
	return u.pages
}
func (u *ui) register() *tview.Form {
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
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "login"), true)
				return
			}
			err = keyring.Set(config.ServiceName, config.CurrentUser.Username, secret)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "login"), true)
				return
			}
			return
		}).AddButton("SignIn", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = u.mediator.SignIn(ctx, config.CurrentUser.Username, secret)
		if err != nil {
			u.pages.AddAndSwitchToPage("error", u.errorModal(err, "login"), true)
			return
		}
		err = keyring.Set(config.ServiceName, config.CurrentUser.Username, secret)
		if err != nil {
			u.pages.AddAndSwitchToPage("error", u.errorModal(err, "login"), true)
			return
		}
		return
	}).AddButton("Quit", func() {
		u.app.Stop()
	})
	form.SetBorder(true).SetTitle("SignUp or login").SetTitleAlign(tview.AlignLeft)
	return form
}

func (u *ui) errorModal(err error, page string) *tview.Modal {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"Retry"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			u.pages.RemovePage("error")
			u.pages.SwitchToPage(page)
		})
	modal.SetBorder(true).SetTitle("Error").SetTitleAlign(tview.AlignLeft)
	return modal
}

func (u *ui) okModal() *tview.Modal {
	modal := tview.NewModal().
		SetText("Success").
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			u.pages.SwitchToPage("menu")
		})
	modal.SetBorder(true).SetTitle("Success").SetTitleAlign(tview.AlignLeft)
	return modal
}

func (u *ui) text() *tview.Form {
	var data = &models.ArbitraryText{}
	form := tview.NewForm().
		AddInputField("Name", "", 30, nil, func(in string) {
			data.Name = in
		}).
		AddInputField("Text", "", 50, nil, func(in string) {
			data.Text = in
		}).
		AddButton("Add", func() {
			if data.Name == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("name is empty"), "text"), true)
				return
			}
			if data.Text == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("text is empty"), "text"), true)
				return
			}
			u.storage.Add(data)
			u.pages.AddAndSwitchToPage("ok", u.okModal(), true)
			return
		}).AddButton("Quit", func() {
		u.app.Stop()
	})
	form.SetBorder(true).SetTitle("Add text").SetTitleAlign(tview.AlignLeft)
	return form
}

func (u *ui) bankCard() *tview.Form {
	var data = &models.BankCard{}
	form := tview.NewForm().
		AddInputField("Name", "", 30, nil, func(in string) {
			data.Name = in
		}).
		AddInputField("Info", "", 30, nil, func(in string) {
			data.Info = in
		}).
		AddInputField("CardNum", "", 30, nil, func(in string) {
			data.CardNum = in
		}).
		AddInputField("CardName", "", 30, nil, func(in string) {
			data.CardName = in
		}).
		AddInputField("CardCvv", "", 30, nil, func(in string) {
			data.CardCvv = in
		}).
		AddInputField("CardExp", "", 30, nil, func(in string) {
			data.CardExp = in
		}).
		AddButton("Add", func() {
			if data.Name == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("name is empty"), "bank_card"), true)
				return
			}
			if data.CardNum == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("text is empty"), "bank_card"), true)
				return
			}
			u.storage.Add(data)
			u.pages.AddAndSwitchToPage("ok", u.okModal(), true)
			return
		}).AddButton("Quit", func() {
		u.app.Stop()
	})
	form.SetBorder(true).SetTitle("Add bank card").SetTitleAlign(tview.AlignLeft)
	return form
}

func (u *ui) binary() *tview.Form {
	var data = &models.Binary{OwnerID: config.CurrentUser.Username}
	form := tview.NewForm().
		AddInputField("Name", "", 30, nil, func(in string) {
			data.Name = in
		}).
		AddInputField("Info", "", 30, nil, func(in string) {
			data.Info = in
		}).
		AddInputField("Path", "", 30, nil, func(in string) {
			file, err := os.Open(in)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "binary"), true)
				return
			}
			defer file.Close()
			readAll, err := io.ReadAll(file)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "binary"), true)
				return
			}
			data.Binary = readAll
		}).
		AddButton("Add", func() {
			if len(data.Binary) == 0 {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("file is empty"), "binary"), true)
				return
			}
			if data.Name == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("name is empty"), "binary"), true)
				return
			}
			u.storage.Add(data)
			u.pages.AddAndSwitchToPage("ok", u.okModal(), true)
			return
		}).AddButton("Quit", func() {
		u.app.Stop()
	})
	form.SetBorder(true).SetTitle("Add binary").SetTitleAlign(tview.AlignLeft)
	return form
}

func (u *ui) login() *tview.Form {
	var data = &models.Login{}
	form := tview.NewForm().
		AddInputField("Name", "", 30, nil, func(in string) {
			data.Name = in
		}).
		AddInputField("Info", "", 30, nil, func(in string) {
			data.Info = in
		}).
		AddInputField("Username", "", 30, nil, func(in string) {
			data.Username = in
		}).
		AddInputField("Password", "", 30, nil, func(in string) {
			data.Password = in
		}).
		AddInputField("OneTimeOrigin", "", 30, nil, func(in string) {
			data.OneTimeOrigin = in
		}).
		AddInputField("RecoveryCodes", "", 30, nil, func(in string) {
			data.RecoveryCodes = in
		}).
		AddButton("Add", func() {
			if data.Name == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("name is empty"), "login"), true)
				return
			}
			if data.Username == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("username is empty"), "login"), true)
				return
			}
			if data.Password == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("password is empty"), "login"), true)
				return
			}
			u.storage.Add(data)
			u.pages.AddAndSwitchToPage("ok", u.okModal(), true)
			return
		}).AddButton("Quit", func() {
		u.app.Stop()
	})
	form.SetBorder(true).SetTitle("Add login").SetTitleAlign(tview.AlignLeft)
	return form
}
func (u *ui) AddData() *tview.Modal {
	modal := tview.NewModal().SetText("Choose What to create").AddButtons([]string{"Text", "Bank Card", "Binary", "Login"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonIndex {
		case 0:
			u.pages.AddAndSwitchToPage("text", u.text(), true)
		case 1:
			u.pages.AddAndSwitchToPage("bank_card", u.bankCard(), true)
		case 2:
			u.pages.AddAndSwitchToPage("binary", u.binary(), true)
		case 3:
			u.pages.AddAndSwitchToPage("login", u.login(), true)
		}
	})
	return modal
}
