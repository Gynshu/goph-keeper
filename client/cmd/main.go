package main

import (
	"fmt"
	"github.com/gynshu-one/goph-keeper/client/UI"
	"github.com/gynshu-one/goph-keeper/client/config"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	if buildVersion == "" {
		buildVersion = "1.0.0"
	}
	if buildDate == "" {
		buildDate = "09.05.2023"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file please check if it exists and is valid" +
			"Config should be in json format and contain SERVER_IP, POLL_TIMER, DUMP_TIMER")
	}

	app := tview.NewApplication()
	newUI := UI.NewUI(app)
	pages := newUI.Pages()

	if err = app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app")
	}
}
