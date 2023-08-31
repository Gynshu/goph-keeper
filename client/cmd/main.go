//package main
//
//import (
//	"context"
//	"github.com/gynshu-one/goph-keeper/client/storage"
//	"github.com/gynshu-one/goph-keeper/client/sync"
//)
//
//func main() {
//	ctx := context.Background()
//	client := sync.NewMediator(storage.NewStorage())
//	client.Sync(ctx)
//}

package main

import (
	"github.com/gynshu-one/goph-keeper/client/UI"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

func main() {
	err := config.NewConfig("client/config.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file please check if it exists and is valid" +
			"Config should be in json format and contain SERVER_IP, POLL_TIMER, DUMP_TIMER")
	}

	//// set password
	//err := keyring.Set(service, user, password)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// get password
	//secret, err := keyring.Get(service, user)
	//if err != nil {
	//	log.Fatal(err)
	//}

	app := tview.NewApplication()

	newUI := UI.NewUI(app)
	form := newUI.RegisterForm()

	form.SetBorder(true).SetTitle("  SignUp or login, (if you already have account on this machine you will be logged in)  ").SetTitleAlign(tview.AlignLeft)
	if err = app.SetRoot(newUI.Pages(), true).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app")
	}
}
