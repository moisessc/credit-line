package calculator

import (
	"context"
	"credit-line/pkg/env"
	"credit-line/pkg/errors"
	"reflect"
	"testing"
)

var ratios *env.Ratios = &env.Ratios{
	CashBalance:    3,
	MonthlyRevenue: 5,
}

func Test_Calculate_Credit_Line_Calculator(t *testing.T) {
	testCases := map[string]struct {
		params struct {
			ctx            context.Context
			foundingType   string
			cashBalance    float64
			monthlyRevenue float64
		}
		expectedLineOfCredit float64
		expectedError        error
	}{
		"credit_line_could_not_be_calculated": {
			params: struct {
				ctx            context.Context
				foundingType   string
				cashBalance    float64
				monthlyRevenue float64
			}{
				ctx:            context.Background(),
				foundingType:   "SMA",
				cashBalance:    435.30,
				monthlyRevenue: 4235.45,
			},
			expectedLineOfCredit: 0,
			expectedError:        errors.ErrInvalidFoundingType,
		},
		"credit_line_calculated_to_SME": {
			params: struct {
				ctx            context.Context
				foundingType   string
				cashBalance    float64
				monthlyRevenue float64
			}{
				ctx:            context.Background(),
				foundingType:   "SME",
				cashBalance:    435.30,
				monthlyRevenue: 4235.45,
			},
			expectedLineOfCredit: 145.10,
			expectedError:        nil,
		},
		"credit_line_calculated_to_Startup_by_monthly_revenue": {
			params: struct {
				ctx            context.Context
				foundingType   string
				cashBalance    float64
				monthlyRevenue float64
			}{
				ctx:            context.Background(),
				foundingType:   "Startup",
				cashBalance:    435.30,
				monthlyRevenue: 4235.45,
			},
			expectedLineOfCredit: 847.0899999999999,
			expectedError:        nil,
		},
		"credit_line_calculated_to_Startup_by_cash_balance": {
			params: struct {
				ctx            context.Context
				foundingType   string
				cashBalance    float64
				monthlyRevenue float64
			}{
				ctx:            context.Background(),
				foundingType:   "Startup",
				cashBalance:    13435.30,
				monthlyRevenue: 4235.45,
			},
			expectedLineOfCredit: 4478.433333333333,
			expectedError:        nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			calculator := NewCreditLine(ratios)
			got, err := calculator.CalculateCreditLine(tc.params.ctx, tc.params.foundingType,
				tc.params.cashBalance, tc.params.monthlyRevenue)

			if tc.expectedError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.expectedError != nil && err == nil {
				t.Fatalf("got nil error expecting: %v", tc.expectedError)
			}

			if tc.expectedError != nil && tc.expectedError.Error() != err.Error() {
				t.Fatalf("unexpected error got: %v expected: %v", err, tc.expectedError)
			}

			if !reflect.DeepEqual(tc.expectedLineOfCredit, got) {
				t.Fatalf("unexpected result, got: %v, expected: %v", got, tc.expectedLineOfCredit)
			}
		})
	}
}
