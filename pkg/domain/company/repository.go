package company

// CompanyRepository defines the interface for accessing and persisting Company aggregates.
// Implementations will handle the underlying data storage (e.g., in-memory, database).
type CompanyRepository interface {
	// FindByTicker retrieves a company by its stock ticker.
	FindByTicker(ticker string) (*Company, error)

	// SearchByScoreRange retrieves companies whose current value score falls within the given range.
	SearchByScoreRange(minScore, maxScore float64) ([]*Company, error)

	// Save creates or updates a company in the repository.
	// If the company with the given ticker already exists, it should be updated.
	// Otherwise, a new company entry should be created.
	Save(company *Company) error

	// Delete removes a company from the repository by its ticker.
	// This method is optional for the initial MVP but good to define.
	Delete(ticker string) error

	// FindAll (Optional) retrieves all companies. Useful for some scenarios.
	// FindAll() ([]*Company, error)

	// FindBySector (Optional) retrieves companies belonging to a specific sector.
	// FindBySector(sector Sector) ([]*Company, error)
}
