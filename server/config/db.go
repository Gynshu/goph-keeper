package config

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDb() *mongo.Database {
	connOpts := options.Client().ApplyURI(GetConfig().MongoURI)

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
	return db
}
