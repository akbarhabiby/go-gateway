package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akbarhabiby/go-gateway/config"
	"github.com/akbarhabiby/go-gateway/constants"
	"github.com/akbarhabiby/go-gateway/log"
	"github.com/akbarhabiby/go-gateway/middlewares"
	"github.com/akbarhabiby/go-gateway/proxy"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const banner = `

                                 .__    .___  .__.__                                     
   ____   ____             _____ |__| __| _/__| _|  |   ______  _  ______ _______  ____  
  / ___\ /  _ \   ______  /     \|  |/ __ |/ __ ||  | _/ __ \ \/ \/ \__  \\_  __ _/ __ \ 
 / /_/  (  <_> ) /_____/ |  Y Y  |  / /_/ / /_/ ||  |_\  ___/\     / / __ \|  | \\  ___/ 
 \___  / \____/          |__|_|  |__\____ \____ ||____/\___  >\/\_/ (____  |__|   \___  >
/_____/                        \/        \/    \/          \/            \/           \/ 



`

func main() {
	fmt.Print(banner)
	e := echo.New()

	// * Middlewares
	middlewares.Setup(e)

	e.GET("/", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, constants.DefaultResponse{
			Status:  constants.STATUS_SUCCESS,
			Message: constants.MESSAGE_SUCCESS,
			Data: constants.HealthCheckResponse{
				Message:    "Server up and running",
				ServerTime: time.Now().Format(time.RFC1123),
				Version:    "v1.0.0",
			},
			Errors: make([]string, 0),
		})
	})

	proxyConfigs := []proxy.ProxyMiddlewareConfig{}
	proxy.Setup(e, proxyConfigs)

	address := config.Config.GetString("app.address")
	port := config.Config.GetString("app.port")

	e.Server.Addr = fmt.Sprintf("%s:%s", address, port)

	log.Info(context.Background(), fmt.Sprintf("[SERVER] h2c server started on %s:%s", address, port))

	// * HTTP/2 Cleartext Server (HTTP2 over HTTP)
	gracehttp.Serve(&http.Server{Addr: e.Server.Addr, Handler: h2c.NewHandler(e, &http2.Server{MaxConcurrentStreams: 500, MaxReadFrameSize: 1048576})})
}
