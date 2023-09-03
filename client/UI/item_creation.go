package UI

import (
	"fmt"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rivo/tview"
	"io"
	"os"
)

func (u *ui) text() *tview.Form {
	var data = &models.ArbitraryText{}
	form := tview.NewForm().
		AddInputField("Name", "", 30, func(in string, last rune) bool {
			data.Name = in
			return true
		}, nil).
		AddInputField("Text", "", 50, func(in string, last rune) bool {
			data.Text = in
			return true
		}, nil).
		AddButton("Add", func() {
			if data.Name == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("name is empty"), "text"), true)
				return
			}
			if data.Text == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("text is empty"), "text"), true)
				return
			}
			err := u.encryptSave(data, models.UserDataModel{
				Name: data.Name,
				Type: models.ArbitraryTextType,
			})
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "text"), true)
				return
			}
			u.goToMenu()
			return
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	form.SetBorder(true).SetTitle("Add text").SetTitleAlign(tview.AlignCenter)
	return form
}

func (u *ui) bankCard() *tview.Form {
	var data = &models.BankCard{}
	form := tview.NewForm().
		AddInputField("Name", "", 30, func(in string, last rune) bool {
			data.Name = in
			return true
		}, nil).
		AddInputField("Info", "", 30, func(in string, last rune) bool {
			data.Info = in
			return true
		}, nil).
		AddInputField("CardNum", "", 30, func(in string, last rune) bool {
			data.CardNum = in
			return true
		}, nil).
		AddInputField("CardName", "", 30, func(in string, last rune) bool {
			data.CardName = in
			return true
		}, nil).
		AddInputField("CardCvv", "", 30, func(in string, last rune) bool {
			data.CardCvv = in
			return true
		}, nil).
		AddInputField("CardExp", "", 30, func(in string, last rune) bool {
			data.CardExp = in
			return true
		}, nil).
		AddButton("Add", func() {
			if data.Name == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("name is empty"), "bank_card"), true)
				return
			}
			if data.CardNum == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("text is empty"), "bank_card"), true)
				return
			}
			err := u.encryptSave(data, models.UserDataModel{
				Name: data.Name,
				Type: models.BankCardType,
			})
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "bank_card"), true)
				return
			}
			u.goToMenu()
			return
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	form.SetBorder(true).SetTitle("Add bank card").SetTitleAlign(tview.AlignCenter)
	return form
}

func (u *ui) binary() *tview.Form {
	var data = &models.Binary{}
	path := ""
	form := tview.NewForm().
		AddInputField("Name", "", 30, func(in string, last rune) bool {
			data.Name = in
			return true
		}, nil).
		AddInputField("Info", "", 30, func(in string, last rune) bool {
			data.Info = in
			return true
		}, nil).
		AddInputField("Path", "", 30, func(in string, last rune) bool {
			path = in
			return true
		}, nil).
		AddButton("Add", func() {
			file, err := os.Open(path)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "binary"), true)
				return
			}
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					u.pages.AddAndSwitchToPage("error", u.errorModal(err, "binary"), true)
					return
				}
			}(file)
			readAll, err := io.ReadAll(file)
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "binary"), true)
				return
			}
			data.Binary = readAll
			if len(data.Binary) == 0 {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("file is empty"), "binary"), true)
				return
			}
			if data.Name == "" {
				u.pages.AddAndSwitchToPage("error", u.errorModal(fmt.Errorf("name is empty"), "binary"), true)
				return
			}
			err = u.encryptSave(data, models.UserDataModel{
				Name: data.Name,
				Type: models.BinaryType,
			})
			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "binary"), true)
				return
			}
			u.goToMenu()
			return
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	form.SetBorder(true).SetTitle("Add binary").SetTitleAlign(tview.AlignCenter)
	return form
}

func (u *ui) login() *tview.Form {
	var data = &models.Login{}
	form := tview.NewForm().
		AddInputField("Name", "", 30, func(in string, last rune) bool {
			data.Name = in
			return true
		}, nil).
		AddInputField("Info", "", 30, func(in string, last rune) bool {
			data.Info = in
			return true
		}, nil).
		AddInputField("Username", "", 30, func(in string, last rune) bool {
			data.Username = in
			return true
		}, nil).
		AddInputField("Password", "", 30, func(in string, last rune) bool {
			data.Password = in
			return true
		}, nil).
		AddInputField("OneTimeOrigin", "", 30, func(in string, last rune) bool {
			data.OneTimeOrigin = in
			return true
		}, nil).
		AddInputField("RecoveryCodes", "", 30, func(in string, last rune) bool {
			data.RecoveryCodes = in
			return true
		}, nil).
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

			err := u.encryptSave(data, models.UserDataModel{
				Name: data.Name,
				Type: models.LoginType,
			})

			if err != nil {
				u.pages.AddAndSwitchToPage("error", u.errorModal(err, "text"), true)
				return
			}
			u.goToMenu()
			return
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	form.SetBorder(true).SetTitle("Add login").SetTitleAlign(tview.AlignCenter)
	return form
}

func (u *ui) AddItemButtons() *tview.Form {
	return tview.NewForm().AddButton("New Text", func() {
		u.pages.SwitchToPage("text")
	}).AddButton("New Bank Card", func() {
		u.pages.SwitchToPage("bank_card")
	}).AddButton("New Binary", func() {
		u.pages.SwitchToPage("binary")
	}).AddButton("New Login", func() {
		u.pages.SwitchToPage("login")
	}).SetButtonsAlign(tview.AlignCenter)
}
