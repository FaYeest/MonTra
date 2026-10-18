package money_test

import (
	"errors"
	"math"
	"testing"

	"github.com/montra/montra/services/api/internal/domain/money"
)

func TestAmountArithmeticAndValidation(t *testing.T) {
	a, err := money.NewPositive(25000, "idr")
	if err != nil {
		t.Fatalf("expected valid amount, got %v", err)
	}
	if a.Currency() != "IDR" || a.Units() != 25000 {
		t.Fatalf("unexpected amount %+v", a)
	}

	b, err := money.NewPositive(15000, "IDR")
	if err != nil {
		t.Fatalf("expected valid amount, got %v", err)
	}

	sum, err := a.Add(b)
	if err != nil || sum.Units() != 40000 {
		t.Fatalf("expected 40000 IDR, got %v, err=%v", sum, err)
	}

	diff, err := a.Sub(b)
	if err != nil || diff.Units() != 10000 {
		t.Fatalf("expected 10000 IDR, got %v, err=%v", diff, err)
	}

	if _, err := money.NewPositive(0, "IDR"); !errors.Is(err, money.ErrNonPositiveAmount) {
		t.Errorf("expected ErrNonPositiveAmount for zero, got %v", err)
	}

	usd, _ := money.NewPositive(10, "USD")
	if _, err := a.Add(usd); !errors.Is(err, money.ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}

	maxAmt, _ := money.New(math.MaxInt64, "IDR")
	oneAmt, _ := money.New(1, "IDR")
	if _, err := maxAmt.Add(oneAmt); !errors.Is(err, money.ErrArithmeticOverflow) {
		t.Errorf("expected ErrArithmeticOverflow, got %v", err)
	}
}

func TestComputeTransferEffectDoesNotDoubleCountIncomeOrExpense(t *testing.T) {
	effect, err := money.ComputeTransferEffect("acc-bank", "acc-ewallet", 500000, 2500, "IDR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if effect.SourceDebit.Units() != 502500 {
		t.Errorf("expected source debit 502500, got %d", effect.SourceDebit.Units())
	}
	if effect.DestinationCredit.Units() != 500000 {
		t.Errorf("expected destination credit 500000, got %d", effect.DestinationCredit.Units())
	}
	if effect.ExternalIncome.Units() != 0 {
		t.Errorf("internal transfer must never count as external income, got %d", effect.ExternalIncome.Units())
	}
	if effect.ExternalExpense.Units() != 2500 {
		t.Errorf("only transfer fee should count as external expense, got %d", effect.ExternalExpense.Units())
	}

	// Net worth change across both accounts must equal -ExternalExpense
	netChange := effect.DestinationCredit.Units() - effect.SourceDebit.Units()
	if netChange != -effect.ExternalExpense.Units() {
		t.Errorf("reconciliation mismatch: netChange=%d, externalExpense=%d", netChange, effect.ExternalExpense.Units())
	}

	if _, err := money.ComputeTransferEffect("acc-1", "acc-1", 10000, 0, "IDR"); !errors.Is(err, money.ErrSameTransferAccount) {
		t.Errorf("expected ErrSameTransferAccount, got %v", err)
	}
}
