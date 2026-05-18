package main

import (
	"context"
	"log"
	"os"
	"time"

	"courier-technical-test/go/internal/courier"
	"courier-technical-test/go/internal/database"
	"courier-technical-test/go/internal/router"
)

func main() {
	mongoURI := env("MONGODB_URI", "mongodb://mongodb:mongodb12345@localhost:27019/?authSource=admin")
	databaseName := env("MONGODB_DATABASE", "gradin-courier")
	port := env("APP_PORT", "8080")

	client, err := database.Connect(mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repository := courier.NewRepository(client.Database(databaseName))
	if err := repository.EnsureIndexes(ctx); err != nil {
		log.Fatal(err)
	}

	app := router.New(courier.NewHandler(courier.NewService(repository)))
	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
