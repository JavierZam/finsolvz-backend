package config

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"finsolvz-backend/internal/utils/errors"
	"finsolvz-backend/internal/utils/log"
)

func ConnectMongoDB(ctx context.Context) (*mongo.Database, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("MONGO_URI_MISSING", "MongoDB URI not configured", 500, nil, nil)
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetMaxPoolSize(10)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)
	clientOptions.SetTimeout(10 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.New("MONGO_CONNECTION_ERROR", "Failed to connect to MongoDB", 500, err, nil)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.New("MONGO_PING_ERROR", "Failed to ping MongoDB", 500, err, nil)
	}

	log.Infof(ctx, "Connected to MongoDB successfully")

	// Return the database instance
	database := client.Database("Finsolvz")
	return database, nil
}
