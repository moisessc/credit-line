package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"credit-line/internal/model"
	"credit-line/pkg/cache"
	"credit-line/pkg/env"
)

// IpRateLimitByTime middleware that provides a rate limit validation by time
func IpRateLimitByTime() echo.MiddlewareFunc {
	cfg := env.RetrieveEnvVariables().Middlewares
	rate := limiter.Rate{
		Period: time.Duration(cfg.ApprovedRateLimitTime) * time.Second,
		Limit:  cfg.ApprovedRateLimitRequest,
	}
	store := memory.NewStore()
	ipRateLimiter := limiter.New(store, rate)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		cacheRequest := cache.RetrieveRequestCache()
		return func(c echo.Context) (err error) {
			ip := c.RealIP()
			limiterCtx, err := ipRateLimiter.Get(c.Request().Context(), ip)
			if cacheRequest.CurrentCreditStatus == string(model.Approved) {
				if err != nil {
					log.Printf("IpRateLimitByTime err: %v, %s on %s", err, ip, c.Request().URL)
					return &echo.HTTPError{
						Code:     middleware.ErrExtractorError.Code,
						Message:  middleware.ErrExtractorError.Message,
						Internal: err,
					}
				}

				if limiterCtx.Reached {
					log.Printf("Many Requests from %s on %s", ip, c.Request().URL)
					return &echo.HTTPError{
						Code:     middleware.ErrRateLimitExceeded.Code,
						Message:  middleware.ErrRateLimitExceeded.Message,
						Internal: err,
					}
				}
			}
			return next(c)
		}
	}
}

// IpRateLimitByFail middleware that provides a rate limit validation bt time when the request is declined
func IpRateLimitByFail() echo.MiddlewareFunc {
	cfg := env.RetrieveEnvVariables().Middlewares
	rate := limiter.Rate{
		Period: time.Duration(cfg.DeclineRateLimitTime) * time.Second,
		Limit:  cfg.DeclineRateLimitRequest,
	}
	store := memory.NewStore()
	ipRateLimiter := limiter.New(store, rate)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		cacheRequest := cache.RetrieveRequestCache()
		return func(c echo.Context) (err error) {
			ip := c.RealIP()
			limiterCtx, err := ipRateLimiter.Get(c.Request().Context(), ip)
			if cacheRequest.CurrentCreditStatus == string(model.Declined) {
				if err != nil {
					log.Printf("IpRateLimitByFail err: %v, %s on %s", err, ip, c.Request().URL)
					return &echo.HTTPError{
						Code:     middleware.ErrExtractorError.Code,
						Message:  middleware.ErrExtractorError.Message,
						Internal: err,
					}
				}

				if limiterCtx.Reached {
					log.Printf("Many Requests from %s on %s", ip, c.Request().URL)
					return &echo.HTTPError{
						Code:     middleware.ErrRateLimitExceeded.Code,
						Message:  middleware.ErrRateLimitExceeded.Message,
						Internal: err,
					}
				}
			}
			return next(c)
		}
	}
}
