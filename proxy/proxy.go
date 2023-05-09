package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/akbarhabiby/go-gateway/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ProxyMiddlewareConfig struct {
	// Name
	Name string

	// Current Server Group Prefix
	Prefix string

	// Target URL
	URL string

	// Timeout
	Timeout time.Duration

	// Rewrite defines URL path rewrite rules. The values captured in asterisk can be
	// retrieved by index e.g. $1, $2 and so on.
	// Examples:
	// "/old":              "/new",
	// "/api/*":            "/$1",
	// "/js/*":             "/public/javascripts/$1",
	// "/users/*/orders/*": "/user/$1/order/$2",
	Rewrite map[string]string

	// ModifyResponse defines function to modify response from ProxyTarget.
	ModifyResponse func(*http.Response) error
}

func Setup(e *echo.Echo, configs []ProxyMiddlewareConfig) *echo.Echo {
	ctx := context.Background()
	for _, config := range configs {
		url, err := url.Parse(config.URL)
		if err != nil {
			panic(err)
		}
		proxyConfig := middleware.ProxyConfig{
			Skipper:        middleware.DefaultSkipper,
			Balancer:       middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{Name: config.Name, URL: url}}),
			Rewrite:        config.Rewrite,
			ModifyResponse: config.ModifyResponse,
		}
		if config.Timeout > 0 {
			// * Must modify Dial, TLSHandshake and ResponseHeader timeout. because there is a difference between this and http.Client.Timeout (we can't use that)
			proxyConfig.Transport = &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   config.Timeout,
					KeepAlive: config.Timeout,
				}).Dial,
				TLSHandshakeTimeout:   config.Timeout,
				ResponseHeaderTimeout: config.Timeout,
			}
		}

		e.Group(config.Prefix, proxyProcessingFunc(config), middleware.ProxyWithConfig(proxyConfig))
		log.Error(ctx, fmt.Sprintf("[PROXY] Created '%s' -> '%s'", config.Prefix, config.URL), config.Rewrite)
	}
	return e
}

func proxyProcessingFunc(config ProxyMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ctx := log.SetProxyConfigFromEchoContext(c, config.URL, config.Rewrite)
			log.Info(ctx, fmt.Sprintf("[PROXY] Rewriting path from '%s' to '%s'", c.Request().URL.String(), config.URL), config.Rewrite)
			return next(c)
		}
	}
}
