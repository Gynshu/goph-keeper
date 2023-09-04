package UI

import (
	"github.com/gynshu-one/goph-keeper/client/storage"
	"github.com/gynshu-one/goph-keeper/client/sync"
	"github.com/gynshu-one/goph-keeper/common/models"
	"github.com/rivo/tview"
)

type UI interface {
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
	u.pages.AddPage("register", u.grid(nil, u.register()), true, true)

	u.pages.AddPage("text", u.grid(u.addItemButtons(), u.text(models.ArbitraryText{}, models.DataWrapper{})), true, false)
	u.pages.AddPage("bank_card", u.grid(u.addItemButtons(), u.bankCard(models.BankCard{}, models.DataWrapper{})), true, false)
	u.pages.AddPage("binary", u.grid(u.addItemButtons(), u.binary(models.Binary{}, models.DataWrapper{})), true, false)
	u.pages.AddPage("login", u.grid(u.addItemButtons(), u.login(models.Login{}, models.DataWrapper{})), true, false)

	return u.pages
}

func (u *ui) throwError(err error, fromPage string) {
	u.pages.AddAndSwitchToPage("error", tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			u.pages.SwitchToPage(fromPage)
		}), false)
	return
}

func (u *ui) goToMenu() {
	err := u.mediator.Sync(nil)
	if err != nil {
		u.throwError(err, "menu")
		return
	}
	u.pages.RemovePage("menu")
	u.pages.AddAndSwitchToPage("menu", u.grid(u.addItemButtons(), u.itemsTable()), true)
}

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
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(tview.NewButton("Quit").SetSelectedFunc(func() {
			u.app.Stop()
		}), 2, 0, 1, 3, 0, 0, false)

	//grid.AddItem(main, 1, 3, 1, 3, 0, 0, false)
	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(newPrimitive("1"), 0, 0, 0, 0, 0, 0, false).
		AddItem(elem, 1, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("1"), 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(newPrimitive("1"), 1, 0, 1, 1, 0, 100, false).
		AddItem(elem, 1, 1, 1, 1, 0, 100, false).
		AddItem(newPrimitive("1"), 1, 2, 1, 1, 0, 100, false)
	return grid
}
