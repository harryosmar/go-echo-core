package middleware

import (
	genericgorm "github.com/harryosmar/generic-gorm"
	"github.com/harryosmar/go-echo-core/auth"
	coreCtx "github.com/harryosmar/go-echo-core/context"
	coreError "github.com/harryosmar/go-echo-core/error"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"strings"
)

func AuthRolesCheckMiddleware(authenticator auth.Authenticator, roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) (err error) {
			req := e.Request()
			ctx := req.Context()
			contextBuilder := coreCtx.NewContextBuilder(ctx)
			logger := contextBuilder.GetLogger()

			authorizationStr := req.Header.Get(echo.HeaderAuthorization)
			tokenParts := strings.Split(authorizationStr, " ")
			if len(tokenParts) != 2 {
				return coreError.ErrUnauthorizedAccess
			}
			if tokenParts[0] != "Bearer" {
				return coreError.ErrUnauthorizedAccess
			}
			token := tokenParts[1]

			checkResult, err := authenticator.Check(ctx, token)
			if err != nil {
				return err
			}
			// set user session context
			logger = logger.WithFields(log.Fields{
				"session_id": checkResult.Jti,
				"user_id":    checkResult.Username,
			})

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
			newCtx := coreCtx.NewContextBuilder(ctx).
				SetSession(coreCtx.NewSession(checkResult)).
				SetLogger(logger).
				Context()

			// gorm context logger
			newCtx = genericgorm.ContextWithLogger(newCtx, logger)

			// set new context to request
			newRequest := req.Clone(newCtx)
			e.SetRequest(newRequest)

			return next(e)
		}
	}
}
