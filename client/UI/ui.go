package UI

import (
	"errors"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/client/sync"
	"github.com/rivo/tview"
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
	go func() {
		for {
			select {
			case err := <-config.ErrChan:
				newErr := errors.New("sorry Some error during Application work " + err.Error())
				u.pages.AddAndSwitchToPage("error", u.errorModal(newErr, "register"), true)
				u.pages.ShowPage("error")
			}
		}
	}()
	return u
}

func (u *ui) Pages() *tview.Pages {
	u.pages.AddPage("register", u.grid(u.register()), true, true)
	u.pages.AddPage("menu", u.grid(u.Menu()), true, false)
	u.pages.AddPage("show_items", u.grid(u.itemsTable()), true, false)
	u.pages.AddPage("add_items", u.grid(u.AddData()), true, false)

	u.pages.AddPage("text", u.grid(u.text()), true, false)
	u.pages.AddPage("bank_card", u.grid(u.bankCard()), true, false)
	u.pages.AddPage("binary", u.grid(u.binary()), true, false)
	u.pages.AddPage("login", u.grid(u.login()), true, false)

	u.pages.AddPage("success", u.grid(u.success()), true, false)
	return u.pages
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

func (u *ui) success() *tview.Modal {
	modal := tview.NewModal().
		SetText("Success").
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			u.pages.SwitchToPage("menu")
		})
	modal.SetBorder(true).SetTitle("Success").SetTitleAlign(tview.AlignLeft)
	return modal
}

func (u *ui) AddData() *tview.Form {
	buttonRow := tview.NewForm().AddButton("Text", func() {
		u.pages.SwitchToPage("text")
	}).AddButton("Bank Card", func() {
		u.pages.SwitchToPage("bank_card")
	}).AddButton("Binary", func() {
		u.pages.SwitchToPage("binary")
	}).AddButton("Login", func() {
		u.pages.SwitchToPage("login")
	}).AddButton("Back", func() {
		u.pages.SwitchToPage("menu")
	})
	return buttonRow
}

func (u *ui) Menu() *tview.Form {
	buttonRow := tview.NewForm().
		AddButton("Add Item", func() {
			u.pages.SwitchToPage("add_items")
		}).AddButton("Show Items", func() {
		u.pages.SwitchToPage("show_items")
	})
	return buttonRow
}

func (u *ui) itemsTable() *tview.List {
	items := u.storage.Get()
	list := tview.NewList()
	list.AddItem("Back", "", 'q', func() {
		u.pages.SwitchToPage("menu")
	})
	for tp, slice := range items {
		for _, item := range slice {
			list.AddItem(item.GetName(), "type of "+string(tp), 0, func() {
				u.pages.SwitchToPage(string(tp))
			})
		}
	}
	//table := tview.NewTable().
	//	SetBorders(true)
	//
	//rows := 0
	//for cols, userDataType := range models.UserDataTypes {
	//	table.SetCell(rows, cols, tview.NewTableCell(string(userDataType)).SetAlign(tview.AlignCenter))
	//	for _, item := range items[userDataType] {
	//		rows++
	//		table.SetCell(rows, cols, tview.NewTableCell(string(item.GetDataID())).
	//			SetTextColor(tcell.ColorWhite).
	//			SetAlign(tview.AlignCenter))
	//	}
	//}
	////for dataType, item := range items {
	////	switch dataType {
	////	case models.ArbitraryTextType:
	////		for i := 0; i < len(item); i++ {
	////			text := item[i].(*models.ArbitraryText)
	////			list.AddItem(text.Name, text.Text, '1', func() {
	////
	////			})
	////		}
	////	case models.BankCardType:
	////		for i := 0; i < len(item); i++ {
	////			bankCard := item[i].(*models.BankCard)
	////			table.SetCell(i, 0, tview.NewTableCell(bankCard.Name))
	////			table.SetCell(i, 1, tview.NewTableCell(bankCard.Info))
	////			table.SetCell(i, 2, tview.NewTableCell(bankCard.CardNum))
	////			table.SetCell(i, 3, tview.NewTableCell(bankCard.CardName))
	////			table.SetCell(i, 4, tview.NewTableCell(bankCard.CardCvv))
	////			table.SetCell(i, 5, tview.NewTableCell(bankCard.CardExp))
	////		}
	////	case models.BinaryType:
	////		for i := 0; i < len(item); i++ {
	////			binary := item[i].(*models.Binary)
	////			table.SetCell(i, 0, tview.NewTableCell(binary.Name))
	////			table.SetCell(i, 1, tview.NewTableCell(binary.Info))
	////			table.SetCell(i, 2, tview.NewTableCell(string(binary.Binary)))
	////		}
	////	case models.LoginType:
	////		for i := 0; i < len(item); i++ {
	////			login := item[i].(*models.Login)
	////			table.SetCell(i, 0, tview.NewTableCell(login.Name))
	////			table.SetCell(i, 1, tview.NewTableCell(login.Info))
	////			table.SetCell(i, 2, tview.NewTableCell(login.Username))
	////			table.SetCell(i, 3, tview.NewTableCell(login.Password))
	////			table.SetCell(i, 4, tview.NewTableCell(login.OneTimeOrigin))
	////			table.SetCell(i, 5, tview.NewTableCell(login.RecoveryCodes))
	////		}
	////	}
	////}
	//if rows == 0 {
	//	table.SetCell(0, 0, tview.NewTableCell("No items"))
	//}
	//table.SetBorder(true).SetTitle("Items").SetTitleAlign(tview.AlignLeft)

	return list
}

func (u *ui) grid(main tview.Primitive) *tview.Grid {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Menu"), 0, 0, 1, 3, 0, 0, false).
		AddItem(tview.NewButton("Quit").SetSelectedFunc(func() {
			u.app.Stop()
		}), 2, 0, 1, 3, 0, 0, false)

	grid.AddItem(main, 1, 0, 1, 3, 0, 0, false)

	return grid
}
