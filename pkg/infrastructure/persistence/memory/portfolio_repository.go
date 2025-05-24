package memory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jizumer/expedition-value/pkg/domain/company"
	"github.com/jizumer/expedition-value/pkg/domain/portfolio"
)

// ErrPortfolioNotFound is returned when a portfolio is not found.
var ErrPortfolioNotFound = errors.New("portfolio not found")

// InMemoryPortfolioRepository is an in-memory implementation of the PortfolioRepository.
type InMemoryPortfolioRepository struct {
	mu           sync.RWMutex
	portfolios   map[string]*portfolio.Portfolio // Keyed by Portfolio ID
	companyRepo  company.CompanyRepository       // For sector lookups
}

// NewInMemoryPortfolioRepository creates a new instance of InMemoryPortfolioRepository.
// It requires a CompanyRepository to look up company sectors for SearchBySector.
func NewInMemoryPortfolioRepository(compRepo company.CompanyRepository) *InMemoryPortfolioRepository {
	return &InMemoryPortfolioRepository{
		portfolios:  make(map[string]*portfolio.Portfolio),
		companyRepo: compRepo,
	}
}

// Save creates or updates a portfolio in the in-memory store.
func (r *InMemoryPortfolioRepository) Save(p *portfolio.Portfolio) error {
	if p == nil {
		return errors.New("portfolio cannot be nil")
	}
	if p.ID == "" {
		return errors.New("portfolio ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.portfolios[p.ID] = p
	return nil
}

// FindByID retrieves a portfolio by its unique identifier.
func (r *InMemoryPortfolioRepository) FindByID(id string) (*portfolio.Portfolio, error) {
	if id == "" {
		return nil, errors.New("portfolio ID cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	portfolio, exists := r.portfolios[id]
	if !exists {
		return nil, ErrPortfolioNotFound
	}
	return portfolio, nil
}

// FindAll retrieves all portfolios.
func (r *InMemoryPortfolioRepository) FindAll() ([]*portfolio.Portfolio, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]*portfolio.Portfolio, 0, len(r.portfolios))
	for _, p := range r.portfolios {
		results = append(results, p)
	}
	return results, nil
}

// SearchByRiskProfile retrieves portfolios matching a specific risk profile.
// (This was defined in the domain interface, adding implementation here)
func (r *InMemoryPortfolioRepository) SearchByRiskProfile(riskProfile portfolio.RiskProfile) ([]*portfolio.Portfolio, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*portfolio.Portfolio
	for _, p := range r.portfolios {
		if p.RiskProfile == riskProfile {
			results = append(results, p)
		}
	}
	return results, nil
}


// SearchBySector retrieves portfolios that hold positions in companies of the given sector.
// This implementation requires looking up company details using the CompanyRepository.
func (r *InMemoryPortfolioRepository) SearchBySector(sector company.Sector) ([]*portfolio.Portfolio, error) {
	if r.companyRepo == nil {
		return nil, errors.New("company repository is not available for sector search")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*portfolio.Portfolio
	seenPortfolios := make(map[string]bool) // To avoid adding the same portfolio multiple times

	for _, p := range r.portfolios {
		if seenPortfolios[p.ID] {
			continue
		}
		for _, holding := range p.Holdings {
			comp, err := r.companyRepo.FindByTicker(holding.CompanyTicker)
			if err != nil {
				// Handle error: log it, or decide if this means the portfolio shouldn't match
				// For now, we'll skip this holding if the company can't be found
				fmt.Printf("Warning: Could not find company with ticker %s during sector search: %v\n", holding.CompanyTicker, err)
				continue
			}
			if comp.Sector == sector {
				results = append(results, p)
				seenPortfolios[p.ID] = true
				break // Found a matching company in this portfolio, move to the next portfolio
			}
		}
	}
	return results, nil
}

// Update is effectively the same as Save for an in-memory repository.
func (r *InMemoryPortfolioRepository) Update(p *portfolio.Portfolio) error {
	return r.Save(p)
}

// Delete removes a portfolio from the repository by its ID.
func (r *InMemoryPortfolioRepository) Delete(id string) error {
	if id == "" {
		return errors.New("portfolio ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.portfolios[id]; !exists {
		return fmt.Errorf("portfolio with ID '%s' not found for deletion: %w", id, ErrPortfolioNotFound)
	}
	delete(r.portfolios, id)
	return nil
}
