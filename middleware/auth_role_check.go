package middleware

import (
	"github.com/harryosmar/go-echo-core/auth"
	coreCtx "github.com/harryosmar/go-echo-core/context"
	coreError "github.com/harryosmar/go-echo-core/error"
	"github.com/labstack/echo/v4"
	"strings"
)

func AuthRolesCheckMiddleware(authenticator auth.Authenticator, roles []string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) (err error) {
			req := e.Request()
			authorizationStr := req.Header.Get(echo.HeaderAuthorization)
			tokenParts := strings.Split(authorizationStr, " ")
			if len(tokenParts) != 2 {
				return coreError.ErrUnauthorizedAccess
			}
			if tokenParts[0] != "Bearer" {
				return coreError.ErrUnauthorizedAccess
			}
			token := tokenParts[1]

			ctx := req.Context()
			checkResult, err := authenticator.Check(ctx, token)
			if err != nil {
				return err
			}

			isRoleValid := false
			for _, role := range roles {
				if checkResult.Role.Code == role {
					isRoleValid = true
					break
				}
			}
			if !isRoleValid {
				return coreError.ErrForbiddenInvalidRole
			}

			// set user session context
			newCtx := coreCtx.NewContextBuilder(ctx).SetSession(coreCtx.NewSession(checkResult)).Context()
			newRequest := req.Clone(newCtx)
			e.SetRequest(newRequest)

			return next(e)
		}
	}
}
