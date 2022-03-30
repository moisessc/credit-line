package bootstrap

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// newEchoRouter builds an instance of the echo router
func newEchoRouter() http.Handler {
	e := echo.New()

	products := e.Group("/api/v1/credits")
	products.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello world")
	})

	return e
}
