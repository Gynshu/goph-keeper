package UI

import (
	"fmt"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rivo/tview"
	"io"
	"os"
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
		AddButton("Submit", func() {
			if wrapper.Name == "" {
				u.throwError(fmt.Errorf("name is empty"), "text")
				return
			}
			if data.Text == "" {
				u.throwError(fmt.Errorf("text is empty"), "text")
				return
			}

			err := u.storage.AddEncrypt(&data, wrapper)
			if err != nil {
				u.throwError(err, "text")
				return
			}
			u.goToMenu()
			return
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	// meaning we are creating item not editing
	if wrapper.ID != "" {
		form.AddButton("Delete", func() {
			err := u.storage.Delete(wrapper.ID)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			err = u.mediator.Sync(nil)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			u.goToMenu()
		})
	}
	form.SetBorder(true).SetTitle("Add text").SetTitleAlign(tview.AlignCenter)
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
		AddButton("Submit", func() {
			if wrapper.Name == "" {
				u.throwError(fmt.Errorf("name is empty"), "bank_card")
				return
			}
			if data.CardNum == "" {
				u.throwError(fmt.Errorf("card number is empty"), "bank_card")
				return
			}
			err := u.storage.AddEncrypt(&data, wrapper)
			if err != nil {
				u.throwError(err, "bank_card")
				return
			}
			u.goToMenu()
			return
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	// meaning we are creating item not editing
	if wrapper.ID != "" {
		form.AddButton("Delete", func() {
			err := u.storage.Delete(wrapper.ID)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			err = u.mediator.Sync(nil)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			u.goToMenu()
		})
	}
	form.SetBorder(true).SetTitle("Add bank card").SetTitleAlign(tview.AlignCenter)
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
		AddButton("Submit", func() {
			file, err := os.Open(path)
			if err != nil {
				u.throwError(err, "binary")
				return
			}
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					u.throwError(err, "binary")
					return
				}
			}(file)
			readAll, err := io.ReadAll(file)
			if err != nil {
				u.throwError(err, "binary")
				return
			}
			data.Binary = readAll
			if len(data.Binary) == 0 {
				u.throwError(fmt.Errorf("binary is empty"), "binary")
				return
			}
			if wrapper.Name == "" {
				u.throwError(fmt.Errorf("name is empty"), "binary")
				return
			}

			err = u.storage.AddEncrypt(&data, wrapper)
			if err != nil {
				u.throwError(err, "binary")
				return
			}
			u.goToMenu()
			return
		}).AddButton("Back", func() {
		u.goToMenu()
	})
	// meaning we are creating item not editing
	if wrapper.ID != "" {
		form.AddButton("Delete", func() {
			err := u.storage.Delete(wrapper.ID)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			err = u.mediator.Sync(nil)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			u.goToMenu()
		})
	}
	form.SetBorder(true).SetTitle("Add binary").SetTitleAlign(tview.AlignCenter)
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
		AddButton("Submit", func() {
			if wrapper.Name == "" {
				u.throwError(fmt.Errorf("name is empty"), "login")
				return
			}
			if data.Username == "" {
				u.throwError(fmt.Errorf("username is empty"), "login")
				return
			}
			if data.Password == "" {
				u.throwError(fmt.Errorf("password is empty"), "login")
				return
			}
			err := u.storage.AddEncrypt(&data, wrapper)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			u.goToMenu()
			return
		}).AddButton("Back", func() {
		u.goToMenu()
	})

	// meaning we are creating item not editing
	if wrapper.ID != "" {
		form.AddButton("Delete", func() {
			err := u.storage.Delete(wrapper.ID)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			err = u.mediator.Sync(nil)
			if err != nil {
				u.throwError(err, "login")
				return
			}
			u.goToMenu()
		})
	}
	form.SetBorder(true).SetTitle("Add login").SetTitleAlign(tview.AlignCenter)
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
