package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/utils/errors"
)

type reportTypeMongoRepository struct {
	collection *mongo.Collection
}

func NewReportTypeMongoRepository(db *mongo.Database) domain.ReportTypeRepository {
	return &reportTypeMongoRepository{
		collection: db.Collection("reporttypes"),
	}
}

func (r *reportTypeMongoRepository) Create(ctx context.Context, reportType *domain.ReportType) error {
	result, err := r.collection.InsertOne(ctx, reportType)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("REPORT_TYPE_ALREADY_EXISTS", "Report type name already exists", 409, err, nil)
		}
		return errors.New("DATABASE_ERROR", "Failed to create report type", 500, err, nil)
	}

	reportType.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *reportTypeMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.ReportType, error) {
	var reportType domain.ReportType
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&reportType)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("REPORT_TYPE_NOT_FOUND", "Report type not found", 404, err, nil)
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to get report type", 500, err, nil)
	}
	return &reportType, nil
}

func (r *reportTypeMongoRepository) GetByName(ctx context.Context, name string) (*domain.ReportType, error) {
	var reportType domain.ReportType
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&reportType)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("REPORT_TYPE_NOT_FOUND", "Report type not found", 404, err, nil)
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to get report type", 500, err, nil)
	}
	return &reportType, nil
}

func (r *reportTypeMongoRepository) GetAll(ctx context.Context) ([]*domain.ReportType, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get report types", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reportTypes []*domain.ReportType
	if err = cursor.All(ctx, &reportTypes); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode report types", 500, err, nil)
	}

	return reportTypes, nil
}

func (r *reportTypeMongoRepository) Update(ctx context.Context, id primitive.ObjectID, reportType *domain.ReportType) error {
	update := bson.M{
		"$set": bson.M{
			"name": reportType.Name,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("REPORT_TYPE_ALREADY_EXISTS", "Report type name already exists", 409, err, nil)
		}
		return errors.New("DATABASE_ERROR", "Failed to update report type", 500, err, nil)
	}

	if result.MatchedCount == 0 {
		return errors.New("REPORT_TYPE_NOT_FOUND", "Report type not found", 404, nil, nil)
	}

	return nil
}

func (r *reportTypeMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.New("DATABASE_ERROR", "Failed to delete report type", 500, err, nil)
	}

	if result.DeletedCount == 0 {
		return errors.New("REPORT_TYPE_NOT_FOUND", "Report type not found", 404, nil, nil)
	}

	return nil
}
