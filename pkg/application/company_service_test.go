package application_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jizumer/expedition-value/pkg/application"
	"github.com/jizumer/expedition-value/pkg/domain/company"
	// "github.com/stretchr/testify/assert" // Example: using testify for assertions
)

// --- Mock CompanyRepository ---

type MockCompanyRepository struct {
	FindByTickerFunc         func(ticker string) (*company.Company, error)
	SearchByScoreRangeFunc   func(minScore, maxScore float64) ([]*company.Company, error)
	SaveFunc                 func(c *company.Company) error
	DeleteFunc               func(ticker string) error
	// Optional methods if needed for other tests
	// FindAllFunc           func() ([]*company.Company, error)
	// FindBySectorFunc      func(sector company.Sector) ([]*company.Company, error)

	// Spy fields (optional, to check if methods were called)
	SaveCalledWith   *company.Company
	FindByTickerCalledWith string
}

func (m *MockCompanyRepository) FindByTicker(ticker string) (*company.Company, error) {
	m.FindByTickerCalledWith = ticker
	if m.FindByTickerFunc != nil {
		return m.FindByTickerFunc(ticker)
	}
	return nil, errors.New("FindByTickerFunc not implemented in mock")
}

func (m *MockCompanyRepository) SearchByScoreRange(minScore, maxScore float64) ([]*company.Company, error) {
	if m.SearchByScoreRangeFunc != nil {
		return m.SearchByScoreRangeFunc(minScore, maxScore)
	}
	return nil, errors.New("SearchByScoreRangeFunc not implemented in mock")
}

func (m *MockCompanyRepository) Save(c *company.Company) error {
	m.SaveCalledWith = c // Spy on the argument
	if m.SaveFunc != nil {
		return m.SaveFunc(c)
	}
	return errors.New("SaveFunc not implemented in mock")
}

func (m *MockCompanyRepository) Delete(ticker string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ticker)
	}
	return errors.New("DeleteFunc not implemented in mock")
}

// --- CompanyService Tests ---

func TestCompanyService_GetCompanyByTicker(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	service := application.NewCompanyService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedCompany, _ := company.NewCompany("AAPL", company.FinancialMetrics{PERatio: 15}, company.Technology)
		mockRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			if ticker == "AAPL" {
				return expectedCompany, nil
			}
			return nil, errors.New("company not found")
		}

		c, err := service.GetCompanyByTicker("AAPL")

		if err != nil {
			t.Errorf("GetCompanyByTicker() error = %v, wantErr nil", err)
		}
		if c == nil {
			t.Fatalf("GetCompanyByTicker() company = nil, want non-nil")
		}
		if c.Ticker != "AAPL" {
			t.Errorf("GetCompanyByTicker() ticker = %s, want AAPL", c.Ticker)
		}
		// Using testify/assert: assert.Equal(t, "AAPL", c.Ticker)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			return nil, errors.New("company not found") // Simulate repository error
		}

		_, err := service.GetCompanyByTicker("UNKNOWN")

		if err == nil {
			t.Errorf("GetCompanyByTicker() with unknown ticker expected error, got nil")
		}
	})

	t.Run("EmptyTicker", func(t *testing.T) {
		_, err := service.GetCompanyByTicker("")
		if err == nil {
			t.Errorf("GetCompanyByTicker() with empty ticker expected error, got nil")
		}
	})
}

func TestCompanyService_CreateCompany(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	service := application.NewCompanyService(mockRepo)

	validMetrics, _ := company.NewFinancialMetrics(20, 3, 0.6) // Reusable valid metrics

	t.Run("Success", func(t *testing.T) {
		ticker := "MSFT"
		sector := company.Technology
		mockRepo.SaveCalledWith = nil // Reset spy

		mockRepo.SaveFunc = func(c *company.Company) error {
			mockRepo.SaveCalledWith = c // Capture the company passed to Save
			return nil
		}

		createdCompany, err := service.CreateCompany(ticker, *validMetrics, sector)

		if err != nil {
			t.Fatalf("CreateCompany() error = %v, wantErr nil", err)
		}
		if createdCompany == nil {
			t.Fatalf("CreateCompany() returned nil company, want non-nil")
		}
		if createdCompany.Ticker != ticker {
			t.Errorf("CreateCompany() Ticker = %s, want %s", createdCompany.Ticker, ticker)
		}
		if createdCompany.FinancialMetrics.PERatio != validMetrics.PERatio {
			t.Errorf("CreateCompany() PERatio = %v, want %v", createdCompany.FinancialMetrics.PERatio, validMetrics.PERatio)
		}
		if createdCompany.Sector != sector {
			t.Errorf("CreateCompany() Sector = %v, want %v", createdCompany.Sector, sector)
		}

		if mockRepo.SaveCalledWith == nil {
			t.Errorf("SaveFunc was not called")
		} else {
			if mockRepo.SaveCalledWith.Ticker != ticker {
				t.Errorf("SaveFunc called with Ticker = %s, want %s", mockRepo.SaveCalledWith.Ticker, ticker)
			}
			if mockRepo.SaveCalledWith.FinancialMetrics.PERatio != validMetrics.PERatio {
				t.Errorf("SaveFunc called with PERatio = %v, want %v", mockRepo.SaveCalledWith.FinancialMetrics.PERatio, validMetrics.PERatio)
			}
			if mockRepo.SaveCalledWith.Sector != sector {
				t.Errorf("SaveFunc called with Sector = %v, want %v", mockRepo.SaveCalledWith.Sector, sector)
			}
		}
	})

	t.Run("SuccessWithDefaultInputs", func(t *testing.T) {
		ticker := "DEFAULT"
		defaultMetrics := company.FinancialMetrics{} // As used by handler
		defaultSector := company.UndefinedSector   // As used by handler
		mockRepo.SaveCalledWith = nil // Reset spy

		mockRepo.SaveFunc = func(c *company.Company) error {
			mockRepo.SaveCalledWith = c
			return nil
		}

		// Assuming company.NewCompany handles UndefinedSector gracefully (e.g., it's a valid defined value in the enum)
		// And that FinancialMetrics{} is also valid for NewCompany
		createdCompany, err := service.CreateCompany(ticker, defaultMetrics, defaultSector)

		if err != nil {
			t.Fatalf("CreateCompany(default inputs) error = %v, wantErr nil", err)
		}
		if createdCompany == nil {
			t.Fatalf("CreateCompany(default inputs) returned nil company, want non-nil")
		}
		if createdCompany.Ticker != ticker {
			t.Errorf("CreateCompany(default inputs) Ticker = %s, want %s", createdCompany.Ticker, ticker)
		}
		if createdCompany.Sector != defaultSector {
			t.Errorf("CreateCompany(default inputs) Sector = %v, want %v", createdCompany.Sector, defaultSector)
		}
		// FinancialMetrics will have zero values, CurrentScore will be 0 initially
		if createdCompany.FinancialMetrics.PERatio != 0 {
			t.Errorf("CreateCompany(default inputs) PERatio = %v, want 0", createdCompany.FinancialMetrics.PERatio)
		}

		if mockRepo.SaveCalledWith == nil {
			t.Errorf("SaveFunc was not called for default inputs")
		} else if mockRepo.SaveCalledWith.Ticker != ticker {
			t.Errorf("SaveFunc called with Ticker = %s for default inputs, want %s", mockRepo.SaveCalledWith.Ticker, ticker)
		}
	})


	t.Run("EmptyTickerDomainError", func(t *testing.T) {
		_, err := service.CreateCompany("", *validMetrics, company.Technology)
		if err == nil {
			t.Errorf("CreateCompany() with empty ticker expected domain error, got nil")
		}
		// Example of more specific error check:
		// if !strings.Contains(err.Error(), "ticker cannot be empty") {
		// 	t.Errorf("Expected error about empty ticker, got: %v", err)
		// }
	})

	t.Run("RepositorySaveError_AlreadyExists", func(t *testing.T) {
		expectedErr := errors.New("company already exists") // Simulate specific error
		mockRepo.SaveFunc = func(c *company.Company) error {
			return expectedErr
		}
		_, err := service.CreateCompany("TSLA", *validMetrics, company.Technology)
		if err == nil {
			t.Errorf("CreateCompany() expected repository save error (already exists), got nil")
		}
		if !errors.Is(err, expectedErr) { // Check if the error is the one we expect
			t.Errorf("CreateCompany() error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("RepositorySaveError_Generic", func(t *testing.T) {
		expectedErr := errors.New("failed to save company for other reasons")
		mockRepo.SaveFunc = func(c *company.Company) error {
			return expectedErr
		}
		_, err := service.CreateCompany("NVDA", *validMetrics, company.Technology)
		if err == nil {
			t.Errorf("CreateCompany() expected generic repository save error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("CreateCompany() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestCompanyService_SearchCompaniesByScore(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	service := application.NewCompanyService(mockRepo)

	comp1, _ := company.NewCompany("C1", company.FinancialMetrics{PERatio: 10}, company.Technology)
	comp1.CurrentScore = 70
	comp2, _ := company.NewCompany("C2", company.FinancialMetrics{PERatio: 12}, company.Industrials)
	comp2.CurrentScore = 85

	t.Run("Success", func(t *testing.T) {
		expectedCompanies := []*company.Company{comp1, comp2}
		mockRepo.SearchByScoreRangeFunc = func(minScore, maxScore float64) ([]*company.Company, error) {
			if minScore == 60 && maxScore == 90 {
				return expectedCompanies, nil
			}
			return nil, errors.New("unexpected score range")
		}

		results, err := service.SearchCompaniesByScore(60, 90)
		if err != nil {
			t.Fatalf("SearchCompaniesByScore() error = %v, wantErr nil", err)
		}
		if len(results) != 2 {
			t.Fatalf("SearchCompaniesByScore() len = %d, want 2", len(results))
		}
		// Using testify/assert: assert.ElementsMatch(t, expectedCompanies, results)
	})

	t.Run("NoResults", func(t *testing.T) {
		mockRepo.SearchByScoreRangeFunc = func(minScore, maxScore float64) ([]*company.Company, error) {
			return []*company.Company{}, nil // Empty slice
		}
		results, err := service.SearchCompaniesByScore(100, 110)
		if err != nil {
			t.Fatalf("SearchCompaniesByScore() for no results error = %v, wantErr nil", err)
		}
		if len(results) != 0 {
			t.Errorf("SearchCompaniesByScore() for no results len = %d, want 0", len(results))
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.SearchByScoreRangeFunc = func(minScore, maxScore float64) ([]*company.Company, error) {
			return nil, errors.New("database error")
		}
		_, err := service.SearchCompaniesByScore(10, 20)
		if err == nil {
			t.Errorf("SearchCompaniesByScore() expected repository error, got nil")
		}
	})

	t.Run("InvalidScoreRange", func(t *testing.T) {
		_, err := service.SearchCompaniesByScore(90, 60) // min > max
		if err == nil {
			t.Errorf("SearchCompaniesByScore() with min > max expected error, got nil")
		}
	})
}

func TestCompanyService_UpdateCompanyMetrics(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	service := application.NewCompanyService(mockRepo)

	existingMetrics, _ := company.NewFinancialMetrics(10, 1, 0.5)
	existingCompany, _ := company.NewCompany("EXT", *existingMetrics, company.Technology)
	
	newMetricsData, _ := company.NewFinancialMetrics(15, 1.5, 0.55)

	t.Run("Success", func(t *testing.T) {
		// Reset spy field for each sub-test if needed, or ensure distinct mock instances
		mockRepo.SaveCalledWith = nil 
		
		mockRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			if ticker == "EXT" {
				// Return a fresh copy to avoid state leakage between tests if the object is modified directly
				clone, _ := company.NewCompany("EXT", *existingMetrics, company.Technology)
				clone.UpdatedAt = existingCompany.UpdatedAt // Preserve original update time for comparison
				return clone, nil
			}
			return nil, errors.New("not found")
		}
		mockRepo.SaveFunc = func(c *company.Company) error {
			mockRepo.SaveCalledWith = c
			return nil
		}

		err := service.UpdateCompanyMetrics("EXT", *newMetricsData)
		if err != nil {
			t.Fatalf("UpdateCompanyMetrics() error = %v, wantErr nil", err)
		}
		if mockRepo.SaveCalledWith == nil {
			t.Fatalf("Save was not called on repository")
		}
		if mockRepo.SaveCalledWith.Ticker != "EXT" {
			t.Errorf("Saved company ticker = %s, want EXT", mockRepo.SaveCalledWith.Ticker)
		}
		if mockRepo.SaveCalledWith.FinancialMetrics.PERatio != newMetricsData.PERatio {
			t.Errorf("Saved company PERatio = %v, want %v", mockRepo.SaveCalledWith.FinancialMetrics.PERatio, newMetricsData.PERatio)
		}
		// Check if UpdatedAt timestamps were advanced (FinancialMetrics.MetricsUpdatedAt and Company.UpdatedAt)
		// This requires comparing with the state *before* the UpdateFinancialMetrics call in the domain object.
		if mockRepo.SaveCalledWith.FinancialMetrics.MetricsUpdatedAt.Before(existingCompany.UpdatedAt) {
			t.Error("FinancialMetrics.MetricsUpdatedAt was not advanced or set correctly")
		}
		if mockRepo.SaveCalledWith.UpdatedAt.Before(existingCompany.UpdatedAt) || mockRepo.SaveCalledWith.UpdatedAt.Equal(existingCompany.UpdatedAt) {
			t.Errorf("Company.UpdatedAt was not advanced. Before: %v, After: %v", existingCompany.UpdatedAt, mockRepo.SaveCalledWith.UpdatedAt)
		}
	})

	t.Run("CompanyNotFound", func(t *testing.T) {
		mockRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			return nil, errors.New("company not found")
		}
		err := service.UpdateCompanyMetrics("UNKNOWN", *newMetricsData)
		if err == nil {
			t.Errorf("UpdateCompanyMetrics() for unknown company expected error, got nil")
		}
	})

	t.Run("EmptyTicker", func(t *testing.T) {
		err := service.UpdateCompanyMetrics("", *newMetricsData)
		if err == nil {
			t.Errorf("UpdateCompanyMetrics() with empty ticker expected error, got nil")
		}
	})

	t.Run("SaveError", func(t *testing.T) {
		mockRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			return existingCompany, nil
		}
		mockRepo.SaveFunc = func(c *company.Company) error {
			return errors.New("failed to save")
		}
		err := service.UpdateCompanyMetrics("EXT", *newMetricsData)
		if err == nil {
			t.Errorf("UpdateCompanyMetrics() expected save error, got nil")
		}
	})
}

func TestCompanyService_RefreshCompany(t *testing.T) {
	mockRepo := &MockCompanyRepository{}
	service := application.NewCompanyService(mockRepo)

	// Company with stale metrics
	staleMetrics, _ := company.NewFinancialMetrics(10,1,1)
	staleMetrics.MetricsUpdatedAt = time.Now().Add(-10 * 24 * time.Hour) // 10 days old
	staleCompany, _ := company.NewCompany("STALE", *staleMetrics, company.Technology)
	originalStaleCompanyUpdateTime := staleCompany.UpdatedAt
	originalStaleMetricsUpdateTime := staleCompany.FinancialMetrics.MetricsUpdatedAt
	
	t.Run("Success_StaleMetrics", func(t *testing.T) {
		mockRepo.SaveCalledWith = nil
		mockRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			if ticker == "STALE" {
				// Return a fresh copy of the stale company for the test
				clone, _ := company.NewCompany("STALE", *staleMetrics, company.Technology)
				clone.UpdatedAt = originalStaleCompanyUpdateTime
				clone.FinancialMetrics.MetricsUpdatedAt = originalStaleMetricsUpdateTime
				return clone, nil
			}
			return nil, errors.New("not found")
		}
		mockRepo.SaveFunc = func(c *company.Company) error {
			mockRepo.SaveCalledWith = c
			return nil
		}

		err := service.RefreshCompany("STALE")
		if err != nil {
			t.Fatalf("RefreshCompany() for stale metrics error = %v, wantErr nil", err)
		}
		if mockRepo.SaveCalledWith == nil {
			t.Fatalf("Save was not called on repository for stale metrics refresh")
		}
		// Check that metrics and company UpdatedAt timestamps were advanced by the domain logic
		if mockRepo.SaveCalledWith.FinancialMetrics.MetricsUpdatedAt.Equal(originalStaleMetricsUpdateTime) ||
		   mockRepo.SaveCalledWith.FinancialMetrics.MetricsUpdatedAt.Before(originalStaleMetricsUpdateTime) {
			t.Errorf("FinancialMetrics.MetricsUpdatedAt not advanced after refresh. Original: %v, Current: %v",
				originalStaleMetricsUpdateTime, mockRepo.SaveCalledWith.FinancialMetrics.MetricsUpdatedAt)
		}
		if mockRepo.SaveCalledWith.UpdatedAt.Equal(originalStaleCompanyUpdateTime) ||
		   mockRepo.SaveCalledWith.UpdatedAt.Before(originalStaleCompanyUpdateTime) {
			t.Errorf("Company.UpdatedAt not advanced after refresh. Original: %v, Current: %v",
				originalStaleCompanyUpdateTime, mockRepo.SaveCalledWith.UpdatedAt)
		}
	})
	
	// Company with recent metrics
	recentMetrics, _ := company.NewFinancialMetrics(12,1.2,0.6)
	recentMetrics.MetricsUpdatedAt = time.Now().Add(-1 * 24 * time.Hour) // 1 day old
	recentCompany, _ := company.NewCompany("RECENT", *recentMetrics, company.Technology)
	originalRecentCompanyUpdateTime := recentCompany.UpdatedAt
	originalRecentMetricsUpdateTime := recentCompany.FinancialMetrics.MetricsUpdatedAt

	t.Run("Success_RecentMetrics", func(t *testing.T) {
		mockRepo.SaveCalledWith = nil
		mockRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			if ticker == "RECENT" {
				clone, _ := company.NewCompany("RECENT", *recentMetrics, company.Technology)
				clone.UpdatedAt = originalRecentCompanyUpdateTime
				clone.FinancialMetrics.MetricsUpdatedAt = originalRecentMetricsUpdateTime
				return clone, nil
			}
			return nil, errors.New("not found")
		}
		mockRepo.SaveFunc = func(c *company.Company) error {
			mockRepo.SaveCalledWith = c
			return nil
		}
		
		err := service.RefreshCompany("RECENT")
		if err != nil {
			t.Fatalf("RefreshCompany() for recent metrics error = %v, wantErr nil", err)
		}
		if mockRepo.SaveCalledWith == nil {
			t.Fatalf("Save was not called on repository for recent metrics refresh")
		}
		// For recent metrics, domain logic placeholder for RefreshStaleMetrics does not update timestamps
		// So, timestamps in SaveCalledWith should be the same as original ones.
		if !mockRepo.SaveCalledWith.FinancialMetrics.MetricsUpdatedAt.Equal(originalRecentMetricsUpdateTime) {
			t.Errorf("FinancialMetrics.MetricsUpdatedAt changed for recent metrics. Original: %v, Current: %v",
				originalRecentMetricsUpdateTime, mockRepo.SaveCalledWith.FinancialMetrics.MetricsUpdatedAt)
		}
		// The current domain placeholder for RefreshStaleMetrics (company.go) doesn't update company.UpdatedAt
		// if metrics are not stale. If it did, this test would need to expect a change.
		if !mockRepo.SaveCalledWith.UpdatedAt.Equal(originalRecentCompanyUpdateTime) {
			t.Errorf("Company.UpdatedAt changed for recent metrics when no refresh occurred. Original: %v, Current: %v",
				originalRecentCompanyUpdateTime, mockRepo.SaveCalledWith.UpdatedAt)
		}
	})


	t.Run("CompanyNotFound", func(t *testing.T) {
		mockRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			return nil, errors.New("company not found")
		}
		err := service.RefreshCompany("UNKNOWN")
		if err == nil {
			t.Errorf("RefreshCompany() for unknown company expected error, got nil")
		}
	})

	t.Run("EmptyTicker", func(t *testing.T) {
		err := service.RefreshCompany("")
		if err == nil {
			t.Errorf("RefreshCompany() with empty ticker expected error, got nil")
		}
	})
}
