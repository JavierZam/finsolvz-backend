package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"finsolvz-backend/internal/utils/log"
)

// CreateIndexes creates all necessary indexes for optimal performance
func CreateIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Users collection indexes
	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "resetPasswordToken", Value: 1}},
			Options: options.Index().SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "company", Value: 1}},
		},
	}

	// Reports collection indexes
	reportIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "company", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "reportType", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "createdBy", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "userAccess", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "reportName", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "year", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: -1}},
		},
		// Compound indexes for common queries
		{
			Keys: bson.D{{Key: "company", Value: 1}, {Key: "reportType", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "company", Value: 1}, {Key: "year", Value: 1}},
		},
	}

	// Companies collection indexes
	companyIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: -1}},
		},
	}

	// ReportTypes collection indexes
	reportTypeIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	// Create indexes
	collections := []struct {
		name    string
		indexes []mongo.IndexModel
	}{
		{"users", userIndexes},
		{"reports", reportIndexes},
		{"companies", companyIndexes},
		{"reporttypes", reportTypeIndexes},
	}

	for _, col := range collections {
		if len(col.indexes) > 0 {
			_, err := db.Collection(col.name).Indexes().CreateMany(ctx, col.indexes)
			if err != nil {
				log.Errorf(ctx, "Failed to create indexes for %s: %v", col.name, err)
				return err
			}
			log.Infof(ctx, "Created %d indexes for %s collection", len(col.indexes), col.name)
		}
	}

	return nil
}
