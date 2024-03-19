package repository

import (
	"context"
	"errors"
	"tech-challenge-order/internal/config"

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
		panic(err)
	}
	db := client.Database(database)
	return db
}
