package UI

import (
	"errors"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rivo/tview"
	"sort"
	"strings"
)

func (u *ui) itemsTable() (list *tview.List) {
	items := u.storage.Get()
	list = tview.NewList()
	// sort items
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})
	for _, item := range items {
		decrypt, wrapper, err := u.storage.FindDecrypt(item.ID)
		if err != nil {
			if errors.Is(err, models.ErrDeleted) {
				list.AddItem(item.Name, "----deleted----", 0, nil)
				continue
			}
			u.throwModal(err, "login")
			return
		}

		switch decrypt.(type) {
		case models.Login:
			login := decrypt.(models.Login)
			list.AddItem(item.Name, item.Type, 0, func() {
				u.pages.AddAndSwitchToPage("login",
					u.grid(u.addItemButtons(), u.login(login, wrapper)), true)
			})
		case models.ArbitraryText:
			text := decrypt.(models.ArbitraryText)
			list.AddItem(item.Name, item.Type, 0, func() {
				u.pages.AddAndSwitchToPage("text",
					u.grid(u.addItemButtons(), u.text(text, wrapper)), true)
			})
		case models.BankCard:
			bank := decrypt.(models.BankCard)
			list.AddItem(item.Name, item.Type, 0, func() {
				u.pages.AddAndSwitchToPage("bank_card",
					u.grid(u.addItemButtons(), u.bankCard(bank, wrapper)), true)
			})
		case models.Binary:
			bin := decrypt.(models.Binary)
			list.AddItem(item.Name, item.Type, 0, func() {
				u.pages.AddAndSwitchToPage("binary",
					u.grid(u.addItemButtons(), u.binary(bin, wrapper)), true)
			})
		}
	}

	return list
}
