package application

import (
	"errors" // Using standard errors for now
	"fmt"    // For error formatting
	"time"   // For setting UpdatedAt if decided here

	// Project packages
	"github.com/jizumer/expedition-value/pkg/domain/company"
	"github.com/jizumer/expedition-value/pkg/domain/portfolio"

	"github.com/google/uuid" // For generating portfolio IDs
)

// RebalanceRecommendation is a DTO for rebalancing suggestions.
// For MVP, it's a simple structure.
type RebalanceRecommendation struct {
	PortfolioID string
	Suggestions []string  // Example: ["Sell 10 shares of AAPL", "Buy 5 shares of MSFT"]
	GeneratedAt time.Time // Timestamp when the recommendation was generated
}

// PortfolioService provides application-level functionalities for managing portfolios.
// It orchestrates domain logic and interacts with portfolio and company repositories.
type PortfolioService struct {
	portfolioRepo portfolio.PortfolioRepository
	companyRepo   company.CompanyRepository // To validate company tickers
}

// NewPortfolioService creates a new instance of PortfolioService.
func NewPortfolioService(pRepo portfolio.PortfolioRepository, cRepo company.CompanyRepository) *PortfolioService {
	return &PortfolioService{
		portfolioRepo: pRepo,
		companyRepo:   cRepo,
	}
}

// CreatePortfolio creates a new Portfolio instance, generates an ID, and saves it.
func (s *PortfolioService) CreatePortfolio(cashBalance portfolio.Money, riskProfile portfolio.RiskProfile) (*portfolio.Portfolio, error) {
	portfolioID := uuid.NewString() // Generate a unique ID for the new portfolio

	newPortfolio, err := portfolio.NewPortfolio(portfolioID, riskProfile, cashBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to create new portfolio in domain: %w", err)
	}

	// Save the new portfolio to the repository
	err = s.portfolioRepo.Save(newPortfolio)
	if err != nil {
		return nil, fmt.Errorf("failed to save portfolio: %w", err)
	}
	return newPortfolio, nil
}

// GetPortfolioDetails retrieves a portfolio by its ID.
func (s *PortfolioService) GetPortfolioDetails(portfolioID string) (*portfolio.Portfolio, error) {
	if portfolioID == "" {
		return nil, errors.New("portfolioID cannot be empty")
	}
	p, err := s.portfolioRepo.FindByID(portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to find portfolio %s: %w", portfolioID, err)
	}
	if p == nil {
		return nil, fmt.Errorf("portfolio %s not found", portfolioID) // More specific error
	}
	return p, nil
}

// AddPosition adds a new position to an existing portfolio.
func (s *PortfolioService) AddPosition(portfolioID string, companyTicker string, shares int, purchasePrice portfolio.Money) error {
	if portfolioID == "" {
		return errors.New("portfolioID cannot be empty")
	}
	if companyTicker == "" {
		return errors.New("companyTicker cannot be empty")
	}
	if shares <= 0 {
		return errors.New("shares must be positive")
	}

	// Optional: Validate company ticker
	if s.companyRepo != nil {
		comp, err := s.companyRepo.FindByTicker(companyTicker)
		if err != nil {
			return fmt.Errorf("failed to verify company ticker %s: %w", companyTicker, err)
		}
		if comp == nil {
			return fmt.Errorf("company with ticker %s not found", companyTicker)
		}
	}

	// Fetch the portfolio
	p, err := s.GetPortfolioDetails(portfolioID) // Use existing method to get portfolio
	if err != nil {
		return err
	}

	// Create the position (using domain constructor if available, or directly)
	newPosition, err := portfolio.NewPosition(companyTicker, shares, purchasePrice)
	if err != nil {
		return fmt.Errorf("failed to create new position: %w", err)
	}

	// Calculate cost (simplified: purchasePrice is per share)
	// Proper money multiplication would be in Money VO
	cost := portfolio.Money{Amount: purchasePrice.Amount * int64(shares), Currency: purchasePrice.Currency}

	// Call domain method to add position
	err = p.AddPosition(*newPosition, cost) // Assuming AddPosition is a method on *Portfolio
	if err != nil {
		return fmt.Errorf("domain error adding position to portfolio %s: %w", portfolioID, err)
	}

	// Save the updated portfolio
	err = s.portfolioRepo.Save(p)
	if err != nil {
		return fmt.Errorf("failed to save updated portfolio %s: %w", portfolioID, err)
	}
	return nil
}

// AdjustPosition modifies an existing position in a portfolio.
// For simplicity, this example assumes adjusting means changing the number of shares.
// A more robust implementation might handle price changes, splits, etc.
func (s *PortfolioService) AdjustPosition(portfolioID string, companyTicker string, newShares int /*, newAveragePrice *portfolio.Money */) error {
	if portfolioID == "" {
		return errors.New("portfolioID cannot be empty")
	}
	if companyTicker == "" {
		return errors.New("companyTicker cannot be empty")
	}
	if newShares <= 0 { // Assuming adjusting to 0 means closing the position
		return errors.New("new shares count must be positive; use RemovePosition to close")
	}

	p, err := s.GetPortfolioDetails(portfolioID)
	if err != nil {
		return err
	}

	// --- Domain logic to adjust position ---
	// This is a simplified placeholder. The actual logic would be more complex:
	// - Find the existing position.
	// - Calculate difference in shares.
	// - If shares increase: calculate cost, check cash balance, update cash balance.
	// - If shares decrease: calculate proceeds, update cash balance.
	// - Update the position's share count and potentially average price.
	// - The Portfolio aggregate should enforce these rules.

	existingPosition, ok := p.Holdings[companyTicker]
	if !ok {
		return fmt.Errorf("position for ticker %s not found in portfolio %s", companyTicker, portfolioID)
	}

	// Simplified: just update shares. Real logic needs cost/proceeds & cash adjustment.
	// This should ideally call a method on `p` like `p.AdjustHolding(companyTicker, newShares)`
	// For now, directly modifying for brevity, but this bypasses domain logic.
	// This is a placeholder for where a proper domain method call would go.
	p.Holdings[companyTicker] = portfolio.Position{
		CompanyTicker: existingPosition.CompanyTicker,
		Shares:        newShares,
		PurchasePrice: existingPosition.PurchasePrice, // Average price would change in reality
	}
	p.UpdatedAt = time.Now()
	// --- End of simplified domain logic placeholder ---

	// Save the updated portfolio
	err = s.portfolioRepo.Save(p)
	if err != nil {
		return fmt.Errorf("failed to save updated portfolio %s after adjusting position: %w", portfolioID, err)
	}
	return nil
}

// RecommendRebalance generates rebalancing recommendations for a portfolio.
func (s *PortfolioService) RecommendRebalance(portfolioID string) (*RebalanceRecommendation, error) {
	if portfolioID == "" {
		return nil, errors.New("portfolioID cannot be empty")
	}

	p, err := s.GetPortfolioDetails(portfolioID)
	if err != nil {
		return nil, err
	}

	// Call domain logic on the portfolio to generate recommendations.
	// The domain method `GenerateRebalanceRecommendations` is a placeholder.
	// It would contain the actual logic for determining what to buy/sell.
	domainRecs, err := p.GenerateRebalanceRecommendations()
	if err != nil {
		// This could be an error like "rebalance not needed" or a real calculation error.
		return nil, fmt.Errorf("domain error generating rebalance recommendations for portfolio %s: %w", portfolioID, err)
	}

	recommendation := &RebalanceRecommendation{
		PortfolioID: portfolioID,
		Suggestions: domainRecs, // Assuming domainRecs is []string as per domain placeholder
		GeneratedAt: time.Now(),
	}

	return recommendation, nil
}

// ExecuteRebalance applies a given rebalancing recommendation to the portfolio.
func (s *PortfolioService) ExecuteRebalance(portfolioID string, recommendation RebalanceRecommendation) error {
	if portfolioID == "" {
		return errors.New("portfolioID cannot be empty")
	}
	if recommendation.PortfolioID != portfolioID {
		return errors.New("recommendation portfolioID does not match provided portfolioID")
	}

	p, err := s.GetPortfolioDetails(portfolioID)
	if err != nil {
		return err
	}

	// --- Domain logic to apply rebalance ---
	// This would involve:
	// - Iterating through recommendation.Suggestions.
	// - For each suggestion (e.g., "Sell 10 AAPL", "Buy 5 MSFT"):
	//   - Parse the action, ticker, shares.
	//   - Call domain methods like `p.RemovePosition` or `p.AddPosition`.
	//   - These domain methods must handle cash adjustments.
	// - This entire process should be transactional within the Portfolio aggregate.
	// For now, this is a placeholder as the domain `ApplyRebalance` is not fully defined.
	// p.ApplyRebalance(recommendation) // This would be the ideal call
	fmt.Printf("Executing rebalance for portfolio %s with %d suggestions (placeholder)\n", portfolioID, len(recommendation.Suggestions))
	p.LastRebalanceTime = time.Now() // Mark as rebalanced
	p.UpdatedAt = time.Now()
	// --- End of placeholder ---

	err = s.portfolioRepo.Save(p)
	if err != nil {
		return fmt.Errorf("failed to save portfolio %s after executing rebalance: %w", portfolioID, err)
	}
	return nil
}
