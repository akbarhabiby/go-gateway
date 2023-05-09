package middlewares

import (
	"github.com/akbarhabiby/go-gateway/log"
	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo) *echo.Echo {
	e.Use(CORS())
	e.Use(RateLimiter())
	e.Use(BodyLimit())
	e.Use(RequestLoggerMiddleware())
	e.Use(TDRLoggerMiddleware())
	e.HTTPErrorHandler = ErrorHandler
	return e
}

func isLoggingSkipForProxy(c echo.Context) bool {
	return log.GetProxyProcessingFromEchoContext(c) != nil
}
