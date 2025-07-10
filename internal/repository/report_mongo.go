package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/utils/errors"
)

type reportMongoRepository struct {
	collection *mongo.Collection
}

func NewReportMongoRepository(db *mongo.Database) domain.ReportRepository {
	return &reportMongoRepository{
		collection: db.Collection("reports"),
	}
}

func (r *reportMongoRepository) Create(ctx context.Context, report *domain.Report) error {
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, report)
	if err != nil {
		return errors.New("DATABASE_ERROR", "Failed to create report", 500, err, nil)
	}

	report.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// getPopulationPipeline creates an optimized aggregation pipeline for populating report references.
func (r *reportMongoRepository) getPopulationPipeline() []bson.M {
	return []bson.M{
		// Single lookup with pipeline for company (more efficient)
		{
			"$lookup": bson.M{
				"from":         "companies",
				"localField":   "company",
				"foreignField": "_id",
				"as":           "company",
				"pipeline": []bson.M{
					{
						"$project": bson.M{
							"_id":            1,
							"name":           1,
							"profilePicture": 1,
							"createdAt":      1,
							"updatedAt":      1,
						},
					},
				},
			},
		},
		// Single lookup with pipeline for reportType
		{
			"$lookup": bson.M{
				"from":         "reporttypes",
				"localField":   "reportType",
				"foreignField": "_id",
				"as":           "reportType",
				"pipeline": []bson.M{
					{
						"$project": bson.M{
							"_id":  1,
							"name": 1,
						},
					},
				},
			},
		},
		// Single lookup with pipeline for createdBy
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "createdBy",
				"foreignField": "_id",
				"as":           "createdBy",
				"pipeline": []bson.M{
					{
						"$project": bson.M{
							"_id":       1,
							"name":      1,
							"email":     1,
							"role":      1,
							"createdAt": 1,
							"updatedAt": 1,
						},
					},
				},
			},
		},
		// Single lookup with pipeline for userAccess
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userAccess",
				"foreignField": "_id",
				"as":           "userAccess",
				"pipeline": []bson.M{
					{
						"$project": bson.M{
							"_id":       1,
							"name":      1,
							"email":     1,
							"role":      1,
							"createdAt": 1,
							"updatedAt": 1,
						},
					},
				},
			},
		},
		// Single project stage to flatten single-item arrays
		{
			"$project": bson.M{
				"_id":        1,
				"reportName": 1,
				"year":       1,
				"currency":   1,
				"reportData": 1,
				"createdAt":  1,
				"updatedAt":  1,
				"company": bson.M{
					"$arrayElemAt": []interface{}{"$company", 0},
				},
				"reportType": bson.M{
					"$arrayElemAt": []interface{}{"$reportType", 0},
				},
				"createdBy": bson.M{
					"$arrayElemAt": []interface{}{"$createdBy", 0},
				},
				"userAccess": 1, // Keep as array
			},
		},
	}
}

func (r *reportMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.PopulatedReport, error) {
	pipeline := append([]bson.M{{"$match": bson.M{"_id": id}}}, r.getPopulationPipeline()...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get report", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode report", 500, err, nil)
	}

	if len(reports) == 0 {
		return nil, errors.New("REPORT_NOT_FOUND", "Report not found", 404, nil, nil)
	}

	return reports[0], nil
}

func (r *reportMongoRepository) GetByName(ctx context.Context, name string) (*domain.PopulatedReport, error) {
	pipeline := append([]bson.M{{"$match": bson.M{"reportName": name}}}, r.getPopulationPipeline()...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get report", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode report", 500, err, nil)
	}

	if len(reports) == 0 {
		return nil, errors.New("REPORT_NOT_FOUND", "Report not found", 404, nil, nil)
	}

	return reports[0], nil
}

func (r *reportMongoRepository) GetAll(ctx context.Context) ([]*domain.PopulatedReport, error) {
	cursor, err := r.collection.Aggregate(ctx, r.getPopulationPipeline())
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get reports", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode reports", 500, err, nil)
	}

	return reports, nil
}

// GetAllPaginated retrieves reports with pagination
func (r *reportMongoRepository) GetAllPaginated(ctx context.Context, skip, limit int) ([]*domain.PopulatedReport, int, error) {
	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, errors.New("DATABASE_ERROR", "Failed to count reports", 500, err, nil)
	}

	// Add pagination to pipeline
	pipeline := r.getPopulationPipeline()
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, errors.New("DATABASE_ERROR", "Failed to get reports", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, 0, errors.New("DATABASE_ERROR", "Failed to decode reports", 500, err, nil)
	}

	return reports, int(total), nil
}

func (r *reportMongoRepository) GetByCompany(ctx context.Context, companyID primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	pipeline := append([]bson.M{{"$match": bson.M{"company": companyID}}}, r.getPopulationPipeline()...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get reports by company", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode reports", 500, err, nil)
	}

	return reports, nil
}

func (r *reportMongoRepository) GetByCompanies(ctx context.Context, companyIDs []primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	pipeline := append([]bson.M{{"$match": bson.M{"company": bson.M{"$in": companyIDs}}}}, r.getPopulationPipeline()...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get reports by companies", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode reports", 500, err, nil)
	}

	return reports, nil
}

func (r *reportMongoRepository) GetByReportType(ctx context.Context, reportTypeID primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	pipeline := append([]bson.M{{"$match": bson.M{"reportType": reportTypeID}}}, r.getPopulationPipeline()...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get reports by report type", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode reports", 500, err, nil)
	}

	return reports, nil
}

func (r *reportMongoRepository) GetByUserAccess(ctx context.Context, userID primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	pipeline := append([]bson.M{{"$match": bson.M{"userAccess": userID}}}, r.getPopulationPipeline()...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get reports by user access", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode reports", 500, err, nil)
	}

	return reports, nil
}

func (r *reportMongoRepository) GetByCreatedBy(ctx context.Context, userID primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	pipeline := append([]bson.M{{"$match": bson.M{"createdBy": userID}}}, r.getPopulationPipeline()...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get reports by created by", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var reports []*domain.PopulatedReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode reports", 500, err, nil)
	}

	return reports, nil
}

func (r *reportMongoRepository) Update(ctx context.Context, id primitive.ObjectID, report *domain.Report) (*domain.PopulatedReport, error) {
	report.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"reportName": report.ReportName,
			"reportType": report.ReportType,
			"year":       report.Year,
			"company":    report.Company,
			"currency":   report.Currency,
			"createdBy":  report.CreatedBy,
			"userAccess": report.UserAccess,
			"reportData": report.ReportData,
			"updatedAt":  report.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to update report", 500, err, nil)
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("REPORT_NOT_FOUND", "Report not found", 404, nil, nil)
	}

	return r.GetByID(ctx, id)
}

func (r *reportMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.New("DATABASE_ERROR", "Failed to delete report", 500, err, nil)
	}

	if result.DeletedCount == 0 {
		return errors.New("REPORT_NOT_FOUND", "Report not found", 404, nil, nil)
	}

	return nil
}
