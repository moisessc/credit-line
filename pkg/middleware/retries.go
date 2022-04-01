package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"credit-line/pkg/cache"
	"credit-line/pkg/env"
)

// ValidateRetries middleware that provides a validation retries
func ValidateRetries() echo.MiddlewareFunc {
	cfg := env.RetrieveEnvVariables().Middlewares
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		cacheRequest := cache.RetrieveRequestCache().RequestFailed
		return func(c echo.Context) (err error) {
			ipRetries := cacheRequest[c.RealIP()]
			if ipRetries >= cfg.DeclineRetriesAllowed {
				return &echo.HTTPError{
					Code:     middleware.ErrRateLimitExceeded.Code,
					Message:  cfg.DeclineRetriesMessage,
					Internal: err,
				}
			}
			return next(c)
		}
	}
}
