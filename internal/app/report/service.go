package report

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/utils/errors"
)

type Service interface {
	CreateReport(ctx context.Context, req CreateReportRequest) (*ReportResponse, error)
	UpdateReport(ctx context.Context, id string, req UpdateReportRequest) (*ReportResponse, error)
	DeleteReport(ctx context.Context, id string) error
	GetReports(ctx context.Context) ([]*ReportResponse, error)
	GetReportByID(ctx context.Context, id string) (*ReportResponse, error)
	GetReportByName(ctx context.Context, name string) (*ReportResponse, error)
	GetReportsByCompany(ctx context.Context, companyID string) ([]*ReportResponse, error)
	GetReportsByCompanies(ctx context.Context, req GetReportsByCompaniesRequest) ([]*ReportResponse, error)
	GetReportsByReportType(ctx context.Context, reportTypeID string) ([]*ReportResponse, error)
	GetReportsByUserAccess(ctx context.Context, userID string) ([]*ReportResponse, error)
	GetReportsByCreatedBy(ctx context.Context, userID string) ([]*ReportResponse, error)
}

type service struct {
	reportRepo domain.ReportRepository
}

func NewService(reportRepo domain.ReportRepository) Service {
	return &service{
		reportRepo: reportRepo,
	}
}

func (s *service) CreateReport(ctx context.Context, req CreateReportRequest) (*ReportResponse, error) {
	reportTypeID, err := primitive.ObjectIDFromHex(req.ReportType)
	if err != nil {
		return nil, errors.New("INVALID_REPORT_TYPE_ID", "Invalid report type ID format", 400, err, nil)
	}

	companyID, err := primitive.ObjectIDFromHex(req.Company)
	if err != nil {
		return nil, errors.New("INVALID_COMPANY_ID", "Invalid company ID format", 400, err, nil)
	}

	createdByID, err := primitive.ObjectIDFromHex(req.CreateBy)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid created by user ID format", 400, err, nil)
	}

	var userAccessIDs []primitive.ObjectID
	for _, userIDStr := range req.UserAccess {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return nil, errors.New("INVALID_USER_ACCESS_ID", "Invalid user access ID format", 400, err, nil)
		}
		userAccessIDs = append(userAccessIDs, userID)
	}

	// Default to empty array if no report data provided
	var reportData interface{}
	if req.ReportData != nil {
		reportData = req.ReportData
	} else {
		reportData = []interface{}{}
	}

	report := &domain.Report{
		ReportName: strings.TrimSpace(req.ReportName),
		ReportType: reportTypeID,
		Year:       strings.TrimSpace(req.Year),
		Company:    companyID,
		Currency:   req.Currency,
		CreatedBy:  createdByID,
		UserAccess: userAccessIDs,
		ReportData: reportData,
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, err
	}

	populatedReport, err := s.reportRepo.GetByID(ctx, report.ID)
	if err != nil {
		return nil, err
	}

	return ToReportResponse(populatedReport), nil
}

func (s *service) UpdateReport(ctx context.Context, id string, req UpdateReportRequest) (*ReportResponse, error) {
	reportID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_REPORT_ID", "Invalid report ID format", 400, err, nil)
	}

	existingReport, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return nil, err
	}

	// Prepare update data from existing report
	updateReport := &domain.Report{
		ID:         existingReport.ID,
		ReportName: existingReport.ReportName,
		ReportType: existingReport.ReportType.ID,
		Year:       existingReport.Year,
		Company:    existingReport.Company.ID,
		Currency:   existingReport.Currency,
		CreatedBy:  existingReport.CreatedBy.ID,
		UserAccess: []primitive.ObjectID{},
		ReportData: existingReport.ReportData,
		CreatedAt:  existingReport.CreatedAt,
	}

	// Convert populated user access back to ObjectIDs
	if existingReport.UserAccess != nil {
		for _, user := range existingReport.UserAccess {
			updateReport.UserAccess = append(updateReport.UserAccess, user.ID)
		}
	}

	if req.ReportName != nil {
		updateReport.ReportName = strings.TrimSpace(*req.ReportName)
	}

	if req.ReportType != nil {
		reportTypeID, err := primitive.ObjectIDFromHex(*req.ReportType)
		if err != nil {
			return nil, errors.New("INVALID_REPORT_TYPE_ID", "Invalid report type ID format", 400, err, nil)
		}
		updateReport.ReportType = reportTypeID
	}

	if req.Year != nil {
		updateReport.Year = strings.TrimSpace(*req.Year)
	}

	if req.Company != nil {
		companyID, err := primitive.ObjectIDFromHex(*req.Company)
		if err != nil {
			return nil, errors.New("INVALID_COMPANY_ID", "Invalid company ID format", 400, err, nil)
		}
		updateReport.Company = companyID
	}

	if req.Currency != nil {
		updateReport.Currency = req.Currency
	}

	if req.UserAccess != nil {
		var userAccessIDs []primitive.ObjectID
		for _, userIDStr := range req.UserAccess {
			userID, err := primitive.ObjectIDFromHex(userIDStr)
			if err != nil {
				return nil, errors.New("INVALID_USER_ACCESS_ID", "Invalid user access ID format", 400, err, nil)
			}
			userAccessIDs = append(userAccessIDs, userID)
		}
		updateReport.UserAccess = userAccessIDs
	}

	if req.ReportData != nil {
		updateReport.ReportData = req.ReportData
	}

	updatedReport, err := s.reportRepo.Update(ctx, reportID, updateReport)
	if err != nil {
		return nil, err
	}

	return ToReportResponse(updatedReport), nil
}

func (s *service) DeleteReport(ctx context.Context, id string) error {
	reportID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("INVALID_REPORT_ID", "Invalid report ID format", 400, err, nil)
	}

	return s.reportRepo.Delete(ctx, reportID)
}

func (s *service) GetReports(ctx context.Context) ([]*ReportResponse, error) {
	reports, err := s.reportRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return ToReportResponseArray(reports), nil
}

func (s *service) GetReportByID(ctx context.Context, id string) (*ReportResponse, error) {
	reportID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_REPORT_ID", "Invalid report ID format", 400, err, nil)
	}

	report, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return nil, err
	}

	return ToReportResponse(report), nil
}

func (s *service) GetReportByName(ctx context.Context, name string) (*ReportResponse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("INVALID_REPORT_NAME", "Report name cannot be empty", 400, nil, nil)
	}

	report, err := s.reportRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return ToReportResponse(report), nil
}

func (s *service) GetReportsByCompany(ctx context.Context, companyID string) ([]*ReportResponse, error) {
	companyObjID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New("INVALID_COMPANY_ID", "Invalid company ID format", 400, err, nil)
	}

	reports, err := s.reportRepo.GetByCompany(ctx, companyObjID)
	if err != nil {
		return nil, err
	}

	return ToReportResponseArray(reports), nil
}

func (s *service) GetReportsByCompanies(ctx context.Context, req GetReportsByCompaniesRequest) ([]*ReportResponse, error) {
	// Business rule: comparison requires at least 2 companies
	if len(req.CompanyIds) < 2 {
		return nil, errors.New("INSUFFICIENT_COMPANIES", "Need 2 or more companies", 400, nil, nil)
	}

	var companyIDs []primitive.ObjectID
	for _, companyIDStr := range req.CompanyIds {
		companyID, err := primitive.ObjectIDFromHex(companyIDStr)
		if err != nil {
			return nil, errors.New("INVALID_COMPANY_ID", "Invalid company ID format", 400, err, nil)
		}
		companyIDs = append(companyIDs, companyID)
	}

	reports, err := s.reportRepo.GetByCompanies(ctx, companyIDs)
	if err != nil {
		return nil, err
	}

	return ToReportResponseArray(reports), nil
}

func (s *service) GetReportsByReportType(ctx context.Context, reportTypeID string) ([]*ReportResponse, error) {
	reportTypeObjID, err := primitive.ObjectIDFromHex(reportTypeID)
	if err != nil {
		return nil, errors.New("INVALID_REPORT_TYPE_ID", "Invalid report type ID format", 400, err, nil)
	}

	reports, err := s.reportRepo.GetByReportType(ctx, reportTypeObjID)
	if err != nil {
		return nil, err
	}

	return ToReportResponseArray(reports), nil
}

func (s *service) GetReportsByUserAccess(ctx context.Context, userID string) ([]*ReportResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid user ID format", 400, err, nil)
	}

	reports, err := s.reportRepo.GetByUserAccess(ctx, userObjID)
	if err != nil {
		return nil, err
	}

	return ToReportResponseArray(reports), nil
}

func (s *service) GetReportsByCreatedBy(ctx context.Context, userID string) ([]*ReportResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid user ID format", 400, err, nil)
	}

	reports, err := s.reportRepo.GetByCreatedBy(ctx, userObjID)
	if err != nil {
		return nil, err
	}

	return ToReportResponseArray(reports), nil
}
