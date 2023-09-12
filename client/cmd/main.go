// Package main implements the main entry point for the application.
package main

import (
	"fmt"

	"github.com/gynshu-one/goph-keeper/client/UI"
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
		buildDate = "11.09.2023"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)

	app := tview.NewApplication()
	newUI := UI.NewUI(app)
	pages := newUI.Pages()

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app")
	}
}
