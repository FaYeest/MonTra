package money

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	// ErrInvalidCurrency indicates that the currency code is not a valid 3-letter uppercase ISO 4217 code.
	ErrInvalidCurrency = errors.New("money: currency must be a 3-letter uppercase ISO 4217 code")
	// ErrCurrencyMismatch indicates an attempt to combine monetary amounts of different currencies.
	ErrCurrencyMismatch = errors.New("money: cannot operate on different currencies")
	// ErrNonPositiveAmount indicates that a transaction or transfer principal amount is <= 0.
	ErrNonPositiveAmount = errors.New("money: transaction amount must be greater than zero")
	// ErrNegativeFee indicates that a transfer fee is < 0.
	ErrNegativeFee = errors.New("money: transfer fee cannot be negative")
	// ErrSameTransferAccount indicates that a transfer source and destination account are identical.
	ErrSameTransferAccount = errors.New("money: source and destination accounts must be distinct")
	// ErrArithmeticOverflow indicates integer overflow during deterministic monetary calculation.
	ErrArithmeticOverflow = errors.New("money: integer overflow in monetary calculation")
)

// Amount represents an authoritative monetary value using exact 64-bit integer arithmetic
// (e.g., whole Rupiah for IDR) paired with an explicit ISO 4217 currency code.
// Floating-point numbers are never used for storage or calculation.
type Amount struct {
	units    int64
	currency string
}

// New creates a validated Amount.
func New(units int64, currency string) (Amount, error) {
	norm, err := ValidateCurrency(currency)
	if err != nil {
		return Amount{}, err
	}
	return Amount{units: units, currency: norm}, nil
}

// NewPositive creates a validated Amount that must be strictly greater than zero (> 0),
// suitable for income, expense, or transfer principal values.
func NewPositive(units int64, currency string) (Amount, error) {
	if units <= 0 {
		return Amount{}, ErrNonPositiveAmount
	}
	return New(units, currency)
}

// ValidateCurrency verifies that a currency string is a 3-letter ASCII code.
func ValidateCurrency(currency string) (string, error) {
	c := strings.TrimSpace(strings.ToUpper(currency))
	if len(c) != 3 {
		return "", ErrInvalidCurrency
	}
	for i := 0; i < 3; i++ {
		if c[i] < 'A' || c[i] > 'Z' {
			return "", ErrInvalidCurrency
		}
	}
	return c, nil
}

// Units returns the exact integer units of the monetary amount.
func (a Amount) Units() int64 {
	return a.units
}

// Currency returns the 3-letter ISO 4217 currency code.
func (a Amount) Currency() string {
	return a.currency
}

// Add returns a + other with currency and integer-overflow checks.
func (a Amount) Add(other Amount) (Amount, error) {
	if a.currency != other.currency {
		return Amount{}, ErrCurrencyMismatch
	}
	if (other.units > 0 && a.units > math.MaxInt64-other.units) ||
		(other.units < 0 && a.units < math.MinInt64-other.units) {
		return Amount{}, ErrArithmeticOverflow
	}
	return Amount{units: a.units + other.units, currency: a.currency}, nil
}

// Sub returns a - other with currency and integer-overflow checks.
func (a Amount) Sub(other Amount) (Amount, error) {
	if a.currency != other.currency {
		return Amount{}, ErrCurrencyMismatch
	}
	if (other.units < 0 && a.units > math.MaxInt64+other.units) ||
		(other.units > 0 && a.units < math.MinInt64+other.units) {
		return Amount{}, ErrArithmeticOverflow
	}
	return Amount{units: a.units - other.units, currency: a.currency}, nil
}

// String formats the amount with its currency code.
func (a Amount) String() string {
	return fmt.Sprintf("%s %d", a.currency, a.units)
}

// TransferEffect represents the deterministic accounting impact of an internal transfer
// between two user-owned accounts. Principal transfer amounts are never counted as
// external income or expense; only an optional fee is recorded as an external expense.
type TransferEffect struct {
	SourceDebit       Amount
	DestinationCredit Amount
	ExternalExpense   Amount
	ExternalIncome    Amount
}

// ComputeTransferEffect validates and calculates the exact account and cash-flow deltas
// for an internal transfer.
func ComputeTransferEffect(sourceAccountID, destAccountID string, principalUnits, feeUnits int64, currency string) (TransferEffect, error) {
	if strings.TrimSpace(sourceAccountID) == "" || strings.TrimSpace(sourceAccountID) == strings.TrimSpace(destAccountID) {
		return TransferEffect{}, ErrSameTransferAccount
	}
	principal, err := NewPositive(principalUnits, currency)
	if err != nil {
		return TransferEffect{}, err
	}
	if feeUnits < 0 {
		return TransferEffect{}, ErrNegativeFee
	}
	fee, err := New(feeUnits, currency)
	if err != nil {
		return TransferEffect{}, err
	}
	sourceDebit, err := principal.Add(fee)
	if err != nil {
		return TransferEffect{}, err
	}
	zeroIncome, _ := New(0, principal.Currency())
	return TransferEffect{
		SourceDebit:       sourceDebit,
		DestinationCredit: principal,
		ExternalExpense:   fee,
		ExternalIncome:    zeroIncome,
	}, nil
}
