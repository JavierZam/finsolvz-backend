package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	ReportName string               `bson:"reportName" json:"reportName"`
	ReportType primitive.ObjectID   `bson:"reportType" json:"reportType"`
	Year       int                  `bson:"year" json:"year"`
	Company    primitive.ObjectID   `bson:"company" json:"company"`
	Currency   *string              `bson:"currency,omitempty" json:"currency"`
	CreatedBy  primitive.ObjectID   `bson:"createdBy" json:"createdBy"`
	UserAccess []primitive.ObjectID `bson:"userAccess" json:"userAccess"`
	ReportData interface{}          `bson:"reportData" json:"reportData"`
	CreatedAt  time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type PopulatedReport struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ReportName string             `bson:"reportName" json:"reportName"`
	ReportType *ReportType        `bson:"reportType" json:"reportType"`
	Year       int                `bson:"year" json:"year"`
	Company    *Company           `bson:"company" json:"company"`
	Currency   *string            `bson:"currency,omitempty" json:"currency"`
	CreatedBy  *User              `bson:"createdBy" json:"createdBy"`
	UserAccess []*User            `bson:"userAccess" json:"userAccess"`
	ReportData interface{}        `bson:"reportData" json:"reportData"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type ReportRepository interface {
	Create(ctx context.Context, report *Report) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*PopulatedReport, error)
	GetByName(ctx context.Context, name string) (*PopulatedReport, error)
	GetAll(ctx context.Context) ([]*PopulatedReport, error)
	GetByCompany(ctx context.Context, companyID primitive.ObjectID) ([]*PopulatedReport, error)
	GetByCompanies(ctx context.Context, companyIDs []primitive.ObjectID) ([]*PopulatedReport, error)
	GetByReportType(ctx context.Context, reportTypeID primitive.ObjectID) ([]*PopulatedReport, error)
	GetByUserAccess(ctx context.Context, userID primitive.ObjectID) ([]*PopulatedReport, error)
	GetByCreatedBy(ctx context.Context, userID primitive.ObjectID) ([]*PopulatedReport, error)
	Update(ctx context.Context, id primitive.ObjectID, report *Report) (*PopulatedReport, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
