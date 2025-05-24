package company

import (
	"time"
)

// Company represents a publicly traded company and its value investment analysis data.
// It is an aggregate root.
type Company struct {
	Ticker           string
	FinancialMetrics FinancialMetrics // Defined in financial_metrics.go
	CurrentScore     float64
	Sector           Sector // Enum defined in sector.go
	UpdatedAt        time.Time
}

// NewCompany creates a new Company instance.
// Additional validation logic can be added here.
func NewCompany(ticker string, metrics FinancialMetrics, sector Sector) (*Company, error) {
	// Basic validation, more can be added.
	if ticker == "" {
		return nil, Errors.New("ticker cannot be empty")
	}
	return &Company{
		Ticker:           ticker,
		FinancialMetrics: metrics,
		Sector:           sector,
		CurrentScore:     0, // Initial score, will be calculated
		UpdatedAt:        time.Now(),
	}, nil
}

// --- Invariant Enforcement Methods (Placeholders) ---

// CheckMetricsAge verifies if the financial metrics are up-to-date.
// This is an example of an invariant.
func (c *Company) CheckMetricsAge() bool {
	// Placeholder: Implement logic to check if FinancialMetrics.UpdatedAt is recent enough.
	// For example, metrics older than 7 days might be considered stale.
	if c.FinancialMetrics.MetricsUpdatedAt.IsZero() { // Assuming MetricsUpdatedAt is a field in FinancialMetrics
		return false // No data
	}
	return time.Since(c.FinancialMetrics.MetricsUpdatedAt) < (7 * 24 * time.Hour)
}

// ValidateScore ensures the CurrentScore is within a logical range (e.g., 0-100).
// This is another example of an invariant.
func (c *Company) ValidateScore() bool {
	// Placeholder: Implement logic to check if CurrentScore is valid.
	return c.CurrentScore >= 0 && c.CurrentScore <= 100
}

// --- Corrective Policy Methods (Placeholders) ---

// RefreshStaleMetrics initiates a process to update financial metrics if they are stale.
// This is an example of a corrective policy.
func (c *Company) RefreshStaleMetrics() error {
	// Placeholder: Implement logic to trigger a refresh of financial metrics.
	// This might involve fetching new data from an external source.
	// After updating, c.FinancialMetrics.MetricsUpdatedAt and c.UpdatedAt should be updated.
	// A Domain Event (e.g., MetricsRefreshedEvent) could be published.
	if !c.CheckMetricsAge() {
		// Simulate metrics refresh
		// c.FinancialMetrics = getNewMetrics()
		c.FinancialMetrics.MetricsUpdatedAt = time.Now() // Update timestamp after refresh
		c.UpdatedAt = time.Now()
		// Publish MetricsRefreshedEvent (details to be implemented)
	}
	return nil
}

// RecalculateScoreOnMetricUpdate recalculates the CurrentScore when financial metrics change.
// This is another corrective policy, often triggered after metrics are updated.
func (c *Company) RecalculateScoreOnMetricUpdate() error {
	// Placeholder: Implement logic to recalculate CurrentScore based on FinancialMetrics.
	// oldScore := c.CurrentScore
	// c.CurrentScore = calculateNewScore(c.FinancialMetrics)
	c.UpdatedAt = time.Now()
	// if oldScore != c.CurrentScore {
	// Publish ScoreRecalculatedEvent
	// }
	return nil
}

// UpdateFinancialMetrics updates the company's financial metrics and triggers a score recalculation.
func (c *Company) UpdateFinancialMetrics(newMetrics FinancialMetrics) error {
	c.FinancialMetrics = newMetrics
	c.FinancialMetrics.MetricsUpdatedAt = time.Now() // Ensure this is set
	c.UpdatedAt = time.Now()
	return c.RecalculateScoreOnMetricUpdate()
}

// --- Domain Event Types (Placeholders) ---

// ScoreRecalculatedEvent indicates that a company's score has been recalculated.
type ScoreRecalculatedEvent struct {
	Ticker    string
	OldScore  float64
	NewScore  float64
	Timestamp time.Time
}

// NewScoreRecalculatedEvent creates a new ScoreRecalculatedEvent.
func NewScoreRecalculatedEvent(ticker string, oldScore, newScore float64) ScoreRecalculatedEvent {
	return ScoreRecalculatedEvent{
		Ticker:    ticker,
		OldScore:  oldScore,
		NewScore:  newScore,
		Timestamp: time.Now(),
	}
}

// MetricsUpdatedEvent indicates that a company's financial metrics have been updated.
type MetricsUpdatedEvent struct {
	Ticker    string
	Timestamp time.Time
}

// NewMetricsUpdatedEvent creates a new MetricsUpdatedEvent.
func NewMetricsUpdatedEvent(ticker string) MetricsUpdatedEvent {
	return MetricsUpdatedEvent{
		Ticker:    ticker,
		Timestamp: time.Now(),
	}
}

// errors is a placeholder for a proper error handling package or built-in errors.
// For now, we'll use a simple error type.
type errors struct{}

func (e *errors) New(text string) error {
	return &customError{text}
}

type customError struct {
	s string
}

func (e *customError) Error() string {
	return e.s
}

var Errors = &errors{}
