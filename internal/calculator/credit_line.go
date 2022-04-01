package calculator

import (
	"context"
	"credit-line/pkg/env"
	"credit-line/pkg/errors"
)

const (
	// SME_FOUNDING_TYPE the SME founding type
	SME_FOUNDING_TYPE = "SME"
	// STARTUP_FOUNDING_TYPE the StartUp founding type
	STARTUP_FOUNDING_TYPE = "Startup"
)

// CreditLineCalculator calculator contracts for the credit line
type CreditLineCalculator interface {
	CalculateCreditLine(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error)
}

// creditLine struct that implement the CreditLineCalculator interface
type creditLine struct {
	ratios *env.Ratios
}

// NewCreditLine creates a new pointer of CreditLine struct
func NewCreditLine(ratios *env.Ratios) *creditLine {
	return &creditLine{
		ratios: ratios,
	}
}

// CalculateCreditLine implement the interface CreditLineCalculator.CalculateCreditLine
func (cl *creditLine) CalculateCreditLine(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error) {
	switch foundingType {
	case SME_FOUNDING_TYPE:
		amount := cashBalance / cl.ratios.CashBalance
		return amount, nil
	case STARTUP_FOUNDING_TYPE:
		amountCb := cashBalance / cl.ratios.CashBalance
		amountMr := monthlyRevenue / cl.ratios.MonthlyRevenue
		if amountCb > amountMr {
			return amountCb, nil
		}
		return amountMr, nil
	default:
		return 0, errors.ErrInvalidFoundingType
	}
}
