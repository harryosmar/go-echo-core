package middleware

import (
	coreCtx "github.com/harryosmar/go-echo-core/context"
	error2 "github.com/harryosmar/go-echo-core/error"
	"github.com/labstack/echo/v4"
)

func PrivilegeCheckMiddleware(expectedPrivileges []string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) (err error) {
			ctx := e.Request().Context()
			session := coreCtx.NewContextBuilder(ctx).GetSession()
			if session == nil {
				return error2.ErrUnauthorizedAccess
			}
			err = session.IsHasPrivileges(expectedPrivileges)
			if err != nil {
				return err
			}

			return next(e)
		}
	}
}
