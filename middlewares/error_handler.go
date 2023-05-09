package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/akbarhabiby/go-gateway/constants"
	"github.com/akbarhabiby/go-gateway/log"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Get("error-handled") != nil {
		return
	}

	c.Set("error-handled", true)

	res := constants.DefaultResponse{
		Status:  constants.STATUS_FAILED,
		Message: err.Error(),
		Data:    struct{}{},
		Errors:  make([]string, 0),
	}

	request := c.Request()
	ctx := log.SetErrorMessageFromEchoContext(c, err.Error())

	if c.Get("request-limit-reached") != nil {
		res.Status = fmt.Sprintf("0%s", strconv.FormatInt(http.StatusTooManyRequests, 10))
		res.Message = http.StatusText(http.StatusTooManyRequests)
		res.Errors = append(res.Errors, err.Error())
	} else if strings.Contains(err.Error(), "code") && strings.Contains(err.Error(), "message") {
		// * echo error
		eMessage := strings.SplitAfter(err.Error(), "message=")
		if len(eMessage) > 1 {
			if strings.Contains(eMessage[1], "remote") && strings.Contains(eMessage[1], "forward") {
				res.Status = fmt.Sprintf("0%s", strconv.FormatInt(http.StatusBadGateway, 10))
				res.Message = http.StatusText(http.StatusBadGateway)
				res.Errors = append(res.Errors, "Please contact administrator.")
			} else {
				res.Message = eMessage[1]
				res.Errors = append(res.Errors, eMessage[1])
				eCode := strings.SplitAfter(eMessage[0], "code=")
				if len(eCode) > 1 {
					res.Status = fmt.Sprintf("0%s", strings.Replace(eCode[1], ", message=", "", 1))
				}
			}
		}
	}

	c.SetRequest(request.WithContext(ctx))

	if isLoggingSkipForProxy(c) {
		rewrite := log.GetProxyProcessingFromEchoContext(c).(map[string]string)
		log.Error(ctx, fmt.Sprintf("[PROXY] Error occurred while proxying request '%s'", c.Request().URL.String()), rewrite, err)
	}

	c.JSON(http.StatusOK, res)
}
