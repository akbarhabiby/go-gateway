package log

import (
	"context"
	"time"

	"github.com/akbarhabiby/go-gateway/constants"
	"github.com/labstack/echo/v4"
)

const (
	CTX_REQUEST_IP       constants.String = "ctx-request-ip"
	CTX_REQUEST_TIME     constants.String = "ctx-request-time"
	CTX_REQUEST_HEADER   constants.String = "ctx-request-header"
	CTX_ERROR_MESSAGE    constants.String = "ctx-error-message"
	CTX_PROXY_PROCESSING constants.String = "ctx-proxy-processing"
	CTX_PROXY_URL        constants.String = "ctx-proxy-url"
	CTX_PROXY_TARGET_URL constants.String = "ctx-proxy-target-url"
)

type ctxKeyLogger struct{}

var ctxKey = ctxKeyLogger{}

type LogContextModel struct {
	ServiceName    string `json:"_app_name"`
	ServiceVersion string `json:"_app_version"`
	ServicePort    int    `json:"_app_port"`
	ThreadID       string `json:"_app_thread_id"`

	ReqMethod string `json:"_app_method"`
	ReqURI    string `json:"_app_uri"`
}

func SetContextFromEchoRequest(c echo.Context) context.Context {
	ctx := c.Request().Context()

	ctx = context.WithValue(ctx, CTX_REQUEST_IP, c.RealIP())
	ctx = context.WithValue(ctx, CTX_REQUEST_TIME, time.Now())
	ctx = context.WithValue(ctx, CTX_REQUEST_HEADER, c.Request().Header)

	return ctx
}

func SetErrorMessageFromEchoContext(c echo.Context, errMessage string) context.Context {
	ctx := c.Request().Context()

	ctx = context.WithValue(ctx, CTX_ERROR_MESSAGE, errMessage)

	return ctx
}

func SetProxyConfigFromEchoContext(c echo.Context, targetURL string, rewrite map[string]string) context.Context {
	ctx := c.Request().Context()

	c.Set(string(CTX_PROXY_PROCESSING), rewrite)
	c.Set(string(CTX_PROXY_URL), c.Request().URL.String())
	c.Set(string(CTX_PROXY_TARGET_URL), targetURL)

	return ctx
}

func GetRequestIPFromContext(ctx context.Context) string {
	s, ok := ctx.Value(CTX_REQUEST_IP).(string)
	if !ok {
		return ""
	}

	return s
}

func GetRequestTimeFromContext(ctx context.Context) time.Time {
	s, ok := ctx.Value(CTX_REQUEST_TIME).(time.Time)
	if !ok {
		return time.Now()
	}

	return s
}

func GetRequestHeaderFromContext(ctx context.Context) interface{} {
	return ctx.Value(CTX_REQUEST_HEADER)
}

func GetProxyProcessingFromEchoContext(c echo.Context) interface{} {
	return c.Get(string(CTX_PROXY_PROCESSING))
}

func GetProxyURLFromEchoContext(c echo.Context) interface{} {
	return c.Get(string(CTX_PROXY_URL))
}

func GetProxyTargetURLFromEchoContext(c echo.Context) interface{} {
	return c.Get(string(CTX_PROXY_TARGET_URL))
}

func GetErrorMessageFromContext(ctx context.Context) string {
	s, ok := ctx.Value(CTX_ERROR_MESSAGE).(string)
	if !ok {
		return ""
	}

	return s
}

func InjectCtx(parent context.Context, ctx LogContextModel) context.Context {
	if parent == nil {
		return InjectCtx(context.Background(), ctx)
	}

	return context.WithValue(parent, ctxKey, ctx)
}

func ExtractCtx(ctx context.Context) LogContextModel {
	if ctx == nil {
		return LogContextModel{}
	}

	val, ok := ctx.Value(ctxKey).(LogContextModel)
	if !ok {
		return LogContextModel{}
	}

	return val
}
