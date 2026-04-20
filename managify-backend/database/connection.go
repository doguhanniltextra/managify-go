package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB     *mongo.Database
	Client *mongo.Client
)

func Connect() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return fmt.Errorf("MONGO_URI environment variable not set")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		return fmt.Errorf("MONGO_DB environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(200).SetMinPoolSize(10)

	var err error
	maxRetries := 5
	backoff := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		Client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			err = Client.Ping(ctx, nil)
		}
		cancel()

		if err == nil {
			fmt.Println("MongoDB connected successfully!")
			DB = Client.Database(dbName)
			return nil
		}

		fmt.Printf("Failed to connect to MongoDB (attempt %d/%d): %v. Retrying in %v...\n", i+1, maxRetries, err, backoff)
		time.Sleep(backoff)
		backoff *= 2 // Exponential backoff
	}

	return fmt.Errorf("could not connect to MongoDB after %d attempts: %w", maxRetries, err)
}

func Disconnect() error {
	if Client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return Client.Disconnect(ctx)
}

func CheckHealth() error {
	if Client == nil {
		return fmt.Errorf("database client not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return Client.Ping(ctx, nil)
}
