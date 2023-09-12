package UI

import (
	"errors"
	"sort"
	"strings"

	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rivo/tview"
)

func (u *ui) itemsTable() (list *tview.List) {
	items := u.storage.Get()
	list = tview.NewList()
	// sort items
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	// add items to list func for clean code
	f := func(name string, form *tview.Form, item models.DataWrapper) {
		list.AddItem(item.Name, item.Type, 0, func() {
			u.pages.AddAndSwitchToPage(name,
				u.grid(u.addItemButtons(), form), true)
		})
	}

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

		switch elem := decrypt.(type) {
		case models.Login:
			f("login", u.login(elem, wrapper), item)
		case models.ArbitraryText:
			f("arbitrary text", u.text(elem, wrapper), item)
		case models.BankCard:
			f("bank card", u.bankCard(elem, wrapper), item)
		case models.Binary:
			f("binary", u.binary(elem, wrapper), item)
		}
	}

	return list
}
