package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	auth "github.com/gynshu-one/goph-keeper/server/internal/auth"
	"github.com/gynshu-one/goph-keeper/server/internal/config"
	router "github.com/gynshu-one/goph-keeper/server/internal/router"
	server "github.com/gynshu-one/goph-keeper/server/internal/server"
	"github.com/gynshu-one/goph-keeper/server/pkg/storage"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	connOpts := options.Client().ApplyURI(config.GetConfig().MongoURI)

	ctx := context.Background()

	client, err := mongo.Connect(ctx, connOpts)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	log.Info().Msg("Successfully connected to MongoDB")
	db := client.Database("goph-keeper")

	// init storage
	storage := storage.NewStorage(db)

	// init handlers
	handlers := server.NewHandlers(db, storage)

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
