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
	u.pages.AddPage("register", u.grid(nil, nil, u.register(), nil), true, true)

	u.pages.AddPage("text", u.grid(u.AddItemButtons(), nil, u.text(), nil), true, false)
	u.pages.AddPage("bank_card", u.grid(u.AddItemButtons(), nil, u.bankCard(), nil), true, false)
	u.pages.AddPage("binary", u.grid(u.AddItemButtons(), nil, u.binary(), nil), true, false)
	u.pages.AddPage("login", u.grid(u.AddItemButtons(), nil, u.login(), nil), true, false)

	u.pages.AddPage("success", u.grid(nil, nil, u.success(), nil), true, false)
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
			u.goToMenu()
		})
	modal.SetBorder(true).SetTitle("Success").SetTitleAlign(tview.AlignLeft)
	return modal
}

func (u *ui) itemsTable() (first, second, third *tview.List) {
	items := u.storage.Get()
	first = tview.NewList()
	second = tview.NewList()
	third = tview.NewList()

	for i, item := range items {
		switch i % 3 {
		case 0:
			second.AddItem(item.Name, item.Type, 0, nil)
		case 1:
			first.AddItem(item.Name, item.Type, 0, nil)
		case 2:
			third.AddItem(item.Name, item.Type, 0, nil)
		}
	}

	return first, second, third
}

func (u *ui) goToMenu() {
	u.pages.RemovePage("menu")
	f, s, t := u.itemsTable()
	u.pages.AddAndSwitchToPage("menu", u.grid(u.AddItemButtons(), f, s, t), true)
}

func (u *ui) grid(header tview.Primitive, first, second, third tview.Primitive) *tview.Grid {
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

	if first == nil {
		first = newPrimitive("First")
	}
	if second == nil {
		second = newPrimitive("Second")
	}
	if third == nil {
		third = newPrimitive("Third")
	}

	//grid.AddItem(main, 1, 3, 1, 3, 0, 0, false)
	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(first, 0, 0, 0, 0, 0, 0, false).
		AddItem(second, 1, 0, 1, 3, 0, 0, false).
		AddItem(third, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(first, 1, 0, 1, 1, 0, 100, true).
		AddItem(second, 1, 1, 1, 1, 0, 100, true).
		AddItem(third, 1, 2, 1, 1, 0, 100, true)
	return grid
}
