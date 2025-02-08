package middleware

import (
	coreCtx "github.com/harryosmar/go-echo-core/context"
	"github.com/labstack/echo/v4"
)

func PrivilegeCheckMiddleware(expectedPrivileges []string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) (err error) {
			ctx := e.Request().Context()
			session := coreCtx.NewContextBuilder(ctx).GetSession()
			err = session.IsHasPrivileges(expectedPrivileges)
			if err != nil {
				return err
			}

			return next(e)
		}
	}
}
