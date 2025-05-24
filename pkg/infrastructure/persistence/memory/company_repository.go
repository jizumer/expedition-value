package memory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jizumer/expedition-value/pkg/domain/company"
)

// ErrCompanyNotFound is returned when a company is not found in the repository.
var ErrCompanyNotFound = errors.New("company not found")

// InMemoryCompanyRepository is an in-memory implementation of the CompanyRepository interface.
// It uses a map to store companies and a RWMutex for concurrent access.
type InMemoryCompanyRepository struct {
	mu        sync.RWMutex
	companies map[string]*company.Company // Keyed by Ticker
}

// NewInMemoryCompanyRepository creates a new instance of InMemoryCompanyRepository.
func NewInMemoryCompanyRepository() *InMemoryCompanyRepository {
	return &InMemoryCompanyRepository{
		companies: make(map[string]*company.Company),
	}
}

// Save creates or updates a company in the in-memory store.
func (r *InMemoryCompanyRepository) Save(c *company.Company) error {
	if c == nil {
		return errors.New("company cannot be nil")
	}
	if c.Ticker == "" {
		return errors.New("company ticker cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.companies[c.Ticker] = c
	return nil
}

// FindByTicker retrieves a company by its stock ticker.
func (r *InMemoryCompanyRepository) FindByTicker(ticker string) (*company.Company, error) {
	if ticker == "" {
		return nil, errors.New("ticker cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	company, exists := r.companies[ticker]
	if !exists {
		return nil, ErrCompanyNotFound
	}
	return company, nil
}

// SearchByScoreRange retrieves companies whose current value score falls within the given range.
func (r *InMemoryCompanyRepository) SearchByScoreRange(minScore, maxScore float64) ([]*company.Company, error) {
	if minScore > maxScore {
		return nil, errors.New("minScore cannot be greater than maxScore")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*company.Company
	for _, c := range r.companies {
		if c.CurrentScore >= minScore && c.CurrentScore <= maxScore {
			results = append(results, c)
		}
	}
	return results, nil
}

// Update is effectively the same as Save for an in-memory repository,
// as Save will overwrite if the key exists.
// This method is here to satisfy the interface if it were to have distinct behavior.
func (r *InMemoryCompanyRepository) Update(c *company.Company) error {
	return r.Save(c)
}

// Delete removes a company from the repository by its ticker.
func (r *InMemoryCompanyRepository) Delete(ticker string) error {
	if ticker == "" {
		return errors.New("ticker cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.companies[ticker]; !exists {
		return fmt.Errorf("company with ticker '%s' not found for deletion: %w", ticker, ErrCompanyNotFound)
	}
	delete(r.companies, ticker)
	return nil
}

// FindAll (Optional method from interface)
// func (r *InMemoryCompanyRepository) FindAll() ([]*company.Company, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
//
// 	companies := make([]*company.Company, 0, len(r.companies))
// 	for _, c := range r.companies {
// 		companies = append(companies, c)
// 	}
// 	return companies, nil
// }

// FindBySector (Optional method from interface)
// For in-memory, this is straightforward if Sector is directly on Company.
// func (r *InMemoryCompanyRepository) FindBySector(sector company.Sector) ([]*company.Company, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
//
// 	var results []*company.Company
// 	for _, c := range r.companies {
// 		if c.Sector == sector {
// 			results = append(results, c)
// 		}
// 	}
// 	return results, nil
// }
