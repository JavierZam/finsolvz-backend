package reporttype

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/utils/errors"
)

type Service interface {
	CreateReportType(ctx context.Context, req CreateReportTypeRequest) (*ReportTypeResponse, error)
	GetReportTypes(ctx context.Context) ([]*ReportTypeResponse, error)
	GetReportTypeByID(ctx context.Context, id string) (*ReportTypeResponse, error)
	GetReportTypeByName(ctx context.Context, name string) (*ReportTypeResponse, error)
	UpdateReportType(ctx context.Context, id string, req UpdateReportTypeRequest) (*ReportTypeResponse, error)
	DeleteReportType(ctx context.Context, id string) error
}

type service struct {
	reportTypeRepo domain.ReportTypeRepository
}

func NewService(reportTypeRepo domain.ReportTypeRepository) Service {
	return &service{
		reportTypeRepo: reportTypeRepo,
	}
}

func (s *service) CreateReportType(ctx context.Context, req CreateReportTypeRequest) (*ReportTypeResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidReportTypeName
	}

	existingReportType, err := s.reportTypeRepo.GetByName(ctx, name)
	if err == nil && existingReportType != nil {
		return nil, ErrReportTypeAlreadyExists
	}

	reportType := &domain.ReportType{
		Name: name,
	}

	if err := s.reportTypeRepo.Create(ctx, reportType); err != nil {
		return nil, err
	}

	response := ToReportTypeResponse(reportType)
	return &response, nil
}

func (s *service) GetReportTypes(ctx context.Context) ([]*ReportTypeResponse, error) {
	reportTypes, err := s.reportTypeRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*ReportTypeResponse, len(reportTypes))
	for i, reportType := range reportTypes {
		response := ToReportTypeResponse(reportType)
		responses[i] = &response
	}

	return responses, nil
}

func (s *service) GetReportTypeByID(ctx context.Context, id string) (*ReportTypeResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_REPORT_TYPE_ID", "Invalid report type ID format", 400, err, nil)
	}

	reportType, err := s.reportTypeRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	response := ToReportTypeResponse(reportType)
	return &response, nil
}

func (s *service) GetReportTypeByName(ctx context.Context, name string) (*ReportTypeResponse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidReportTypeName
	}

	reportType, err := s.reportTypeRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	response := ToReportTypeResponse(reportType)
	return &response, nil
}

func (s *service) UpdateReportType(ctx context.Context, id string, req UpdateReportTypeRequest) (*ReportTypeResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_REPORT_TYPE_ID", "Invalid report type ID format", 400, err, nil)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidReportTypeName
	}

	reportType, err := s.reportTypeRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Check name uniqueness when being changed
	if name != reportType.Name {
		existingReportType, err := s.reportTypeRepo.GetByName(ctx, name)
		if err == nil && existingReportType != nil {
			return nil, ErrReportTypeAlreadyExists
		}
	}

	reportType.Name = name

	if err := s.reportTypeRepo.Update(ctx, objectID, reportType); err != nil {
		return nil, err
	}

	response := ToReportTypeResponse(reportType)
	return &response, nil
}

func (s *service) DeleteReportType(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("INVALID_REPORT_TYPE_ID", "Invalid report type ID format", 400, err, nil)
	}

	_, err = s.reportTypeRepo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}

	return s.reportTypeRepo.Delete(ctx, objectID)
}
