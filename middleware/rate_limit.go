package middleware

import (
	"github.com/labstack/echo/v4"
	"strings"

	"github.com/harryosmar/go-echo-core/auth"
	coreError "github.com/harryosmar/go-echo-core/error"
	"github.com/labstack/echo/v4/middleware"
)

func RateLimitMiddleware(store middleware.RateLimiterStore, authenticator auth.Authenticator) echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store:   store,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			req := ctx.Request()
			authHeader := req.Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				return id, nil
			}
			splittedAuthHeader := strings.Split(authHeader, " ")
			if len(splittedAuthHeader) != 2 {
				return id, nil
			}
			if splittedAuthHeader[0] != "Bearer" {
				return id, nil
			}
			token := splittedAuthHeader[1]
			return token, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return err
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return coreError.ErrTooManyRequests
		},
	})
}
