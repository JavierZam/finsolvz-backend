package report

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
)

// Mock repository for testing
type mockReportRepository struct {
	reports []domain.PopulatedReport
}

func (m *mockReportRepository) Create(ctx context.Context, report *domain.Report) error {
	report.ID = primitive.NewObjectID()
	return nil
}

func (m *mockReportRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.PopulatedReport, error) {
	for _, r := range m.reports {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, nil
}

func (m *mockReportRepository) GetByName(ctx context.Context, name string) (*domain.PopulatedReport, error) {
	return &m.reports[0], nil
}

func (m *mockReportRepository) GetAll(ctx context.Context) ([]*domain.PopulatedReport, error) {
	var result []*domain.PopulatedReport
	for i := range m.reports {
		result = append(result, &m.reports[i])
	}
	return result, nil
}

func (m *mockReportRepository) GetAllPaginated(ctx context.Context, skip, limit int) ([]*domain.PopulatedReport, int, error) {
	total := len(m.reports)
	end := skip + limit
	if end > total {
		end = total
	}

	var result []*domain.PopulatedReport
	if skip < total {
		for i := skip; i < end; i++ {
			result = append(result, &m.reports[i])
		}
	}

	return result, total, nil
}

func (m *mockReportRepository) GetByCompany(ctx context.Context, companyID primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	return []*domain.PopulatedReport{&m.reports[0]}, nil
}

func (m *mockReportRepository) GetByCompanies(ctx context.Context, companyIDs []primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	return []*domain.PopulatedReport{&m.reports[0]}, nil
}

func (m *mockReportRepository) GetByReportType(ctx context.Context, reportTypeID primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	return []*domain.PopulatedReport{&m.reports[0]}, nil
}

func (m *mockReportRepository) GetByUserAccess(ctx context.Context, userID primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	return []*domain.PopulatedReport{&m.reports[0]}, nil
}

func (m *mockReportRepository) GetByCreatedBy(ctx context.Context, userID primitive.ObjectID) ([]*domain.PopulatedReport, error) {
	return []*domain.PopulatedReport{&m.reports[0]}, nil
}

func (m *mockReportRepository) Update(ctx context.Context, id primitive.ObjectID, report *domain.Report) (*domain.PopulatedReport, error) {
	return &m.reports[0], nil
}

func (m *mockReportRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func TestService_GetReportsPaginated(t *testing.T) {
	// Setup mock data
	mockRepo := &mockReportRepository{
		reports: []domain.PopulatedReport{
			{
				ID:         primitive.NewObjectID(),
				ReportName: "Test Report 1",
				Year:       2024,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				ID:         primitive.NewObjectID(),
				ReportName: "Test Report 2",
				Year:       2024,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}

	service := NewService(mockRepo)

	// Test pagination
	reports, total, err := service.GetReportsPaginated(context.Background(), 0, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(reports) != 1 {
		t.Fatalf("Expected 1 report, got %d", len(reports))
	}

	if total != 2 {
		t.Fatalf("Expected total 2, got %d", total)
	}
}

func TestService_GetReportByID_Performance(t *testing.T) {
	mockRepo := &mockReportRepository{
		reports: []domain.PopulatedReport{
			{
				ID:         primitive.NewObjectID(),
				ReportName: "Performance Test Report",
				Year:       2024,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}

	service := NewService(mockRepo)
	reportID := mockRepo.reports[0].ID.Hex()

	// Measure performance
	start := time.Now()
	_, err := service.GetReportByID(context.Background(), reportID)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should complete within 100ms
	if duration > 100*time.Millisecond {
		t.Fatalf("GetReportByID took too long: %v", duration)
	}

	// Test cache hit
	start = time.Now()
	_, err = service.GetReportByID(context.Background(), reportID)
	cachedDuration := time.Since(start)

	if err != nil {
		t.Fatalf("Expected no error on cached request, got %v", err)
	}

	// Cached request should be much faster
	if cachedDuration > 10*time.Millisecond {
		t.Fatalf("Cached request took too long: %v", cachedDuration)
	}
}
