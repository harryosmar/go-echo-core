package middleware

import (
	"github.com/google/uuid"
	coreCtx "github.com/harryosmar/go-echo-core/context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

var RequestIdMiddleware = middleware.RequestIDWithConfig(middleware.RequestIDConfig{
	Generator: func() string {
		return uuid.New().String()
	},
	RequestIDHandler: func(c echo.Context, requestId string) {
		request := c.Request()
		ctx := request.Context()
		loggerEntry := log.NewEntry(coreCtx.CustomLogger).WithField("request_id", requestId)

		newCtx := coreCtx.NewContextBuilder(ctx).SetRequestId(requestId).SetLogger(loggerEntry).Context()
		newRequest := request.Clone(newCtx)
		c.SetRequest(newRequest)
	},
})
