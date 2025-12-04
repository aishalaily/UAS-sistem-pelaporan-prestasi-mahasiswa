package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Database

func ConnectMongo() error {
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_URI")).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect MongoDB: %w", err)
	}

	Mongo = client.Database(os.Getenv("MONGO_DB"))
	return nil
}
