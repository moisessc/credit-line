package bootstrap

import (
	"net/http"

	pv "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"credit-line/internal/controller"
	"credit-line/pkg/middleware"
	"credit-line/pkg/validator"
)

// newEchoRouter builds an instance of the echo router
func newEchoRouter(clh *controller.CreditLineHandler) http.Handler {
	e := echo.New()
	e.Validator = validator.New(pv.New())

	products := e.Group("/api/v1/credits")
	products.GET("/calculate/limit",
		clh.CreditLine, middleware.ValidateRetries(), middleware.IpRateLimitByTime(), middleware.IpRateLimitByFail())

	return e
}
