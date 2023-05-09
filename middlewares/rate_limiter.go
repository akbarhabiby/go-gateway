package middlewares

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter() echo.MiddlewareFunc {
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  60,
	}
	limiter := limiter.New(memory.NewStore(), rate)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			identifier := c.RealIP()

			l, _ := limiter.Get(c.Request().Context(), identifier)

			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(l.Limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(l.Remaining, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(l.Reset, 10))

			if l.Reached {
				c.Set("request-limit-reached", true)
				err = fmt.Errorf("too many requests on %s", c.Request().URL.String())
				return
			}

			return next(c)
		}
	}
}
