package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"credit-line/pkg/validator"
)

const (
	// UnmarshallErr identify an unmarshall error type
	UnmarshallErr ErrorType = "UNMARSHALL_ERROR"
	// ValidationErr identify a validation error type
	ValidationErr ErrorType = "VALIDATION_ERROR"
	// DomainErr identify a domain error type
	DomainErr ErrorType = "DOMAIN_ERROR"

	// internalServerErrorCode code to represent an internal server error
	internalServerErrorCode = "INTERNAL_SERVER_ERROR"
	// invalidRequestCode code to represent an invalid request
	invalidRequestCode = "INVALID_REQUEST"
)

// ErrorType type to specify an error type
type ErrorType string

// ApiResponse struct for error responses in the API
type ApiResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// MapError transform an error into custom error response
func MapError(err error, errType ErrorType) (*ApiResponse, int) {
	var msg, code string
	var statusCode int
	switch errType {
	case UnmarshallErr:
		msg = retrieveUnmarshalErrorMessage(err)
		code = invalidRequestCode
	case ValidationErr:
		msg = validator.RetrieveValidationErrorMessage(err)
		code = invalidRequestCode
	case DomainErr:
		msg = err.Error()
		statusCode, code = retrieveDomainErrorCode(err)
	}

	return &ApiResponse{Message: msg, Code: code}, statusCode
}

// retrieveUnmarshalErrorInformation retrieves the information when the bind method fails
func retrieveUnmarshalErrorMessage(err error) string {
	var field, expected, got string
	if err, ok := err.(*echo.HTTPError); ok {
		ierr := err.Internal
		if ute, ok := ierr.(*json.UnmarshalTypeError); ok {
			field = ute.Field
			expected = ute.Type.Name()
			got = ute.Value
		}
	}
	if strings.Contains(expected, "int") || strings.Contains(expected, "float") {
		expected = "number"
	}
	return fmt.Sprint("unmarshal error data type, got: ", got, ", expected: ", expected, " in ", field, " param")
}

// retrieveDomainErrorCode retrieves the error code of one domain error
func retrieveDomainErrorCode(err error) (int, string) {
	switch {
	case errors.Is(err, ErrInvalidFoundingType):
		return http.StatusBadRequest, invalidRequestCode
	default:
		return http.StatusInternalServerError, internalServerErrorCode
	}
}
