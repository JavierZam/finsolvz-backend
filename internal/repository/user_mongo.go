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

type userMongoRepository struct {
	collection *mongo.Collection
}

func NewUserMongoRepository(db *mongo.Database) domain.UserRepository {
	return &userMongoRepository{
		collection: db.Collection("users"),
	}
}

func (r *userMongoRepository) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("USER_ALREADY_EXISTS", "Email already registered", 409, err, nil)
		}
		return errors.New("DATABASE_ERROR", "Failed to create user", 500, err, nil)
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *userMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("USER_NOT_FOUND", "User not found", 404, err, nil)
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to get user", 500, err, nil)
	}
	return &user, nil
}

func (r *userMongoRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	// Include password field for authentication
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("USER_NOT_FOUND", "User not found", 404, err, nil)
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to get user", 500, err, nil)
	}
	return &user, nil
}

// ✅ ENHANCED: GetAll method with comprehensive company field handling
func (r *userMongoRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	// Enhanced aggregation pipeline to handle ALL legacy data scenarios
	pipeline := []bson.M{
		{
			"$project": bson.M{
				"_id":       1,
				"name":      1,
				"email":     1,
				"role":      1,
				"createdAt": 1,
				"updatedAt": 1,
				// Enhanced company field handling for ALL scenarios
				"company": bson.M{
					"$switch": bson.M{
						"branches": []bson.M{
							{
								// Case 1: Field doesn't exist
								"case": bson.M{"$eq": []interface{}{bson.M{"$type": "$company"}, "missing"}},
								"then": []primitive.ObjectID{}, // Return empty array
							},
							{
								// Case 2: Field is null
								"case": bson.M{"$eq": []interface{}{"$company", nil}},
								"then": []primitive.ObjectID{}, // Return empty array
							},
							{
								// Case 3: Field is string (legacy format)
								"case": bson.M{"$eq": []interface{}{bson.M{"$type": "$company"}, "string"}},
								"then": []primitive.ObjectID{}, // Convert string to empty array for now
							},
							{
								// Case 4: Field is proper ObjectId array
								"case": bson.M{"$isArray": "$company"},
								"then": "$company", // Use as-is
							},
							{
								// Case 5: Field is single ObjectId
								"case": bson.M{"$eq": []interface{}{bson.M{"$type": "$company"}, "objectId"}},
								"then": []interface{}{"$company"}, // Wrap in array
							},
						},
						"default": []primitive.ObjectID{}, // Fallback to empty array
					},
				},
				// Remove unwanted fields from legacy data
				// Don't include: password, __v, resetPasswordToken, etc.
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to get users", 500, err, nil)
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, errors.New("DATABASE_ERROR", "Failed to decode users", 500, err, nil)
	}

	return users, nil
}

func (r *userMongoRepository) Update(ctx context.Context, id primitive.ObjectID, user *domain.User) error {
	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"email":     user.Email,
			"role":      user.Role,
			"company":   user.Company,
			"updatedAt": user.UpdatedAt,
		},
	}

	// Only update password if it's provided
	if user.Password != "" {
		update["$set"].(bson.M)["password"] = user.Password
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("EMAIL_ALREADY_EXISTS", "Email already used by another user", 409, err, nil)
		}
		return errors.New("DATABASE_ERROR", "Failed to update user", 500, err, nil)
	}

	if result.MatchedCount == 0 {
		return errors.New("USER_NOT_FOUND", "User not found", 404, nil, nil)
	}

	return nil
}

func (r *userMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.New("DATABASE_ERROR", "Failed to delete user", 500, err, nil)
	}

	if result.DeletedCount == 0 {
		return errors.New("USER_NOT_FOUND", "User not found", 404, nil, nil)
	}

	return nil
}

func (r *userMongoRepository) SetResetToken(ctx context.Context, email, token string, expires time.Time) error {
	update := bson.M{
		"$set": bson.M{
			"resetPasswordToken":   token,
			"resetPasswordExpires": expires,
			"updatedAt":            time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, update)
	if err != nil {
		return errors.New("DATABASE_ERROR", "Failed to set reset token", 500, err, nil)
	}

	if result.MatchedCount == 0 {
		return errors.New("USER_NOT_FOUND", "User not found", 404, nil, nil)
	}

	return nil
}

func (r *userMongoRepository) GetByResetToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{
		"resetPasswordToken":   token,
		"resetPasswordExpires": bson.M{"$gt": time.Now()},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("INVALID_TOKEN", "Invalid or expired token", 400, err, nil)
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to get user by reset token", 500, err, nil)
	}

	return &user, nil
}