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

type companyMongoRepository struct {
	collection *mongo.Collection
}

func NewCompanyMongoRepository(db *mongo.Database) domain.CompanyRepository {
	return &companyMongoRepository{
		collection: db.Collection("companies"),
	}
}

func (r *companyMongoRepository) Create(ctx context.Context, company *domain.Company) error {
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("COMPANY_ALREADY_EXISTS", "Company name already exists", 409, err, nil)
		}
		return errors.New("DATABASE_ERROR", "Failed to create company", 500, err, nil)
	}

	company.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *companyMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Company, error) {
	var company domain.Company
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("COMPANY_NOT_FOUND", "Company not found", 404, err, nil)
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to get company", 500, err, nil)
	}
	return &company, nil
}

func (r *companyMongoRepository) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	var company domain.Company
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("COMPANY_NOT_FOUND", "Company not found", 404, err, nil)
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to get company", 500, err, nil)
	}
	return &company, nil
}

func (r *companyMongoRepository) GetAll(ctx context.Context) ([]*domain.Company, error) {
	// Enhanced aggregation pipeline untuk populate user data
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user",
				"foreignField": "_id",
				"as":           "userDetails",
			},
		},
		{
			"$project": bson.M{
				"_id":             1,
				"name":            1,
				"profilePicture":  1,
				"user":            1,
				"createdAt":       1,
				"updatedAt":       1,
				"userDetails": bson.M{
					"$map": bson.M{
						"input": "$userDetails",
						"as":    "user",
						"in": bson.M{
							"_id":  "$$user._id",
							"name": "$$user.name",
						},
					},
				},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get companies", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var companies []*domain.Company
	if err = cursor.All(ctx, &companies); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode companies", 500, err, nil)
	}

	return companies, nil
}

func (r *companyMongoRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Company, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user": userID})
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get user companies", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var companies []*domain.Company
	if err = cursor.All(ctx, &companies); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode companies", 500, err, nil)
	}

	return companies, nil
}

func (r *companyMongoRepository) Update(ctx context.Context, id primitive.ObjectID, company *domain.Company) error {
	company.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":           company.Name,
			"profilePicture": company.ProfilePicture,
			"user":           company.User,
			"updatedAt":      company.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("COMPANY_ALREADY_EXISTS", "Company name already exists", 409, err, nil)
		}
		return errors.New("DATABASE_ERROR", "Failed to update company", 500, err, nil)
	}

	if result.MatchedCount == 0 {
		return errors.New("COMPANY_NOT_FOUND", "Company not found", 404, nil, nil)
	}

	return nil
}

func (r *companyMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.New("DATABASE_ERROR", "Failed to delete company", 500, err, nil)
	}

	if result.DeletedCount == 0 {
		return errors.New("COMPANY_NOT_FOUND", "Company not found", 404, nil, nil)
	}

	return nil
}