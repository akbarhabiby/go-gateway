package middlewares

import (
	"fmt"
	"strconv"
	"time"

	"github.com/akbarhabiby/go-gateway/config"
	"github.com/akbarhabiby/go-gateway/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var uid string
			u, err := uuid.NewUUID()
			if err != nil {
				uid = strconv.FormatInt(time.Now().Unix(), 10)
			} else {
				uid = u.String()
			}
			ctxLogger := log.LogContextModel{
				ServiceName:    config.Config.GetString("app.name"),
				ServiceVersion: config.Config.GetString("app.version"),
				ServicePort:    config.Config.GetInt("app.port"),
				ThreadID:       uid,
				ReqMethod:      c.Request().Method,
				ReqURI:         c.Request().URL.String(),
			}

			request := c.Request()

			ctx := log.SetContextFromEchoRequest(c)
			ctx = log.InjectCtx(ctx, ctxLogger)
			c.SetRequest(request.WithContext(ctx))

			return next(c)
		}
	}
}

func TDRLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		ctx := c.Request().Context()
		if isLoggingSkipForProxy(c) {
			rewrite := log.GetProxyProcessingFromEchoContext(c).(map[string]string)
			url := log.GetProxyURLFromEchoContext(c).(string)
			targetURL := log.GetProxyTargetURLFromEchoContext(c).(string)
			log.Info(ctx, fmt.Sprintf("[PROXY] %s '%s' to '%s", c.Request().Method, url, targetURL), rewrite)
			return
		}

		if c.Get("skip-body-logging") != nil {
			log.TDR(ctx, []byte{}, resBody)
			return
		}

		log.TDR(ctx, reqBody, resBody)
	})
}
