package UI

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rivo/tview"
)

func (u *ui) text(data models.ArbitraryText, wrapper models.DataWrapper) *tview.Form {
	if data.Text == "" {
		data = models.ArbitraryText{}
	}
	if wrapper.ID == "" {
		wrapper = models.DataWrapper{
			Type: models.ArbitraryTextType,
		}
	}
	form := tview.NewForm().
		AddInputField("Name", wrapper.Name, 30, nil, func(in string) {
			wrapper.Name = in
		}).
		AddInputField("Text", data.Text, 50, nil, func(in string) {
			data.Text = in
		}).
		AddButton("Save", func() {
			if wrapper.Name == "" {
				u.throwModal(fmt.Errorf("name is empty"), "text")
				return
			}
			if data.Text == "" {
				u.throwModal(fmt.Errorf("text is empty"), "text")
				return
			}

			err := u.storage.AddEncrypt(&data, wrapper)
			if err != nil {
				u.throwModal(err, "text")
				return
			}
			u.goToMenu()
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	// meaning we are creating item not editing
	if wrapper.ID != "" {
		form.AddButton("Delete", func() {
			err := u.storage.Delete(wrapper.ID)
			if err != nil {
				u.throwModal(err, "text")
				return
			}
			err = u.mediator.Sync(context.Background())
			if err != nil {
				u.throwModal(err, "text")
				return
			}
			u.throwModal(fmt.Errorf("item will be deleted from server in 30 days"), "text")
		})
		form.SetTitle(" Edit text ")
	} else {
		form.SetTitle(" Add text ")
	}
	form.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	return form
}

func (u *ui) bankCard(data models.BankCard, wrapper models.DataWrapper) *tview.Form {
	if data.CardNum == "" {
		data = models.BankCard{}
	}
	if wrapper.ID == "" {
		wrapper = models.DataWrapper{
			Type: models.BankCardType,
		}
	}
	form := tview.NewForm().
		AddInputField("Name", wrapper.Name, 30, nil, func(in string) {
			wrapper.Name = in
		}).
		AddInputField("Info", data.Info, 30, nil, func(in string) {
			data.Info = in
		}).
		AddInputField("CardNum", data.CardNum, 30, nil, func(in string) {
			data.CardNum = in
		}).
		AddInputField("CardName", data.CardName, 30, nil, func(in string) {
			data.CardName = in
		}).
		AddInputField("CardCvv", data.CardCvv, 30, nil, func(in string) {
			data.CardCvv = in
		}).
		AddInputField("CardExp", data.CardExp, 30, nil, func(in string) {
			data.CardExp = in
		}).
		AddButton("Save", func() {
			if wrapper.Name == "" {
				u.throwModal(fmt.Errorf("name is empty"), "bank_card")
				return
			}
			if data.CardNum == "" {
				u.throwModal(fmt.Errorf("card number is empty"), "bank_card")
				return
			}
			err := u.storage.AddEncrypt(&data, wrapper)
			if err != nil {
				u.throwModal(err, "bank_card")
				return
			}
			u.goToMenu()
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	// meaning we are creating item not editing
	if wrapper.ID != "" {
		form.AddButton("Delete", func() {
			err := u.storage.Delete(wrapper.ID)
			if err != nil {
				u.throwModal(err, "bank_card")
				return
			}
			err = u.mediator.Sync(context.Background())
			if err != nil {
				u.throwModal(err, "bank_card")
				return
			}
			u.throwModal(fmt.Errorf("item will be deleted from server in 30 days"), "bank_card")
		})
		form.SetTitle(" Edit bank card ")
	} else {
		form.SetTitle(" Add bank card ")
	}
	form.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	return form
}

func (u *ui) binary(data models.Binary, wrapper models.DataWrapper) *tview.Form {
	if data.Binary == nil {
		data = models.Binary{}
	}
	if wrapper.ID == "" {
		wrapper = models.DataWrapper{
			Type: models.BinaryType,
		}
	}
	path := ""
	form := tview.NewForm().
		AddInputField("Name", wrapper.Name, 30, nil, func(in string) {
			wrapper.Name = in
		}).
		AddInputField("Info", data.Info, 30, nil, func(in string) {
			data.Info = in
		}).
		AddInputField("Path", "", 30, nil, func(in string) {
			path = in
		}).
		AddButton("Save", func() {
			file, err := os.Open(path)
			if err != nil {
				u.throwModal(err, "binary")
				return
			}
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					u.throwModal(err, "binary")
					return
				}
			}(file)
			readAll, err := io.ReadAll(file)
			if err != nil {
				u.throwModal(err, "binary")
				return
			}
			data.Binary = readAll
			if len(data.Binary) == 0 {
				u.throwModal(fmt.Errorf("binary is empty"), "binary")
				return
			}
			if wrapper.Name == "" {
				u.throwModal(fmt.Errorf("name is empty"), "binary")
				return
			}

			err = u.storage.AddEncrypt(&data, wrapper)
			if err != nil {
				u.throwModal(err, "binary")
				return
			}
			u.goToMenu()
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	// meaning we are creating item not editing
	if wrapper.ID != "" {
		form.AddButton("Delete", func() {
			err := u.storage.Delete(wrapper.ID)
			if err != nil {
				u.throwModal(err, "binary")
				return
			}
			err = u.mediator.Sync(context.Background())
			if err != nil {
				u.throwModal(err, "binary")
				return
			}
			u.throwModal(fmt.Errorf("item will be deleted from server in 30 days"), "binary")
		})
		form.SetTitle(" Edit binary ")
	} else {
		form.SetTitle(" Add binary ")
	}
	form.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	return form
}

func (u *ui) login(data models.Login, wrapper models.DataWrapper) *tview.Form {
	if data.Username == "" {
		data = models.Login{}
	}
	if wrapper.ID == "" {
		wrapper = models.DataWrapper{
			Type: models.LoginType,
		}
	}
	form := tview.NewForm().
		AddInputField("Name", wrapper.Name, 30, nil, func(in string) {
			wrapper.Name = in
		}).
		AddInputField("Info", data.Info, 30, nil, func(in string) {
			data.Info = in
		}).
		AddInputField("Username", data.Username, 30, nil, func(in string) {
			data.Username = in
		}).
		AddInputField("Password", data.Password, 30, nil, func(in string) {
			data.Password = in
		}).
		AddInputField("OneTimeOrigin", data.OneTimeOrigin, 30, nil, func(in string) {
			data.OneTimeOrigin = in
		}).
		AddInputField("RecoveryCodes", data.RecoveryCodes, 30, nil, func(in string) {
			data.RecoveryCodes = in
		}).
		AddButton("Save", func() {
			if wrapper.Name == "" {
				u.throwModal(fmt.Errorf("name is empty"), "login")
				return
			}
			if data.Username == "" {
				u.throwModal(fmt.Errorf("username is empty"), "login")
				return
			}
			if data.Password == "" {
				u.throwModal(fmt.Errorf("password is empty"), "login")
				return
			}
			if data.OneTimeOrigin != "" {
				_, _, err := data.GenerateOneTimePassword()
				if err != nil {
					u.throwModal(fmt.Errorf("invalid one time origin"), "login")
					return
				}
			}
			err := u.storage.AddEncrypt(&data, wrapper)
			if err != nil {
				u.throwModal(err, "login")
				return
			}
			u.goToMenu()
		}).AddButton("Back", func() {
		u.goToMenu()
	})

	// meaning we are creating item not editing
	if wrapper.ID != "" {
		form.AddButton("Delete", func() {
			err := u.storage.Delete(wrapper.ID)
			if err != nil {
				u.throwModal(err, "login")
				return
			}
			err = u.mediator.Sync(context.Background())
			if err != nil {
				u.throwModal(err, "login")
				return
			}
			u.throwModal(fmt.Errorf("item will be deleted from server in 30 days"), "login")
		})
	}
	if data.OneTimeOrigin != "" {
		form.AddButton("Get One Time Password", func() {
			oneTime, _, err := data.GenerateOneTimePassword()
			if err != nil {
				u.throwModal(err, "login")
				return
			}
			u.pages.AddAndSwitchToPage("one_time", tview.NewModal().
				SetText(oneTime).
				AddButtons([]string{"Ok"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					u.goToMenu()
				}), false)
		})
		form.SetTitle(" Edit login ")
	} else {
		form.SetTitle(" Add login ")
	}
	form.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	return form
}

func (u *ui) addItemButtons() *tview.Form {
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
