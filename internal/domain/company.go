package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name           string               `bson:"name" json:"name"`
	ProfilePicture *string              `bson:"profilePicture,omitempty" json:"profilePicture"`
	User           []primitive.ObjectID `bson:"user" json:"user"`
	CreatedAt      time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Company, error)
	GetByName(ctx context.Context, name string) (*Company, error)
	SearchByName(ctx context.Context, name string) ([]*Company, error)
	GetAll(ctx context.Context) ([]*Company, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Company, error)
	Update(ctx context.Context, id primitive.ObjectID, company *Company) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
