package UI

import (
	"context"

	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/client/sync"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rivo/tview"
)

// UI is an interface for the UI
type UI interface {
	Pages() *tview.Pages
}

type ui struct {
	pages    *tview.Pages
	app      *tview.Application
	mediator sync.Mediator
	storage  storage.Storage
}

// NewUI creates a new UI instance with the given application
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

// Pages returns the predefined dynamic pages such as register, menu, new item pages etc.
// This is basically UI routing
func (u *ui) Pages() *tview.Pages {
	u.pages.AddPage("register", u.grid(nil, u.register()), true, true)

	u.pages.AddPage("text", u.grid(u.addItemButtons(), u.text(models.ArbitraryText{}, models.DataWrapper{})), true, false)
	u.pages.AddPage("bank_card", u.grid(u.addItemButtons(), u.bankCard(models.BankCard{}, models.DataWrapper{})), true, false)
	u.pages.AddPage("binary", u.grid(u.addItemButtons(), u.binary(models.Binary{}, models.DataWrapper{})), true, false)
	u.pages.AddPage("login", u.grid(u.addItemButtons(), u.login(models.Login{}, models.DataWrapper{})), true, false)

	return u.pages
}

// throwModal throws a modal with the given message and redirects to the given page
func (u *ui) throwModal(message error, redirect string) {
	u.pages.AddAndSwitchToPage("error", tview.NewModal().
		SetText(message.Error()).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			u.pages.SwitchToPage(redirect)
		}), false)
}

// goToMenu redirects to the menu page
func (u *ui) goToMenu() {
	err := u.mediator.Sync(context.Background())
	if err != nil {
		u.throwModal(err, "menu")
		return
	}
	u.pages.RemovePage("menu")
	u.pages.AddAndSwitchToPage("menu", u.grid(u.addItemButtons(), u.itemsTable()), true)
}

// grid creates a grid lyout with the given header and element
func (u *ui) grid(header tview.Primitive, elem tview.Primitive) *tview.Grid {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	if header == nil {
		header = newPrimitive("Header")
	}
	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewButton("Quit").SetSelectedFunc(func() {
			u.app.Stop()
		}), 2, 0, 1, 1, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(elem, 1, 0, 1, 1, 0, 0, false)

	return grid
}
