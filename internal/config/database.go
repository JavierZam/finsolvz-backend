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

	// Set client options optimized for production
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetMaxPoolSize(50)                    // Increased from 10
	clientOptions.SetMinPoolSize(5)                     // Maintain minimum connections
	clientOptions.SetMaxConnIdleTime(10 * time.Minute) // Longer idle time
	clientOptions.SetTimeout(5 * time.Second)          // Faster timeout for failed connections
	clientOptions.SetMaxConnecting(10)                 // Limit concurrent connections

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
	
	// Create indexes for optimal performance (async, don't block startup)
	go func() {
		if err := CreateIndexes(database); err != nil {
			log.Warnf(context.Background(), "Failed to create some indexes: %v", err)
		}
	}()
	
	return database, nil
}
