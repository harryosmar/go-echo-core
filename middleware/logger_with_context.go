package middleware

import (
	"fmt"
	"github.com/harryosmar/go-echo-core/context"
	"github.com/harryosmar/go-echo-core/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

var (
	LoggerWithContextMiddleware = []echo.MiddlewareFunc{
		middleware.RequestLoggerWithConfig(
			middleware.RequestLoggerConfig{
				LogURI:    true,
				LogStatus: true,
				BeforeNextFunc: func(c echo.Context) {
					request := c.Request()
					loggerEntry := ctx.NewContextBuilder(request.Context()).GetLogger()

					//payload, err := httputil.DumpRequest(request, true)
					payload, err := utils.GetCopyPayloadFromRequest(request)
					payloadStr := ""
					if err == nil {
						payloadStr = string(payload)
					}

					loggerEntry.WithFields(log.Fields{
						"URI":    request.RequestURI,
						"method": request.Method,
						"headers": func() string {
							if request == nil {
								return ""
							}
							if request.Header == nil {
								return ""
							}
							return fmt.Sprintf("%+v", request.Header.Clone())
						}(),
						"payload": payloadStr,
					}).Info("request")
				},
				LogError: true,
				LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
					request := c.Request()
					loggerEntry := ctx.NewContextBuilder(request.Context()).GetLogger()
					if values.Error != nil {
						loggerEntry.Error(values.Error)
						return nil
					}
					loggerEntry.Info("OK")
					return nil
				},
			},
		),
		//middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		// disable create duplicated content
		//	request := c.Request()
		//	loggerEntry := ctx.NewContextBuilder(request.Context()).GetLogger()
		//	loggerEntry.WithFields(log.Fields{
		//		"content": string(resBody),
		//	}).Info("res_body")
		//}),
	}
)
