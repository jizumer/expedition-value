package application

import (
	"errors" // Using standard errors for now
	"time"   // For setting UpdatedAt if decided here

	"github.com/user/project/pkg/domain/company" // Assuming this module path
)

// CompanyService provides application-level functionalities for managing companies.
// It orchestrates domain logic and interacts with the company repository.
type CompanyService struct {
	companyRepo company.CompanyRepository
}

// NewCompanyService creates a new instance of CompanyService.
func NewCompanyService(repo company.CompanyRepository) *CompanyService {
	return &CompanyService{
		companyRepo: repo,
	}
}

// GetCompanyByTicker retrieves a company by its stock ticker.
func (s *CompanyService) GetCompanyByTicker(ticker string) (*company.Company, error) {
	if ticker == "" {
		return nil, errors.New("ticker cannot be empty")
	}
	return s.companyRepo.FindByTicker(ticker)
}

// SearchCompaniesByScore retrieves companies whose current value score falls within the given range.
func (s *CompanyService) SearchCompaniesByScore(minScore, maxScore float64) ([]*company.Company, error) {
	if minScore > maxScore {
		return nil, errors.New("minScore cannot be greater than maxScore")
	}
	return s.companyRepo.SearchByScoreRange(minScore, maxScore)
}

// CreateCompany creates a new Company instance, validates it, and saves it to the repository.
func (s *CompanyService) CreateCompany(ticker string, metrics company.FinancialMetrics, sector company.Sector) (*company.Company, error) {
	// Create new company instance using the domain constructor
	newCompany, err := company.NewCompany(ticker, metrics, sector)
	if err != nil {
		return nil, err // Error from domain validation (e.g., empty ticker)
	}

	// The domain's NewCompany already sets UpdatedAt and initial score.
	// We can call domain methods for further validation if needed here.
	// For example:
	// if !newCompany.ValidateScore() {
	// return nil, errors.New("initial score is invalid")
	// }

	// Save the new company to the repository
	err = s.companyRepo.Save(newCompany)
	if err != nil {
		return nil, err
	}
	return newCompany, nil
}

// UpdateCompanyMetrics updates the financial metrics for a given company and triggers score recalculation.
func (s *CompanyService) UpdateCompanyMetrics(ticker string, newMetrics company.FinancialMetrics) error {
	if ticker == "" {
		return errors.New("ticker cannot be empty")
	}

	// Fetch the existing company
	existingCompany, err := s.companyRepo.FindByTicker(ticker)
	if err != nil {
		return err // Company not found or other repository error
	}
	if existingCompany == nil {
		return errors.New("company not found") // Should be covered by repo error, but good practice
	}

	// Call domain method to update metrics and recalculate score
	err = existingCompany.UpdateFinancialMetrics(newMetrics)
	if err != nil {
		return err // Error from domain logic during update
	}

	// Save the updated company
	// The CompanyRepository's Save method should handle both create and update.
	return s.companyRepo.Save(existingCompany)
}

// RefreshCompany triggers a refresh of a company's data, potentially involving external sources.
// For the MVP, this is a placeholder that calls domain logic for refreshing stale metrics.
func (s *CompanyService) RefreshCompany(ticker string) error {
	if ticker == "" {
		return errors.New("ticker cannot be empty")
	}

	// Fetch the existing company
	c, err := s.companyRepo.FindByTicker(ticker)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("company not found")
	}

	// Call domain method to refresh stale metrics
	// This method might update the company's state (e.g., FinancialMetrics, UpdatedAt)
	err = c.RefreshStaleMetrics()
	if err != nil {
		return err // Error from domain logic (e.g., failed to refresh)
	}

	// If RefreshStaleMetrics modified the company, it should have updated its internal state.
	// Now, save the potentially updated company back to the repository.
	return s.companyRepo.Save(c)
}

// InitializeGoModule is a helper to create a go.mod file if it doesn't exist.
// This is not part of the CompanyService itself but a utility for the agent.
// It should be called separately if needed.
func InitializeGoModule(modulePath string) error {
	// This function would use os/exec to run 'go mod init <modulePath>'
	// For the purpose of this exercise, we'll assume it's handled or not strictly needed by the tool environment.
	// If there are compilation errors due to missing go.mod, this would be the place to call it from.
	// For now, this is a conceptual placeholder.
	return nil
}
