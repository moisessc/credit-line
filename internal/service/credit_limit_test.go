package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"credit-line/internal/calculator"
	"credit-line/internal/model"
	"credit-line/pkg/errors"
)

type mockCreditLineCalculator struct {
	calculateCreditLine func(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error)
}

func (mclc *mockCreditLineCalculator) CalculateCreditLine(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error) {
	return mclc.calculateCreditLine(ctx, foundingType, cashBalance, monthlyRevenue)
}

func Test_Determine_Credit_Limit_Service(t *testing.T) {
	testCases := map[string]struct {
		calculator calculator.CreditLineCalculator
		params     struct {
			ctx        context.Context
			ip         string
			creditLine *model.CreditLine
		}
		expectedResponse *model.CreditLineResponse
		expectedError    error
	}{
		"credit_line_could_not_be_determined": {
			calculator: &mockCreditLineCalculator{
				func(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error) {
					return 0, errors.ErrInvalidFoundingType
				},
			},
			params: struct {
				ctx        context.Context
				ip         string
				creditLine *model.CreditLine
			}{
				ctx:        context.Background(),
				ip:         "167.222.20.251",
				creditLine: model.NewCreditLine("SME", "2021-07-19T16:32:59.860Z", 435.30, 4235.45, 100),
			},
			expectedError: fmt.Errorf("determination failed: %w", errors.ErrInvalidFoundingType),
		},
		"credit_line_approved_SME": {
			calculator: &mockCreditLineCalculator{
				func(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error) {
					return 145.10, nil
				},
			},
			params: struct {
				ctx        context.Context
				ip         string
				creditLine *model.CreditLine
			}{
				ctx:        context.Background(),
				ip:         "167.222.20.251",
				creditLine: model.NewCreditLine("SME", "2021-07-19T16:32:59.860Z", 435.30, 4235.45, 100),
			},
			expectedResponse: model.NewCreditLineResponse(model.Approved, "145.10"),
		},
		"credit_line_declined_SME": {
			calculator: &mockCreditLineCalculator{
				func(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error) {
					return 145.10, nil
				},
			},
			params: struct {
				ctx        context.Context
				ip         string
				creditLine *model.CreditLine
			}{
				ctx:        context.Background(),
				ip:         "167.222.20.251",
				creditLine: model.NewCreditLine("SME", "2021-07-19T16:32:59.860Z", 435.30, 4235.45, 1000),
			},
			expectedResponse: model.NewCreditLineResponse(model.Declined, "0.00"),
		},
		"credit_line_approved_Startup": {
			calculator: &mockCreditLineCalculator{
				func(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error) {
					return 847.09, nil
				},
			},
			params: struct {
				ctx        context.Context
				ip         string
				creditLine *model.CreditLine
			}{
				ctx:        context.Background(),
				ip:         "167.222.20.251",
				creditLine: model.NewCreditLine("Startup", "2021-07-19T16:32:59.860Z", 435.30, 4235.45, 100),
			},
			expectedResponse: model.NewCreditLineResponse(model.Approved, "847.09"),
		},
		"credit_line_declined_Startup": {
			calculator: &mockCreditLineCalculator{
				func(ctx context.Context, foundingType string, cashBalance, monthlyRevenue float64) (float64, error) {
					return 847.09, nil
				},
			},
			params: struct {
				ctx        context.Context
				ip         string
				creditLine *model.CreditLine
			}{
				ctx:        context.Background(),
				ip:         "167.222.20.251",
				creditLine: model.NewCreditLine("Startup", "2021-07-19T16:32:59.860Z", 435.30, 4235.45, 1000),
			},
			expectedResponse: model.NewCreditLineResponse(model.Declined, "0.00"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			service := NewCreditLine(tc.calculator)
			got, err := service.DetermineCreditLimit(tc.params.ctx, tc.params.ip, tc.params.creditLine)

			if tc.expectedError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.expectedError != nil && err == nil {
				t.Fatalf("got nil error expecting: %v", tc.expectedError)
			}

			if tc.expectedError != nil && tc.expectedError.Error() != err.Error() {
				t.Fatalf("unexpected error got: %v expected: %v", err, tc.expectedError)
			}

			if !reflect.DeepEqual(tc.expectedResponse, got) {
				t.Fatalf("unexpected result, got: %v, expected: %v", got, tc.expectedResponse)
			}
		})
	}
}
