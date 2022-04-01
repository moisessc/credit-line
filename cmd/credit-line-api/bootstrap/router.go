package bootstrap

import (
	"net/http"

	pv "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"

	"credit-line/internal/controller"
	"credit-line/pkg/middleware"
	"credit-line/pkg/validator"
)

// newEchoRouter builds an instance of the echo router
func newEchoRouter(clh *controller.CreditLineHandler) http.Handler {
	e := echo.New()
	e.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format:           "${time_custom} ip=${remote_ip} method=${method}, uri=${uri}, status=${status} latency=${latency_human} \n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
	e.Validator = validator.New(pv.New())

	products := e.Group("/api/v1/credits")
	products.POST("/calculate/limit",
		clh.CreditLine, middleware.ValidateRetries(), middleware.IpRateLimitByTime(), middleware.IpRateLimitByFail())

	return e
}
