package portfolio

// Import company.Sector if direct type usage is intended and allowed.
// For now, we'll assume sector is a string that can be matched or
// that a more sophisticated cross-context communication mechanism would be used later.
// import "github.com/path-to-your-repo/pkg/domain/company"

// PortfolioRepository defines the interface for accessing and persisting Portfolio aggregates.
type PortfolioRepository interface {
	// FindByID retrieves a portfolio by its unique identifier.
	FindByID(id string) (*Portfolio, error)

	// FindAll retrieves all portfolios.
	// This might be resource-intensive and should be used judiciously or with pagination in a real system.
	FindAll() ([]*Portfolio, error)

	// SearchByRiskProfile retrieves portfolios matching a specific risk profile.
	SearchByRiskProfile(riskProfile RiskProfile) ([]*Portfolio, error)

	// Save creates a new portfolio or updates an existing one in the repository.
	// Implementations should handle the logic for differentiating between create and update.
	Save(portfolio *Portfolio) error

	// Delete removes a portfolio from the repository by its ID.
	// This method is optional for the initial MVP but good to define for completeness.
	Delete(id string) error

	// Note: SearchBySector from the prompt implies a dependency on the Company context.
	// For a pure DDD approach, this might involve:
	// 1. The Portfolio context storing only company tickers.
	// 2. An Application Service querying the Company context for tickers in a sector,
	//    then querying the Portfolio context for portfolios holding those tickers.
	// 3. Or, a denormalized read model updated by events.
	// For MVP, if a direct query is needed, the repository might take a simple string for sector
	// and the infrastructure layer would handle the join or multi-step query.
	// Example: SearchByCompanySector(sectorName string) ([]*Portfolio, error)
}
