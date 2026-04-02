package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() error {
	// MongoDB Atlas URI
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return fmt.Errorf("MONGO_URI environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(200).SetMinPoolSize(10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("MongoDB connected successfully!")

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		return fmt.Errorf("MONGO_DB environment variable not set")
	}

	DB = client.Database(dbName)
	return nil
}
