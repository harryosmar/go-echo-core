package middleware

import (
	"github.com/google/uuid"
	genericgorm "github.com/harryosmar/generic-gorm"
	coreCtx "github.com/harryosmar/go-echo-core/context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
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

func GetClientIP(c echo.Context) string {
	// 1. CF-Connecting-IP (Cloudflare) – trusted source
	if ip := strings.TrimSpace(c.Request().Header.Get("CF-Connecting-IP")); ip != "" {
		return normalizeIP(ip)
	}

	// 2. X-Forwarded-For (ALB may include this)
	// Format: "client, proxy1, proxy2"
	if xff := strings.TrimSpace(c.Request().Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return normalizeIP(strings.TrimSpace(parts[0]))
		}
	}

	// 3. X-Real-IP (optional fallback)
	if ip := strings.TrimSpace(c.Request().Header.Get("X-Real-IP")); ip != "" {
		return normalizeIP(ip)
	}

	// 4. Fallback: RemoteAddr (localhost or direct request)
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request().RemoteAddr))
	if err != nil {
		return normalizeIP(c.Request().RemoteAddr)
	}
	return normalizeIP(ip)
}

func normalizeIP(ip string) string {
	// Handle IPv6 loopback
	if ip == "::1" {
		return "127.0.0.1"
	}
	// Optionally strip zone identifiers from IPv6 (e.g. fe80::1%eth0)
	if strings.Contains(ip, "%") {
		ip = strings.Split(ip, "%")[0]
	}
	return ip
}

var RequestIdMiddlewareV2 = func(appEnv string, appVersion string) echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
		RequestIDHandler: func(c echo.Context, requestId string) {
			request := c.Request()
			ctx := request.Context()
			contextBuilder := coreCtx.NewContextBuilder(ctx)

			loggerEntry := contextBuilder.GetLogger().WithFields(log.Fields{
				"environment": appEnv,
				"version":     appVersion,
				"request_id":  requestId,
				"ip_address":  GetClientIP(c),
			})

			// session context logger
			newCtx := contextBuilder.
				SetRequestId(requestId).
				SetLogger(loggerEntry).
				Context()

			// gorm context logger
			newCtx = genericgorm.ContextWithLogger(newCtx, loggerEntry)

			// set new context to request
			newRequest := request.Clone(newCtx)
			c.SetRequest(newRequest)
		},
	})
}
