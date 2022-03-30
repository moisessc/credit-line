package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// validatorHandler struct for the custom request validator handler
type validatorHandler struct {
	validator *validator.Validate
}

// creditLineFields variable to identify the product attribute with the json tag
var creditLineFields = map[string]string{
	"FoundingType":        "foundingType",
	"CashBalance":         "cashBalance",
	"MonthlyRevenue":      "monthlyRevenue",
	"RequestedCreditLine": "requestedCreditLine",
	"RequestedDate":       "requestedDate",
}

// New creates a new instance of validatorHandler struct
func New(validator *validator.Validate) *validatorHandler {
	return &validatorHandler{
		validator: validator,
	}
}

// Validate validates an API requests
func (vh *validatorHandler) Validate(i interface{}) error {
	return vh.validator.Struct(i)
}

// RetrieveValidationErrorMessage retrieves a custom error message when the validaton fails
func RetrieveValidationErrorMessage(err error) string {
	fields := make([]string, 0, 5)
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, v := range err {
			fields = append(fields, fmt.Sprint(creditLineFields[v.Field()], ","))
		}
		lastField := fields[len(fields)-1]
		fields[len(fields)-1] = lastField[:len(lastField)-1]
	}
	return fmt.Sprintf("malformed request, please check the following parameters in the request: %v", fields)
}
