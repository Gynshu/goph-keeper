package main

import (
	auth "github.com/gynshu-one/goph-keeper/server/api/auth"
	server "github.com/gynshu-one/goph-keeper/server/api/handlers"
	"github.com/gynshu-one/goph-keeper/server/api/router"
	"github.com/gynshu-one/goph-keeper/server/config"
	"github.com/gynshu-one/goph-keeper/server/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	db := config.NewDb()

	// init storage
	newStorage := storage.NewStorage(db.Collection("user-data"), db.Collection("users"))

	// init handlers
	handlers := server.NewHandlers(newStorage)

	// init sessions
	auth.Sessions = auth.NewSessionManager()

	r := router.NewRouter(handlers)

	log.Info().Msg("Starting server")

	go func() {
		log.Info().Msgf("Listening on %s", config.GetConfig().HttpServerPort)
		//err := http.ListenAndServeTLS(":"+config.GetConfig().HttpServerPort, "cert/server-cert.pem", "cert/server-key.pem", r)
		err := http.ListenAndServeTLS(
			":"+config.GetConfig().HttpServerPort,
			config.GetConfig().CertFilePath,
			config.GetConfig().KeyFilePath,
			r)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to listen")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutdown Server ...")
	//TODO:
	log.Info().Msg("Server exiting")

}
