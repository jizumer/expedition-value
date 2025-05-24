package portfolio_test

import (
	"testing"
	"time"

	"github.com/jizumer/expedition-value/pkg/domain/portfolio"
	// "github.com/google/uuid" // If testing ID generation specifically
)

func TestNewPortfolio(t *testing.T) {
	initialCash, _ := portfolio.NewMoney(100000, "USD") // 1000.00 USD
	riskProfile := portfolio.Moderate

	t.Run("ValidPortfolioCreation", func(t *testing.T) {
		id := "test-portfolio-123"
		p, err := portfolio.NewPortfolio(id, riskProfile, *initialCash)

		if err != nil {
			t.Fatalf("NewPortfolio() error = %v, wantErr nil", err)
		}
		if p == nil {
			t.Fatalf("NewPortfolio() returned nil portfolio, want non-nil")
		}
		if p.ID != id {
			t.Errorf("NewPortfolio() ID = %s, want %s", p.ID, id)
		}
		if len(p.Holdings) != 0 {
			t.Errorf("NewPortfolio() Holdings count = %d, want 0", len(p.Holdings))
		}
		if p.CashBalance.Amount != initialCash.Amount || p.CashBalance.Currency != initialCash.Currency {
			t.Errorf("NewPortfolio() CashBalance = %v, want %v", p.CashBalance, *initialCash)
		}
		if p.RiskProfile != riskProfile {
			t.Errorf("NewPortfolio() RiskProfile = %v, want %v", p.RiskProfile, riskProfile)
		}
		if !p.LastRebalanceTime.IsZero() {
			t.Errorf("NewPortfolio() LastRebalanceTime should be zero, got %v", p.LastRebalanceTime)
		}
		if p.UpdatedAt.IsZero() {
			t.Errorf("NewPortfolio() UpdatedAt was not set")
		}
	})

	t.Run("EmptyIDValidation", func(t *testing.T) {
		_, err := portfolio.NewPortfolio("", riskProfile, *initialCash)
		if err == nil {
			t.Errorf("NewPortfolio() with empty ID expected error, got nil")
		}
	})

	t.Run("NegativeInitialCashValidation", func(t *testing.T) {
		negativeCash, _ := portfolio.NewMoney(-100, "USD")
		_, err := portfolio.NewPortfolio("test-id", riskProfile, *negativeCash)
		if err == nil {
			t.Errorf("NewPortfolio() with negative initial cash expected error, got nil")
		}
	})
}

func TestPortfolio_ValidateCashBalance(t *testing.T) {
	p, _ := portfolio.NewPortfolio("test", portfolio.Conservative, portfolio.Money{Amount: 100, Currency: "USD"})

	t.Run("PositiveCashBalance", func(t *testing.T) {
		p.CashBalance.Amount = 5000
		if !p.ValidateCashBalance() {
			t.Errorf("ValidateCashBalance() returned false for positive balance, want true")
		}
	})

	t.Run("ZeroCashBalance", func(t *testing.T) {
		p.CashBalance.Amount = 0
		if !p.ValidateCashBalance() {
			t.Errorf("ValidateCashBalance() returned false for zero balance, want true")
		}
	})

	t.Run("NegativeCashBalance", func(t *testing.T) {
		p.CashBalance.Amount = -100
		if p.ValidateCashBalance() {
			t.Errorf("ValidateCashBalance() returned true for negative balance, want false")
		}
	})
}

func TestPortfolio_CheckRebalanceTrigger(t *testing.T) {
	p, _ := portfolio.NewPortfolio("test", portfolio.Moderate, portfolio.Money{Amount: 10000, Currency: "USD"})

	t.Run("NeverRebalanced", func(t *testing.T) {
		p.LastRebalanceTime = time.Time{} 
		if !p.CheckRebalanceTrigger() {
			t.Errorf("CheckRebalanceTrigger() returned false for never rebalanced portfolio, want true")
		}
	})

	t.Run("RecentlyRebalanced", func(t *testing.T) {
		p.LastRebalanceTime = time.Now().Add(-30 * 24 * time.Hour) 
		if p.CheckRebalanceTrigger() { 
			t.Errorf("CheckRebalanceTrigger() returned true for recently rebalanced portfolio, want false")
		}
	})

	t.Run("NeedsRebalanceDueToTime", func(t *testing.T) {
		p.LastRebalanceTime = time.Now().Add(-4 * 30 * 24 * time.Hour) 
		if !p.CheckRebalanceTrigger() {
			t.Errorf("CheckRebalanceTrigger() returned false for portfolio needing rebalance due to time, want true")
		}
	})
}

func TestPortfolio_AddPosition(t *testing.T) {
	initialCash, _ := portfolio.NewMoney(100000, "USD") // 1000.00 USD

	t.Run("SuccessfulAdd", func(t *testing.T) {
		pFresh, _ := portfolio.NewPortfolio("pFresh", portfolio.Aggressive, *initialCash)
		originalUpdatedAt := pFresh.UpdatedAt
		time.Sleep(1 * time.Millisecond) 
		
		pos1Price, _ := portfolio.NewMoney(10000, "USD") 
		pos1, _ := portfolio.NewPosition("AAPL", 5, *pos1Price) 
		cost, _ := portfolio.NewMoney(pos1.PurchasePrice.Amount*int64(pos1.Shares), pos1.PurchasePrice.Currency)
		
		err := pFresh.AddPosition(*pos1, *cost)

		if err != nil {
			t.Fatalf("AddPosition() error = %v, wantErr nil", err)
		}
		if len(pFresh.Holdings) != 1 {
			t.Errorf("Holdings count = %d, want 1", len(pFresh.Holdings))
		}
		if _, ok := pFresh.Holdings["AAPL"]; !ok {
			t.Errorf("Holdings does not contain AAPL ticker")
		}
		expectedCash := initialCash.Amount - cost.Amount
		if pFresh.CashBalance.Amount != expectedCash {
			t.Errorf("CashBalance = %d, want %d", pFresh.CashBalance.Amount, expectedCash)
		}
		if pFresh.UpdatedAt.Equal(originalUpdatedAt) || pFresh.UpdatedAt.Before(originalUpdatedAt) {
			t.Errorf("UpdatedAt was not advanced. Initial: %v, Current: %v", originalUpdatedAt, pFresh.UpdatedAt)
		}
	})

	t.Run("InsufficientCash", func(t *testing.T) {
		pFresh, _ := portfolio.NewPortfolio("pFresh", portfolio.Aggressive, portfolio.Money{Amount: 100, Currency: "USD"}) 
		
		expensivePosPrice, _ := portfolio.NewMoney(5000, "USD") 
		expensivePos, _ := portfolio.NewPosition("TSLA", 10, *expensivePosPrice) 
		
		cost, _ := portfolio.NewMoney(expensivePos.PurchasePrice.Amount*int64(expensivePos.Shares), expensivePos.PurchasePrice.Currency)
		err := pFresh.AddPosition(*expensivePos, *cost)

		if err == nil {
			t.Errorf("AddPosition() with insufficient cash expected error, got nil")
		}
		if len(pFresh.Holdings) != 0 {
			t.Errorf("Holdings count should be 0 after failed add, got %d", len(pFresh.Holdings))
		}
	})
}

func TestPortfolio_RemovePosition(t *testing.T) {
	t.Run("SuccessfulRemove", func(t *testing.T) {
		testPortfolio, _ := portfolio.NewPortfolio("testRemove", portfolio.Conservative, portfolio.Money{Amount: 10000, Currency: "USD"})
		priceVal, _ := portfolio.NewMoney(500, "USD") 
		pos, _ := portfolio.NewPosition("MSFT", 10, *priceVal)
		cost, _ := portfolio.NewMoney(pos.PurchasePrice.Amount*int64(pos.Shares), pos.PurchasePrice.Currency)
		if err := testPortfolio.AddPosition(*pos, *cost); err != nil {
			t.Fatalf("Setup: AddPosition failed: %v", err)
		}

		originalCash := testPortfolio.CashBalance.Amount
		originalUpdatedAt := testPortfolio.UpdatedAt
		time.Sleep(1 * time.Millisecond) 

		sharesToRemove := 5
		proceedsFromSale, _ := portfolio.NewMoney(int64(sharesToRemove)*priceVal.Amount, "USD") 

		err := testPortfolio.RemovePosition("MSFT", sharesToRemove, *proceedsFromSale)
		if err != nil {
			t.Fatalf("RemovePosition() error = %v, wantErr nil", err)
		}

		if testPortfolio.CashBalance.Amount != originalCash+proceedsFromSale.Amount {
			t.Errorf("CashBalance after remove = %d, want %d", testPortfolio.CashBalance.Amount, originalCash+proceedsFromSale.Amount)
		}
		if testPortfolio.UpdatedAt.Equal(originalUpdatedAt) || testPortfolio.UpdatedAt.Before(originalUpdatedAt) {
			t.Errorf("UpdatedAt not advanced after RemovePosition. Initial: %v, Current: %v", originalUpdatedAt, testPortfolio.UpdatedAt)
		}
	})
}

func TestPortfolio_GenerateRebalanceRecommendations(t *testing.T) {
	p, _ := portfolio.NewPortfolio("test", portfolio.Moderate, portfolio.Money{Amount: 1000, Currency: "USD"})

	t.Run("RebalanceTriggered", func(t *testing.T) {
		p.LastRebalanceTime = time.Time{} 
		recs, err := p.GenerateRebalanceRecommendations()

		if err != nil {
			t.Fatalf("GenerateRebalanceRecommendations() error = %v, wantErr nil (for triggered rebalance)", err)
		}
		if len(recs) == 0 { 
			t.Errorf("Expected recommendations, got empty slice")
		}
		if recs[0] != "Recommendation: Sell X, Buy Y" { 
			t.Errorf("Unexpected recommendation content: %s", recs[0])
		}
	})

	t.Run("RebalanceNotTriggered", func(t *testing.T) {
		p.LastRebalanceTime = time.Now().Add(-10 * 24 * time.Hour) 
		_, err := p.GenerateRebalanceRecommendations()
		if err == nil {
			t.Errorf("GenerateRebalanceRecommendations() expected error for non-triggered rebalance, got nil")
		}
	})
}

func TestPortfolio_UpdateRiskProfile(t *testing.T) {
	p, _ := portfolio.NewPortfolio("test", portfolio.Conservative, portfolio.Money{Amount: 1000, Currency: "USD"})
	originalUpdatedAt := p.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	newProfile := portfolio.Aggressive
	p.UpdateRiskProfile(newProfile)

	if p.RiskProfile != newProfile {
		t.Errorf("RiskProfile not updated. Got %v, want %v", p.RiskProfile, newProfile)
	}
	if p.UpdatedAt.Equal(originalUpdatedAt) || p.UpdatedAt.Before(originalUpdatedAt) {
		t.Errorf("UpdatedAt not advanced after UpdateRiskProfile. Initial: %v, Current: %v", originalUpdatedAt, p.UpdatedAt)
	}
}

func TestNewPosition(t *testing.T) {
	price, _ := portfolio.NewMoney(15000, "USD") 
	t.Run("ValidPosition", func(t *testing.T) {
		pos, err := portfolio.NewPosition("GOOG", 10, *price)
		if err != nil {
			t.Fatalf("NewPosition() error = %v, wantErr nil", err)
		}
		if pos.CompanyTicker != "GOOG" {
			t.Errorf("Ticker = %s, want GOOG", pos.CompanyTicker)
		}
		if pos.Shares != 10 {
			t.Errorf("Shares = %d, want 10", pos.Shares)
		}
	})
	t.Run("EmptyTicker", func(t *testing.T) {
		_, err := portfolio.NewPosition("", 10, *price)
		if err == nil {
			t.Error("NewPosition() with empty ticker expected error, got nil")
		}
	})
	t.Run("NonPositiveShares", func(t *testing.T) {
		_, err := portfolio.NewPosition("MSFT", 0, *price)
		if err == nil {
			t.Error("NewPosition() with zero shares expected error, got nil")
		}
		_, err = portfolio.NewPosition("MSFT", -5, *price)
		if err == nil {
			t.Error("NewPosition() with negative shares expected error, got nil")
		}
	})
}

func TestNewMoney(t *testing.T) {
	t.Run("ValidMoney", func(t *testing.T) {
		m, err := portfolio.NewMoney(10050, "USD") 
		if err != nil {
			t.Fatalf("NewMoney() error = %v, wantErr nil", err)
		}
		if m.Amount != 10050 {
			t.Errorf("Amount = %d, want 10050", m.Amount)
		}
		if m.Currency != "USD" {
			t.Errorf("Currency = %s, want USD", m.Currency)
		}
	})
	t.Run("EmptyCurrency", func(t *testing.T) {
		_, err := portfolio.NewMoney(1000, "")
		if err == nil {
			t.Error("NewMoney() with empty currency expected error, got nil")
		}
	})
}

func TestMoney_Arithmetic(t *testing.T) {
	m1, _ := portfolio.NewMoney(100, "USD")
	m2, _ := portfolio.NewMoney(50, "USD")
	m3, _ := portfolio.NewMoney(30, "EUR")

	t.Run("AddSuccess", func(t *testing.T) {
		sum, err := m1.Add(*m2)
		if err != nil {
			t.Fatalf("Add() error = %v, wantErr nil", err)
		}
		if sum.Amount != 150 {
			t.Errorf("Add() Amount = %d, want 150", sum.Amount)
		}
		if sum.Currency != "USD" {
			t.Errorf("Add() Currency = %s, want USD", sum.Currency)
		}
	})
	t.Run("AddCurrencyMismatch", func(t *testing.T) {
		_, err := m1.Add(*m3)
		if err == nil {
			t.Error("Add() with currency mismatch expected error, got nil")
		}
	})
	t.Run("SubtractSuccess", func(t *testing.T) {
		diff, err := m1.Subtract(*m2)
		if err != nil {
			t.Fatalf("Subtract() error = %v, wantErr nil", err)
		}
		if diff.Amount != 50 {
			t.Errorf("Subtract() Amount = %d, want 50", diff.Amount)
		}
	})
	t.Run("SubtractCurrencyMismatch", func(t *testing.T) {
		_, err := m1.Subtract(*m3)
		if err == nil {
			t.Error("Subtract() with currency mismatch expected error, got nil")
		}
	})
}

func TestMoney_Checks(t *testing.T) {
	zero, _ := portfolio.NewMoney(0, "USD")
	positive, _ := portfolio.NewMoney(10, "USD")
	negative, _ := portfolio.NewMoney(-10, "USD")

	if !zero.IsZero() { t.Error("IsZero() failed for zero amount") }
	if positive.IsZero() { t.Error("IsZero() failed for positive amount") }
	if !positive.IsPositive() { t.Error("IsPositive() failed for positive amount") }
	if zero.IsPositive() { t.Error("IsPositive() failed for zero amount") }
	if !negative.IsNegative() { t.Error("IsNegative() failed for negative amount") }
	if zero.IsNegative() { t.Error("IsNegative() failed for zero amount") }
}

func TestRiskProfileEnum(t *testing.T) {
	testCases := []struct {
		profileVal portfolio.RiskProfile
		strVal     string
	}{
		{portfolio.Conservative, "Conservative"},
		{portfolio.Moderate, "Moderate"},
		{portfolio.Aggressive, "Aggressive"},
		{portfolio.UndefinedProfile, "UndefinedProfile"},
		{portfolio.RiskProfile(99), "UndefinedProfile"},
	}
	for _, tc := range testCases {
		t.Run("RiskProfileToString_"+tc.strVal, func(t *testing.T) {
			if str := tc.profileVal.String(); str != tc.strVal {
				t.Errorf("RiskProfile(%d).String() = %q, want %q", tc.profileVal, str, tc.strVal)
			}
		})
		if tc.profileVal != portfolio.UndefinedProfile && tc.profileVal <= portfolio.Aggressive {
			t.Run("ParseRiskProfile_"+tc.strVal, func(t *testing.T) {
				if p := portfolio.ParseRiskProfile(tc.strVal); p != tc.profileVal {
					t.Errorf("ParseRiskProfile(%q) = %v, want %v", tc.strVal, p, tc.profileVal)
				}
			})
		}
	}
	t.Run("ParseRiskProfile_Unknown", func(t *testing.T) {
		if p := portfolio.ParseRiskProfile("Unknown"); p != portfolio.UndefinedProfile {
			t.Errorf("ParseRiskProfile('Unknown') = %v, want UndefinedProfile", p)
		}
	})
}
