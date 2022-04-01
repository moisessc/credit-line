package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	pv "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"credit-line/internal/model"
	"credit-line/internal/service"
	"credit-line/pkg/errors"
	"credit-line/pkg/validator"
)

type mockCreditLineService struct {
	determineCreditLimit func(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error)
}

func (mcls *mockCreditLineService) DetermineCreditLimit(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error) {
	return mcls.determineCreditLimit(ctx, ip, creditLine)
}

func Test_Determine_Credit_Limit_Controller(t *testing.T) {
	testCases := map[string]struct {
		service            service.CreditLineService
		request            []byte
		expectedBody       interface{}
		expectedStatusCode int
	}{
		"unmarshal_error": {
			service: &mockCreditLineService{
				determineCreditLimit: func(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error) {
					return nil, nil
				},
			},
			request: []byte(`{
				"foundingType": "SME",
				"cashBalance": "435.30",
				"monthlyRevenue": 4235.45,
				"requestedCreditLine": 100,
				"requestedDate": "2021-07-19T16:32:59.860Z"
			}`),

			expectedBody: errors.ApiResponse{
				Message: "unmarshal error data type, got: string, expected: number in cashBalance param",
				Code:    "INVALID_REQUEST",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		"validation_error": {
			service: &mockCreditLineService{
				determineCreditLimit: func(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error) {
					return nil, nil
				},
			},
			request: []byte(`{
				"cashBalance": 435.30,
				"requestedCreditLine": 100,
				"requestedDate": "2021-07-19T16:32:59.860Z"
			}`),

			expectedBody: errors.ApiResponse{
				Message: "malformed request, please check the following parameters in the request: [foundingType, monthlyRevenue]",
				Code:    "INVALID_REQUEST",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		"credit_line_could_not_be_determined": {
			service: &mockCreditLineService{
				determineCreditLimit: func(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error) {
					return nil, fmt.Errorf("determination failed: %w", errors.ErrInvalidFoundingType)
				},
			},
			request: []byte(`{
				"foundingType": "SME",
				"cashBalance": 435.30,
				"monthlyRevenue": 4235.45,
				"requestedCreditLine": 100,
				"requestedDate": "2021-07-19T16:32:59.860Z"
			}`),
			expectedBody: errors.ApiResponse{
				Message: "determination failed: invalid foundingType",
				Code:    "INVALID_REQUEST",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		"credit_line_approved": {
			service: &mockCreditLineService{
				determineCreditLimit: func(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error) {
					return model.NewCreditLineResponse(model.Approved, "145.10"), nil
				},
			},
			request: []byte(`{
				"foundingType": "SME",
				"cashBalance": 435.30,
				"monthlyRevenue": 4235.45,
				"requestedCreditLine": 100,
				"requestedDate": "2021-07-19T16:32:59.860Z"
			}`),
			expectedBody:       model.NewCreditLineResponse(model.Approved, "145.10"),
			expectedStatusCode: http.StatusOK,
		},
		"credit_line_declined": {
			service: &mockCreditLineService{
				determineCreditLimit: func(ctx context.Context, ip string, creditLine *model.CreditLine) (*model.CreditLineResponse, error) {
					return model.NewCreditLineResponse(model.Declined, "0.00"), nil
				},
			},
			request: []byte(`{
				"foundingType": "SME",
				"cashBalance": 435.30,
				"monthlyRevenue": 4235.45,
				"requestedCreditLine": 1000,
				"requestedDate": "2021-07-19T16:32:59.860Z"
			}`),
			expectedBody:       model.NewCreditLineResponse(model.Declined, "0.00"),
			expectedStatusCode: http.StatusOK,
		},
	}

	e := echo.New()
	e.Validator = validator.New(pv.New())
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tc.request))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			ctx := e.NewContext(r, w)

			handler := NewCreditLineHandler(tc.service)
			err := handler.CreditLine(ctx)
			if err != nil {
				t.Errorf("unexpected error, got: %v", err)
			}

			gotStatusCode := w.Code
			if tc.expectedStatusCode != gotStatusCode {
				t.Errorf("unexpected status code, got: %v, expected: %v", gotStatusCode, tc.expectedStatusCode)
			}

			if name != "credit_line_approved" && name != "credit_line_declined" {
				var gotBody errors.ApiResponse
				err = json.NewDecoder(w.Body).Decode(&gotBody)
				if err != nil {
					t.Errorf("unexpected unmarshall error, got: %v", err)
				}

				if !reflect.DeepEqual(tc.expectedBody, gotBody) {
					t.Errorf("unexpected response, got: %v, expected: %v", gotBody, tc.expectedBody)
				}
			} else {
				var gotBody model.CreditLineResponse
				err = json.NewDecoder(w.Body).Decode(&gotBody)
				if err != nil {
					t.Errorf("unexpected unmarshall error, got: %v", err)
				}

				if !reflect.DeepEqual(tc.expectedBody, &gotBody) {
					t.Errorf("unexpected response, got: %v, expected: %v", gotBody, tc.expectedBody)
				}
			}
		})
	}
}
