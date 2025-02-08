package middleware

import (
	error2 "github.com/harryosmar/go-echo-core/error"
	"github.com/labstack/echo/v4"
)

const (
	HeaderInternalToken = "X-Internal-Token"
)

func InternalTokenMiddleware(expectedInternalToken string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) (err error) {
			req := e.Request()
			actualInternalToken := req.Header.Get(HeaderInternalToken)
			if actualInternalToken == "" {
				return error2.ErrUnauthorizedAccess
			}
			if actualInternalToken != expectedInternalToken {
				return error2.ErrUnauthorizedAccess
			}

			return next(e)
		}
	}
}
