package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	server "github.com/gynshu-one/goph-keeper/server/api/handlers"
	"github.com/gynshu-one/goph-keeper/server/api/router"
	"github.com/gynshu-one/goph-keeper/server/config"
	"github.com/gynshu-one/goph-keeper/server/storage"

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
	db := config.NewDb()

	// Init storage
	newStorage := storage.NewStorage(db.Collection("user-data"), db.Collection("users"))

	// Init handlers
	handlers := server.NewHandlers(newStorage)

	r := router.NewRouter(handlers)

	log.Info().Msg("Starting server")

	go func() {
		// Run https server
		log.Info().Msgf("Listening https on %s", config.GetConfig().HttpServerPort)
		err := http.ListenAndServeTLS(
			":"+config.GetConfig().HttpServerPort,
			config.GetConfig().CertFilePath,
			config.GetConfig().KeyFilePath,
			r)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to listen")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutdown Server ...")
	log.Info().Msg("Server exiting")

}
