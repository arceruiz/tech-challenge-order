package repository

import (
	"context"
	"errors"
	"tech-challenge-order/internal/config"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrorNotFound = errors.New("entity not found")
	database      = "order"
)

func NewMongo() *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Get().DB.ConnectionString))
	if err != nil {
		log.Fatal().Err(err).Msg("an error occurred when try to connect to mongo")
	}

	return client.Database(database)
}
