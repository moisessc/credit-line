package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"credit-line/internal/model"
	"credit-line/internal/service"
	"credit-line/pkg/errors"
)

// CreditLineHandler struct that contains the service for the CreditLine entity
type CreditLineHandler struct {
	service service.CreditLineService
}

// CreditLineRequest struct that represents the CreditLine request
type CreditLineRequest struct {
	FoundingType        string  `json:"foundingType" validate:"required"`
	CashBalance         float64 `json:"cashBalance" validate:"required"`
	MonthlyRevenue      float64 `json:"monthlyRevenue" validate:"required"`
	RequestedCreditLine float64 `json:"requestedCreditLine" validate:"required"`
	RequestedDate       string  `json:"requestedDate" validate:"required"`
}

// NewCreditLineHandler creates a new pointer of CreditLineHandler struct
func NewCreditLineHandler(service service.CreditLineService) *CreditLineHandler {
	return &CreditLineHandler{
		service: service,
	}
}

// CreditLine invokes the echo handler to calculate the credit line
func (clh *CreditLineHandler) CreditLine(c echo.Context) error {
	var request CreditLineRequest

	if err := c.Bind(&request); err != nil {
		errResponse, _ := errors.MapError(err, errors.UnmarshallErr)
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	if err := c.Validate(request); err != nil {
		errResponse, _ := errors.MapError(err, errors.ValidationErr)
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	creditLine := model.NewCreditLine(request.FoundingType, request.RequestedDate,
		request.CashBalance, request.MonthlyRevenue, request.RequestedCreditLine)

	creditLineResponse, err := clh.service.DetermineCreditLimit(c.Request().Context(), creditLine)
	if err != nil {
		errResponse, code := errors.MapError(err, errors.DomainErr)
		return c.JSON(code, errResponse)
	}
	return c.JSON(http.StatusOK, creditLineResponse)
}
