package company_test

import (
	"testing"
	"time"

	"github.com/jizumer/expedition-value/pkg/domain/company"
)

func TestNewCompany(t *testing.T) {
	t.Run("ValidCompanyCreation", func(t *testing.T) {
		ticker := "AAPL"
		metrics, _ := company.NewFinancialMetrics(15.0, 2.0, 0.5)
		sector := company.Technology

		c, err := company.NewCompany(ticker, *metrics, sector)

		if err != nil {
			t.Errorf("NewCompany() error = %v, wantErr nil", err)
		}
		if c == nil {
			t.Errorf("NewCompany() returned nil company, want non-nil")
			return // Guard against nil pointer dereference
		}
		if c.Ticker != ticker {
			t.Errorf("NewCompany() Ticker = %v, want %v", c.Ticker, ticker)
		}
		if c.FinancialMetrics.PERatio != 15.0 { // Basic check
			t.Errorf("NewCompany() FinancialMetrics.PERatio = %v, want %v", c.FinancialMetrics.PERatio, 15.0)
		}
		if c.Sector != sector {
			t.Errorf("NewCompany() Sector = %v, want %v", c.Sector, sector)
		}
		if c.CurrentScore != 0 {
			t.Errorf("NewCompany() CurrentScore = %v, want %v", c.CurrentScore, 0)
		}
		if c.UpdatedAt.IsZero() {
			t.Errorf("NewCompany() UpdatedAt was not set")
		}
	})

	t.Run("EmptyTickerValidation", func(t *testing.T) {
		metrics, _ := company.NewFinancialMetrics(15.0, 2.0, 0.5)
		_, err := company.NewCompany("", *metrics, company.Technology)
		if err == nil {
			t.Errorf("NewCompany() with empty ticker expected error, got nil")
		}
		// Check for specific error if your NewCompany returns a typed error
		// For now, company.Errors.New("ticker cannot be empty") is not exported or typed in a way we can directly check
		// So we check if an error is returned.
	})
}

func TestCompany_CheckMetricsAge(t *testing.T) {
	metrics, _ := company.NewFinancialMetrics(10, 1, 1)

	t.Run("MetricsAreRecent", func(t *testing.T) {
		metrics.MetricsUpdatedAt = time.Now().Add(-24 * time.Hour) // 1 day old
		c, _ := company.NewCompany("TEST", *metrics, company.Technology)
		if !c.CheckMetricsAge() {
			t.Errorf("CheckMetricsAge() returned false for recent metrics, want true")
		}
	})

	t.Run("MetricsAreStale", func(t *testing.T) {
		metrics.MetricsUpdatedAt = time.Now().Add(-10 * 24 * time.Hour) // 10 days old
		c, _ := company.NewCompany("TEST", *metrics, company.Technology)
		if c.CheckMetricsAge() {
			t.Errorf("CheckMetricsAge() returned true for stale metrics, want false")
		}
	})

	t.Run("MetricsUpdateDateIsZero", func(t *testing.T) {
		metrics.MetricsUpdatedAt = time.Time{} // Zero time
		c, _ := company.NewCompany("TEST", *metrics, company.Technology)
		if c.CheckMetricsAge() {
			t.Errorf("CheckMetricsAge() returned true for zero time metrics, want false")
		}
	})
}

func TestCompany_ValidateScore(t *testing.T) {
	metrics, _ := company.NewFinancialMetrics(10, 1, 1)
	c, _ := company.NewCompany("TEST", *metrics, company.Technology)

	t.Run("ValidScore", func(t *testing.T) {
		c.CurrentScore = 50
		if !c.ValidateScore() {
			t.Errorf("ValidateScore() returned false for valid score %v, want true", c.CurrentScore)
		}
	})

	t.Run("ScoreTooLow", func(t *testing.T) {
		c.CurrentScore = -10
		if c.ValidateScore() {
			t.Errorf("ValidateScore() returned true for low score %v, want false", c.CurrentScore)
		}
	})

	t.Run("ScoreTooHigh", func(t *testing.T) {
		c.CurrentScore = 110
		if c.ValidateScore() {
			t.Errorf("ValidateScore() returned true for high score %v, want false", c.CurrentScore)
		}
	})
}

func TestCompany_RecalculateScoreOnMetricUpdate(t *testing.T) {
	// This test is illustrative, as the actual recalculation logic is a placeholder.
	// It will primarily test that the UpdatedAt field is modified.
	metrics, _ := company.NewFinancialMetrics(10, 1, 1)
	c, _ := company.NewCompany("TEST", *metrics, company.Technology)
	initialUpdateTime := c.UpdatedAt

	time.Sleep(1 * time.Millisecond) // Ensure time progresses for UpdatedAt check

	err := c.RecalculateScoreOnMetricUpdate()
	if err != nil {
		t.Fatalf("RecalculateScoreOnMetricUpdate() returned error: %v", err)
	}

	if c.UpdatedAt.Equal(initialUpdateTime) || c.UpdatedAt.Before(initialUpdateTime) {
		t.Errorf("UpdatedAt not advanced after RecalculateScoreOnMetricUpdate. Initial: %v, Current: %v", initialUpdateTime, c.UpdatedAt)
	}
	// When actual score calculation logic is added, assert c.CurrentScore changes as expected.
	// For example: if c.CurrentScore was 0, and metrics imply a score of X, then c.CurrentScore should be X.
}

func TestCompany_RefreshStaleMetrics(t *testing.T) {
	// This test is also illustrative for the placeholder logic.
	// It checks if FinancialMetrics.MetricsUpdatedAt and UpdatedAt are updated if metrics were stale.
	staleMetrics, _ := company.NewFinancialMetrics(10, 1, 1)
	staleMetrics.MetricsUpdatedAt = time.Now().Add(-10 * 24 * time.Hour) // 10 days old
	
	cStale, _ := company.NewCompany("STALE", *staleMetrics, company.Technology)
	initialCompanyUpdateTimeStale := cStale.UpdatedAt
	initialMetricsUpdateTimeStale := cStale.FinancialMetrics.MetricsUpdatedAt

	time.Sleep(1 * time.Millisecond) // Ensure time progresses

	err := cStale.RefreshStaleMetrics()
	if err != nil {
		t.Fatalf("RefreshStaleMetrics() for stale metrics returned error: %v", err)
	}

	if cStale.FinancialMetrics.MetricsUpdatedAt.Equal(initialMetricsUpdateTimeStale) || cStale.FinancialMetrics.MetricsUpdatedAt.Before(initialMetricsUpdateTimeStale) {
		t.Errorf("FinancialMetrics.MetricsUpdatedAt not advanced for stale metrics. Initial: %v, Current: %v", initialMetricsUpdateTimeStale, cStale.FinancialMetrics.MetricsUpdatedAt)
	}
	if cStale.UpdatedAt.Equal(initialCompanyUpdateTimeStale) || cStale.UpdatedAt.Before(initialCompanyUpdateTimeStale) {
		t.Errorf("Company.UpdatedAt not advanced for stale metrics. Initial: %v, Current: %v", initialCompanyUpdateTimeStale, cStale.UpdatedAt)
	}

	// Test with non-stale metrics
	recentMetrics, _ := company.NewFinancialMetrics(12, 1.2, 0.6)
	recentMetrics.MetricsUpdatedAt = time.Now().Add(-1 * 24 * time.Hour) // 1 day old
	
	cRecent, _ := company.NewCompany("RECENT", *recentMetrics, company.Technology)
	initialCompanyUpdateTimeRecent := cRecent.UpdatedAt
	initialMetricsUpdateTimeRecent := cRecent.FinancialMetrics.MetricsUpdatedAt
	
	time.Sleep(1 * time.Millisecond)

	err = cRecent.RefreshStaleMetrics()
	if err != nil {
		t.Fatalf("RefreshStaleMetrics() for recent metrics returned error: %v", err)
	}
	// For recent metrics, the timestamps should NOT change
	if !cRecent.FinancialMetrics.MetricsUpdatedAt.Equal(initialMetricsUpdateTimeRecent) {
		t.Errorf("FinancialMetrics.MetricsUpdatedAt changed for recent metrics. Initial: %v, Current: %v", initialMetricsUpdateTimeRecent, cRecent.FinancialMetrics.MetricsUpdatedAt)
	}
	// Company.UpdatedAt might still be updated if the method always touches it,
	// or it might not if there's no actual refresh.
	// The current placeholder logic in company.go does not update UpdatedAt if metrics are not stale.
	// If the intent is to always update UpdatedAt on call, this test part would need adjustment.
	// For now, assuming UpdatedAt is only touched if FinancialMetrics are actually refreshed.
	if !cRecent.UpdatedAt.Equal(initialCompanyUpdateTimeRecent) {
		t.Errorf("Company.UpdatedAt changed for recent metrics when no refresh occurred. Initial: %v, Current: %v", initialCompanyUpdateTimeRecent, cRecent.UpdatedAt)
	}
}

func TestCompany_UpdateFinancialMetrics(t *testing.T) {
	initialMetrics, _ := company.NewFinancialMetrics(10, 1, 0.5)
	c, _ := company.NewCompany("TEST", *initialMetrics, company.Technology)
	oldCompanyUpdateTs := c.UpdatedAt
	
	time.Sleep(1 * time.Millisecond) // Ensure time can advance

	newMetrics, _ := company.NewFinancialMetrics(20, 2, 0.6)
	// Explicitly set a different MetricsUpdatedAt for the new set of metrics,
	// although UpdateFinancialMetrics should set it to time.Now()
	newMetrics.MetricsUpdatedAt = time.Now().Add(-1 * time.Hour) 

	err := c.UpdateFinancialMetrics(*newMetrics)
	if err != nil {
		t.Fatalf("UpdateFinancialMetrics() returned error: %v", err)
	}

	if c.FinancialMetrics.PERatio != newMetrics.PERatio {
		t.Errorf("PERatio not updated. Got %v, want %v", c.FinancialMetrics.PERatio, newMetrics.PERatio)
	}
	if c.FinancialMetrics.PBRatio != newMetrics.PBRatio {
		t.Errorf("PBRatio not updated. Got %v, want %v", c.FinancialMetrics.PBRatio, newMetrics.PBRatio)
	}
	if c.FinancialMetrics.DebtToEquity != newMetrics.DebtToEquity {
		t.Errorf("DebtToEquity not updated. Got %v, want %v", c.FinancialMetrics.DebtToEquity, newMetrics.DebtToEquity)
	}

	if c.FinancialMetrics.MetricsUpdatedAt.Equal(newMetrics.MetricsUpdatedAt) {
		t.Errorf("FinancialMetrics.MetricsUpdatedAt was not set to current time by UpdateFinancialMetrics. Got %v", c.FinancialMetrics.MetricsUpdatedAt)
	}
	if c.FinancialMetrics.MetricsUpdatedAt.Before(oldCompanyUpdateTs) {
		t.Errorf("FinancialMetrics.MetricsUpdatedAt is older than the previous company update time. Got %v", c.FinancialMetrics.MetricsUpdatedAt)
	}
	if c.UpdatedAt.Equal(oldCompanyUpdateTs) || c.UpdatedAt.Before(oldCompanyUpdateTs) {
		t.Errorf("Company.UpdatedAt was not advanced. Initial: %v, Current: %v", oldCompanyUpdateTs, c.UpdatedAt)
	}
	// Further tests could assert that RecalculateScoreOnMetricUpdate was effectively called
	// (e.g., by checking score if logic existed, or by using a spy/mock if the method was an interface).
}

// Test for Domain Event Constructors - simple value checks
func TestScoreRecalculatedEvent(t *testing.T) {
	ticker := "EVTEST"
	oldScore := 50.0
	newScore := 65.5
	event := company.NewScoreRecalculatedEvent(ticker, oldScore, newScore)

	if event.Ticker != ticker {
		t.Errorf("NewScoreRecalculatedEvent Ticker = %s, want %s", event.Ticker, ticker)
	}
	if event.OldScore != oldScore {
		t.Errorf("NewScoreRecalculatedEvent OldScore = %f, want %f", event.OldScore, oldScore)
	}
	if event.NewScore != newScore {
		t.Errorf("NewScoreRecalculatedEvent NewScore = %f, want %f", event.NewScore, newScore)
	}
	if event.Timestamp.IsZero() {
		t.Error("NewScoreRecalculatedEvent Timestamp was not set")
	}
}

func TestMetricsUpdatedEvent(t *testing.T) {
	ticker := "EVTEST2"
	event := company.NewMetricsUpdatedEvent(ticker)

	if event.Ticker != ticker {
		t.Errorf("NewMetricsUpdatedEvent Ticker = %s, want %s", event.Ticker, ticker)
	}
	if event.Timestamp.IsZero() {
		t.Error("NewMetricsUpdatedEvent Timestamp was not set")
	}
}

// Test for FinancialMetrics constructor
func TestNewFinancialMetrics(t *testing.T) {
    pe := 10.0
    pb := 1.0
    de := 0.5

    fm, err := company.NewFinancialMetrics(pe, pb, de)
    if err != nil {
        t.Fatalf("NewFinancialMetrics returned an unexpected error: %v", err)
    }
    if fm == nil {
        t.Fatal("NewFinancialMetrics returned nil but no error")
    }
    if fm.PERatio != pe {
        t.Errorf("Expected PERatio %v, got %v", pe, fm.PERatio)
    }
    if fm.PBRatio != pb {
        t.Errorf("Expected PBRatio %v, got %v", pb, fm.PBRatio)
    }
    if fm.DebtToEquity != de {
        t.Errorf("Expected DebtToEquity %v, got %v", de, fm.DebtToEquity)
    }
    if fm.MetricsUpdatedAt.IsZero() {
        t.Error("Expected MetricsUpdatedAt to be set, but it was zero")
    }

    // Example of adding validation tests if NewFinancialMetrics had them
    // For instance, if PE ratio couldn't be negative:
    // _, err = company.NewFinancialMetrics(-1.0, pb, de)
    // if err == nil {
    //    t.Error("Expected error for negative PE ratio, got nil")
    // }
}

// Test for Sector enum String() and ParseSector()
func TestSectorEnum(t *testing.T) {
	testCases := []struct {
		sectorVal company.Sector
		strVal    string
	}{
		{company.Technology, "Technology"},
		{company.Healthcare, "Healthcare"},
		{company.Financials, "Financials"},
		{company.ConsumerDiscretionary, "Consumer Discretionary"},
		{company.ConsumerStaples, "Consumer Staples"},
		{company.Industrials, "Industrials"},
		{company.Energy, "Energy"},
		{company.Utilities, "Utilities"},
		{company.RealEstate, "Real Estate"},
		{company.Materials, "Materials"},
		{company.TelecommunicationServices, "Telecommunication Services"},
		{company.UndefinedSector, "UndefinedSector"},
		{company.Sector(99), "UndefinedSector"}, // Test out-of-range value
	}

	for _, tc := range testCases {
		t.Run("SectorToString_"+tc.strVal, func(t *testing.T) {
			if str := tc.sectorVal.String(); str != tc.strVal {
				t.Errorf("Sector(%d).String() = %q, want %q", tc.sectorVal, str, tc.strVal)
			}
		})
		// Only parse valid strings, skip "UndefinedSector" for direct parsing if it's a default
		if tc.sectorVal != company.UndefinedSector && tc.sectorVal <= company.TelecommunicationServices {
			t.Run("ParseSector_"+tc.strVal, func(t *testing.T) {
				if sector := company.ParseSector(tc.strVal); sector != tc.sectorVal {
					t.Errorf("ParseSector(%q) = %v, want %v", tc.strVal, sector, tc.sectorVal)
				}
			})
		}
	}

	t.Run("ParseSector_UnknownString", func(t *testing.T) {
		unknownStr := "NonExistentSector"
		if sector := company.ParseSector(unknownStr); sector != company.UndefinedSector {
			t.Errorf("ParseSector(%q) = %v, want %v (UndefinedSector)", unknownStr, sector, company.UndefinedSector)
		}
	})
}
