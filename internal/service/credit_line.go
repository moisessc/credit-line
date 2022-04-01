package service

import (
	"context"
	"fmt"

	"credit-line/internal/calculator"
	"credit-line/internal/model"
	"credit-line/pkg/cache"
)

// CreditLineService services contracts for the credit line entity
type CreditLineService interface {
	DetermineCreditLimit(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error)
}

// creditLine struct that implement the CreditLineService interface
type creditLine struct {
	calculator calculator.CreditLineCalculator
}

// NewCreditLine creates a new pointer of creditLine struct
func NewCreditLine(calculator calculator.CreditLineCalculator) *creditLine {
	return &creditLine{
		calculator: calculator,
	}
}

// DetermineCreditLimit implement the interface CreditLineService.DetermineCreditLimit
func (cl *creditLine) DetermineCreditLimit(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error) {
	amount, err := cl.calculator.CalculateCreditLine(ctx, creditLine.FoundingType(), creditLine.CashBalance(), creditLine.MonthlyRevenue())
	if err != nil {
		return nil, fmt.Errorf("determination failed: %w", err)
	}

	if amount > creditLine.RequestedCreditLine() {
		cache.UpdateRequestCache(model.Approved, ip)
		return model.NewCreditLineResponse(model.Approved, fmt.Sprintf("%.2f", amount)), nil
	}

	cache.UpdateRequestCache(model.Declined, ip)
	return model.NewCreditLineResponse(model.Declined, "0.00"), nil
}
