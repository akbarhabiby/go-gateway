package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BodyLimit() echo.MiddlewareFunc {
	return middleware.BodyLimit("10M")
}
