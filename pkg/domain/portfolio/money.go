package portfolio

import "errors"

// Money represents a monetary value, including currency.
// This is a value object.
type Money struct {
	Amount   int64  // Amount in the smallest currency unit (e.g., cents for USD)
	Currency string // Currency code (e.g., "USD", "EUR")
}

// NewMoney creates a new Money instance.
func NewMoney(amount int64, currency string) (*Money, error) {
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}
	// Potentially add more validation for currency codes if a specific list is supported.
	return &Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

// Add returns a new Money object representing the sum of m and other.
// It returns an error if the currencies do not match.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("currency mismatch")
	}
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}

// Subtract returns a new Money object representing the difference of m and other.
// It returns an error if the currencies do not match.
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("currency mismatch")
	}
	return Money{Amount: m.Amount - other.Amount, Currency: m.Currency}, nil
}

// IsZero checks if the monetary amount is zero.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsPositive checks if the monetary amount is positive.
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// IsNegative checks if the monetary amount is negative.
func (m Money) IsNegative() bool {
	return m.Amount < 0
}
