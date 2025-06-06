package application_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jizumer/expedition-value/pkg/application"
	"github.com/jizumer/expedition-value/pkg/domain/company"
	"github.com/jizumer/expedition-value/pkg/domain/portfolio"
	// "github.com/stretchr/testify/assert"
)

// --- Mock PortfolioRepository ---
type MockPortfolioRepository struct {
	FindByIDFunc            func(id string) (*portfolio.Portfolio, error)
	FindAllFunc             func() ([]*portfolio.Portfolio, error)
	SearchByRiskProfileFunc func(riskProfile portfolio.RiskProfile) ([]*portfolio.Portfolio, error)
	SearchBySectorFunc      func(sector company.Sector) ([]*portfolio.Portfolio, error) // Added
	SaveFunc                func(p *portfolio.Portfolio) error
	DeleteFunc              func(id string) error

	SaveCalledWith *portfolio.Portfolio
}

func (m *MockPortfolioRepository) FindByID(id string) (*portfolio.Portfolio, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc not implemented in mock")
}

func (m *MockPortfolioRepository) FindAll() ([]*portfolio.Portfolio, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, errors.New("FindAllFunc not implemented in mock")
}

func (m *MockPortfolioRepository) SearchByRiskProfile(riskProfile portfolio.RiskProfile) ([]*portfolio.Portfolio, error) {
	if m.SearchByRiskProfileFunc != nil {
		return m.SearchByRiskProfileFunc(riskProfile)
	}
	return nil, errors.New("SearchByRiskProfileFunc not implemented in mock")
}

// SearchBySector is part of the interface, so it needs to be on the mock
func (m *MockPortfolioRepository) SearchBySector(sector company.Sector) ([]*portfolio.Portfolio, error) {
	if m.SearchBySectorFunc != nil {
		return m.SearchBySectorFunc(sector)
	}
	return nil, errors.New("SearchBySectorFunc not implemented in mock")
}

func (m *MockPortfolioRepository) Save(p *portfolio.Portfolio) error {
	m.SaveCalledWith = p
	if m.SaveFunc != nil {
		return m.SaveFunc(p)
	}
	return errors.New("SaveFunc not implemented in mock")
}

func (m *MockPortfolioRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return errors.New("DeleteFunc not implemented in mock")
}

// MockCompanyRepository is already defined in company_service_test.go
// We can reuse it if it's in the same package `application_test`, or redefine/copy if needed.
// For this exercise, assuming it can be implicitly reused or we'd copy it.
// If these test files are compiled into the same test binary (common for `go test ./...`),
// then MockCompanyRepository from company_service_test.go might be accessible if not for naming conflicts.
// To be safe and explicit, especially if running tests per-package, we'll define it here.
type MinimalMockCompanyRepository struct {
	FindByTickerFunc func(ticker string) (*company.Company, error)
}

func (m *MinimalMockCompanyRepository) FindByTicker(ticker string) (*company.Company, error) {
	if m.FindByTickerFunc != nil {
		return m.FindByTickerFunc(ticker)
	}
	return nil, errors.New("FindByTickerFunc not implemented in minimal mock company repo")
}
func (m *MinimalMockCompanyRepository) SearchByScoreRange(minScore, maxScore float64) ([]*company.Company, error) { return nil, nil }
func (m *MinimalMockCompanyRepository) Save(c *company.Company) error { return nil }
func (m *MinimalMockCompanyRepository) Delete(ticker string) error    { return nil }


// --- PortfolioService Tests ---

func TestPortfolioService_CreatePortfolio(t *testing.T) {
	mockPortfolioRepo := &MockPortfolioRepository{}
	// CompanyRepo is not strictly needed for CreatePortfolio, can be nil or a minimal mock
	mockCompanyRepo := &MinimalMockCompanyRepository{}
	service := application.NewPortfolioService(mockPortfolioRepo, mockCompanyRepo)

	cash, _ := portfolio.NewMoney(500000, "USD") // 5000.00 USD
	risk := portfolio.Conservative

	t.Run("Success", func(t *testing.T) {
		mockPortfolioRepo.SaveFunc = func(p *portfolio.Portfolio) error {
			return nil // Simulate successful save
		}

		p, err := service.CreatePortfolio(*cash, risk)

		if err != nil {
			t.Fatalf("CreatePortfolio() error = %v, wantErr nil", err)
		}
		if p == nil {
			t.Fatalf("CreatePortfolio() portfolio = nil, want non-nil")
		}
		if p.ID == "" { // UUID should be generated
			t.Errorf("CreatePortfolio() ID is empty, want a generated UUID")
		}
		if p.CashBalance.Amount != cash.Amount {
			t.Errorf("CreatePortfolio() CashBalance = %d, want %d", p.CashBalance.Amount, cash.Amount)
		}
		if p.RiskProfile != risk {
			t.Errorf("CreatePortfolio() RiskProfile = %v, want %v", p.RiskProfile, risk)
		}
		if mockPortfolioRepo.SaveCalledWith == nil {
			t.Errorf("SaveFunc was not called on portfolio repository")
		} else if mockPortfolioRepo.SaveCalledWith.ID != p.ID {
			t.Errorf("SaveFunc called with ID = %s, want %s", mockPortfolioRepo.SaveCalledWith.ID, p.ID)
		}
	})

	t.Run("DomainValidationError", func(t *testing.T) {
		invalidCash, _ := portfolio.NewMoney(-100, "USD") // Negative cash
		_, err := service.CreatePortfolio(*invalidCash, risk)
		if err == nil {
			t.Errorf("CreatePortfolio() with invalid domain data expected error, got nil")
		}
		// Expected error message: "failed to create new portfolio in domain: initial cash balance cannot be negative"
	})

	t.Run("RepositorySaveError", func(t *testing.T) {
		mockPortfolioRepo.SaveFunc = func(p *portfolio.Portfolio) error {
			return errors.New("database constraint failed")
		}
		_, err := service.CreatePortfolio(*cash, risk)
		if err == nil {
			t.Errorf("CreatePortfolio() expected repository save error, got nil")
		}
		// Expected error message: "failed to save portfolio: database constraint failed"
	})
}

func TestPortfolioService_GetPortfolioDetails(t *testing.T) {
	mockPortfolioRepo := &MockPortfolioRepository{}
	mockCompanyRepo := &MinimalMockCompanyRepository{} // Not used in this method
	service := application.NewPortfolioService(mockPortfolioRepo, mockCompanyRepo)

	portfolioID := uuid.NewString()
	expectedPortfolio, _ := portfolio.NewPortfolio(portfolioID, portfolio.Moderate, portfolio.Money{Amount: 1000, Currency: "USD"})

	t.Run("Success", func(t *testing.T) {
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			if id == portfolioID {
				return expectedPortfolio, nil
			}
			return nil, errors.New("portfolio not found in mock")
		}

		p, err := service.GetPortfolioDetails(portfolioID)
		if err != nil {
			t.Fatalf("GetPortfolioDetails() error = %v, wantErr nil", err)
		}
		if p == nil {
			t.Fatalf("GetPortfolioDetails() portfolio = nil, want non-nil")
		}
		if p.ID != portfolioID {
			t.Errorf("GetPortfolioDetails() ID = %s, want %s", p.ID, portfolioID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			return nil, errors.New("db: no rows in result set") // Simulate repo error
		}
		_, err := service.GetPortfolioDetails(uuid.NewString())
		if err == nil {
			t.Errorf("GetPortfolioDetails() for non-existent ID expected error, got nil")
		}
		// Expected: "failed to find portfolio ...: db: no rows in result set"
	})

	t.Run("NotFound_RepoReturnsNilNil", func(t *testing.T) {
		// Test the service's specific nil check after repository call
		nonExistentID := uuid.NewString()
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			if id == nonExistentID {
				return nil, nil // Simulate repository returning no error but also no portfolio
			}
			return nil, errors.New("unexpected ID in mock")
		}
		_, err := service.GetPortfolioDetails(nonExistentID)
		if err == nil {
			t.Errorf("GetPortfolioDetails() for non-existent ID (repo nil,nil) expected error, got nil")
		}
		expectedErrorMsg := "portfolio " + nonExistentID + " not found"
		if err != nil && err.Error() != expectedErrorMsg {
			t.Errorf("GetPortfolioDetails() error = %q, want %q", err.Error(), expectedErrorMsg)
		}
	})

	t.Run("EmptyID", func(t *testing.T) {
		_, err := service.GetPortfolioDetails("")
		if err == nil {
			t.Errorf("GetPortfolioDetails() with empty ID expected error, got nil")
		}
		// Expected: "portfolioID cannot be empty"
	})
}

func TestPortfolioService_AddPosition(t *testing.T) {
	mockPortfolioRepo := &MockPortfolioRepository{}
	mockCompanyRepo := &MinimalMockCompanyRepository{}
	service := application.NewPortfolioService(mockPortfolioRepo, mockCompanyRepo)

	portfolioID := uuid.NewString()
	// Adjusted initialCash to be sufficient for the test position
	initialCash, _ := portfolio.NewMoney(200000, "USD") // 2000.00 
	existingPortfolio, _ := portfolio.NewPortfolio(portfolioID, portfolio.Aggressive, *initialCash)

	companyTicker := "AAPL"
	shares := 10
	purchasePrice, _ := portfolio.NewMoney(15000, "USD") // 150.00 per share

	sampleCompany, _ := company.NewCompany(companyTicker, company.FinancialMetrics{}, company.Technology)

	t.Run("Success", func(t *testing.T) {
		// Reset state for this sub-test
		freshPortfolio, _ := portfolio.NewPortfolio(portfolioID, portfolio.Aggressive, *initialCash)
		mockPortfolioRepo.SaveCalledWith = nil

		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			if id == portfolioID {
				return freshPortfolio, nil // Return the modifiable portfolio
			}
			return nil, errors.New("portfolio not found")
		}
		mockCompanyRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			if ticker == companyTicker {
				return sampleCompany, nil
			}
			return nil, errors.New("company not found")
		}
		mockPortfolioRepo.SaveFunc = func(p *portfolio.Portfolio) error {
			mockPortfolioRepo.SaveCalledWith = p // Capture for assertion
			return nil
		}

		err := service.AddPosition(portfolioID, companyTicker, shares, *purchasePrice)
		if err != nil {
			t.Fatalf("AddPosition() error = %v, wantErr nil", err)
		}

		if mockPortfolioRepo.SaveCalledWith == nil {
			t.Fatalf("Save was not called on portfolio repository")
		}
		savedPortfolio := mockPortfolioRepo.SaveCalledWith
		if len(savedPortfolio.Holdings) != 1 {
			t.Errorf("Expected 1 holding, got %d", len(savedPortfolio.Holdings))
		}
		pos, ok := savedPortfolio.Holdings[companyTicker]
		if !ok {
			t.Errorf("Holding for %s not found", companyTicker)
		} else {
			if pos.Shares != shares {
				t.Errorf("Shares for %s = %d, want %d", companyTicker, pos.Shares, shares)
			}
		}
		expectedCash := initialCash.Amount - (purchasePrice.Amount * int64(shares))
		if savedPortfolio.CashBalance.Amount != expectedCash {
			t.Errorf("CashBalance = %d, want %d", savedPortfolio.CashBalance.Amount, expectedCash)
		}
	})

	t.Run("PortfolioNotFound", func(t *testing.T) {
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			return nil, errors.New("portfolio not found error")
		}
		err := service.AddPosition(uuid.NewString(), companyTicker, shares, *purchasePrice)
		if err == nil {
			t.Errorf("AddPosition() with non-existent portfolio ID expected error, got nil")
		}
	})

	t.Run("CompanyNotFound", func(t *testing.T) {
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			return existingPortfolio, nil
		}
		mockCompanyRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			return nil, errors.New("company ticker not found in DB") // Simulate company not found
		}
		err := service.AddPosition(portfolioID, "UNKNOWNCO", shares, *purchasePrice)
		if err == nil {
			t.Errorf("AddPosition() with non-existent company ticker expected error, got nil")
		}
		// Expected: "failed to verify company ticker UNKNOWNCO: company ticker not found in DB" or "company with ticker UNKNOWNCO not found"
	})

	t.Run("InsufficientFunds", func(t *testing.T) {
		smallCash, _ := portfolio.NewMoney(100, "USD") // 1.00 USD
		poorPortfolio, _ := portfolio.NewPortfolio(portfolioID, portfolio.Aggressive, *smallCash)
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			return poorPortfolio, nil
		}
		mockCompanyRepo.FindByTickerFunc = func(ticker string) (*company.Company, error) {
			return sampleCompany, nil
		}
		// shares (10) * purchasePrice (150.00) = 1500.00 USD needed
		err := service.AddPosition(portfolioID, companyTicker, shares, *purchasePrice)
		if err == nil {
			t.Errorf("AddPosition() with insufficient funds expected domain error, got nil")
		}
		// Expected: "domain error adding position ...: insufficient cash balance to add position"
	})
	
	t.Run("EmptyPortfolioID", func(t *testing.T) {
		err := service.AddPosition("", companyTicker, shares, *purchasePrice)
		if err == nil { t.Error("Expected error for empty portfolio ID") }
	})
	t.Run("EmptyCompanyTicker", func(t *testing.T) {
		err := service.AddPosition(portfolioID, "", shares, *purchasePrice)
		if err == nil { t.Error("Expected error for empty company ticker") }
	})
	t.Run("NonPositiveShares", func(t *testing.T) {
		err := service.AddPosition(portfolioID, companyTicker, 0, *purchasePrice)
		if err == nil { t.Error("Expected error for zero shares") }
	})

}

// TestPortfolioService_AdjustPosition - Placeholder, as domain logic is very basic
func TestPortfolioService_AdjustPosition(t *testing.T) {
	mockPortfolioRepo := &MockPortfolioRepository{}
	mockCompanyRepo := &MinimalMockCompanyRepository{} // Not directly used by AdjustPosition's current simplified logic
	service := application.NewPortfolioService(mockPortfolioRepo, mockCompanyRepo)

	portfolioID := uuid.NewString()
	initialCash, _ := portfolio.NewMoney(100000, "USD")
	existingPortfolio, _ := portfolio.NewPortfolio(portfolioID, portfolio.Aggressive, *initialCash)
	
	// Pre-add a position
	ticker := "MSFT"
	oldShares := 10
	price, _ := portfolio.NewMoney(5000, "USD")
	pos, _ := portfolio.NewPosition(ticker, oldShares, *price)
	cost, _ := portfolio.NewMoney(price.Amount * int64(oldShares), price.Currency)
	_ = existingPortfolio.AddPosition(*pos, *cost) // Add directly for test setup convenience


	t.Run("Success_AdjustShares", func(t *testing.T) {
		// Important: GetPortfolioDetails returns a *copy* or the *actual object* based on FindByIDFunc.
		// For this test, we want the service to operate on the 'existingPortfolio' we've set up.
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			if id == portfolioID {
				// Return the portfolio instance that has the position already
				return existingPortfolio, nil
			}
			return nil, errors.New("not found")
		}
		mockPortfolioRepo.SaveFunc = func(p *portfolio.Portfolio) error {
			mockPortfolioRepo.SaveCalledWith = p
			return nil
		}
		
		newShares := 15
		err := service.AdjustPosition(portfolioID, ticker, newShares)
		if err != nil {
			t.Fatalf("AdjustPosition() error = %v, wantErr nil", err)
		}
		if mockPortfolioRepo.SaveCalledWith == nil {
			t.Fatal("Save was not called")
		}
		if savedPos, ok := mockPortfolioRepo.SaveCalledWith.Holdings[ticker]; !ok {
			t.Errorf("Position for %s not found after adjustment", ticker)
		} else if savedPos.Shares != newShares {
			t.Errorf("Shares for %s = %d, want %d", ticker, savedPos.Shares, newShares)
		}
		// Note: Cash balance adjustment is NOT part of this simplified AdjustPosition service method.
	})
	
	t.Run("PortfolioNotFound", func(t *testing.T) {
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) { return nil, errors.New("not found"); }
		err := service.AdjustPosition(uuid.NewString(), "ANY", 5)
		if err == nil { t.Error("Expected error for non-existent portfolio") }
	})

	t.Run("PositionNotFoundInPortfolio", func(t *testing.T) {
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) { 
			freshP, _ := portfolio.NewPortfolio(portfolioID, portfolio.Aggressive, *initialCash) // Portfolio without the position
			return freshP, nil
		}
		err := service.AdjustPosition(portfolioID, "NONEXISTENT", 5)
		if err == nil { t.Error("Expected error for non-existent position in portfolio") }
		// Expected: "position for ticker NONEXISTENT not found in portfolio..."
	})
}


func TestPortfolioService_RecommendRebalance(t *testing.T) {
	mockPortfolioRepo := &MockPortfolioRepository{}
	service := application.NewPortfolioService(mockPortfolioRepo, nil) // CompanyRepo not needed for this method

	portfolioID := uuid.NewString()
	pInstance, _ := portfolio.NewPortfolio(portfolioID, portfolio.Moderate, portfolio.Money{Amount:1000, Currency:"USD"})

	t.Run("Success_Triggered", func(t *testing.T) {
		pInstance.LastRebalanceTime = time.Time{} // Ensure rebalance is triggered in domain logic
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			return pInstance, nil
		}

		rec, err := service.RecommendRebalance(portfolioID)
		if err != nil {
			t.Fatalf("RecommendRebalance() error = %v, wantErr nil", err)
		}
		if rec == nil {
			t.Fatalf("RecommendRebalance() recommendation = nil, want non-nil")
		}
		if rec.PortfolioID != portfolioID {
			t.Errorf("Recommendation PortfolioID = %s, want %s", rec.PortfolioID, portfolioID)
		}
		if len(rec.Suggestions) == 0 { // Based on placeholder domain logic
			t.Error("Expected suggestions, got empty")
		}
	})

	t.Run("Success_NotTriggeredErrorFromDomain", func(t *testing.T) {
		pInstance.LastRebalanceTime = time.Now().Add(-10 * 24 * time.Hour) // Recently rebalanced
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			return pInstance, nil
		}
		_, err := service.RecommendRebalance(portfolioID)
		if err == nil {
			t.Errorf("RecommendRebalance() expected error when not triggered by domain, got nil")
		}
		// Expected: "domain error generating rebalance recommendations ...: rebalance not currently triggered"
	})
}

func TestPortfolioService_ExecuteRebalance(t *testing.T) {
	mockPortfolioRepo := &MockPortfolioRepository{}
	service := application.NewPortfolioService(mockPortfolioRepo, nil)

	portfolioID := uuid.NewString()
	pInstance, _ := portfolio.NewPortfolio(portfolioID, portfolio.Moderate, portfolio.Money{Amount:1000, Currency:"USD"})
	originalLastRebalanceTime := pInstance.LastRebalanceTime

	recommendation := application.RebalanceRecommendation{
		PortfolioID: portfolioID,
		Suggestions: []string{"Sell AAPL", "Buy MSFT"},
		GeneratedAt: time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		mockPortfolioRepo.SaveCalledWith = nil
		mockPortfolioRepo.FindByIDFunc = func(id string) (*portfolio.Portfolio, error) {
			// Return a fresh instance to ensure LastRebalanceTime is as expected pre-call
			freshP, _ := portfolio.NewPortfolio(portfolioID, portfolio.Moderate, portfolio.Money{Amount:1000, Currency:"USD"})
			freshP.LastRebalanceTime = originalLastRebalanceTime
			return freshP, nil
		}
		mockPortfolioRepo.SaveFunc = func(p *portfolio.Portfolio) error {
			mockPortfolioRepo.SaveCalledWith = p
			return nil
		}

		err := service.ExecuteRebalance(portfolioID, recommendation)
		if err != nil {
			t.Fatalf("ExecuteRebalance() error = %v, wantErr nil", err)
		}
		if mockPortfolioRepo.SaveCalledWith == nil {
			t.Fatal("Save was not called")
		}
		// Check if LastRebalanceTime was updated (placeholder logic in service does this)
		if mockPortfolioRepo.SaveCalledWith.LastRebalanceTime.Equal(originalLastRebalanceTime) {
			t.Errorf("LastRebalanceTime was not updated. Original: %v, Current: %v",
				originalLastRebalanceTime, mockPortfolioRepo.SaveCalledWith.LastRebalanceTime)
		}
	})

	t.Run("MismatchedPortfolioID", func(t *testing.T) {
		wrongRec := application.RebalanceRecommendation{PortfolioID: "wrong-id"}
		err := service.ExecuteRebalance(portfolioID, wrongRec)
		if err == nil {
			t.Error("Expected error for mismatched portfolio ID in recommendation")
		}
	})
}
