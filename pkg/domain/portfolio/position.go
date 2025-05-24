package portfolio

// Position represents a holding of a specific company's stock within a portfolio.
// This is a value object when considered within the context of a single portfolio,
// but might be an entity if it had its own lifecycle and identity across portfolios (not the case here).
type Position struct {
	CompanyTicker string // Stock ticker of the company
	Shares        int    // Number of shares held
	PurchasePrice Money  // Average purchase price per share for this position
	// CurrentMarketValue could be added if needed, but might be calculated dynamically.
}

// NewPosition creates a new Position instance.
// Basic validation can be added here.
func NewPosition(ticker string, shares int, purchasePrice Money) (*Position, error) {
	if ticker == "" {
		return nil, Errors.New("company ticker cannot be empty")
	}
	if shares <= 0 {
		return nil, Errors.New("shares must be positive")
	}
	// Add more validation for purchasePrice if necessary (e.g., positive amount)
	return &Position{
		CompanyTicker: ticker,
		Shares:        shares,
		PurchasePrice: purchasePrice,
	}, nil
}

// errors is a placeholder for a proper error handling package or built-in errors.
// For now, we'll use a simple error type.
// Custom error handling (if any specific to Position logic) should ideally use
// the 'Errors' instance from the portfolio.go file within this package,
// or the standard 'errors' package for generic errors.
// The duplicated custom errors struct has been removed from here.
