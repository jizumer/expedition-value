package company

import "time"

// FinancialMetrics holds key financial ratios and data for a company.
// This is a value object.
type FinancialMetrics struct {
	PERatio         float64   // Price-to-Earnings Ratio
	PBRatio         float64   // Price-to-Book Ratio
	DebtToEquity    float64   // Debt-to-Equity Ratio
	MetricsUpdatedAt time.Time // Timestamp of when these metrics were last updated
	// Add other relevant financial metrics as needed for value calculation
}

// NewFinancialMetrics creates and returns a new FinancialMetrics instance.
// Basic validation can be added here.
func NewFinancialMetrics(pe, pb, de float64) (*FinancialMetrics, error) {
	// Placeholder: Add validation if necessary (e.g., ratios cannot be negative)
	return &FinancialMetrics{
		PERatio:         pe,
		PBRatio:         pb,
		DebtToEquity:    de,
		MetricsUpdatedAt: time.Now(), // Set to current time on creation or update
	}, nil
}
