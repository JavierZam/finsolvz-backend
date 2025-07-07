package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportType struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type ReportTypeRepository interface {
	Create(ctx context.Context, reportType *ReportType) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*ReportType, error)
	GetByName(ctx context.Context, name string) (*ReportType, error)
	GetAll(ctx context.Context) ([]*ReportType, error)
	Update(ctx context.Context, id primitive.ObjectID, reportType *ReportType) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}