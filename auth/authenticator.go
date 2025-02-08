package auth

import (
	"context"
	"github.com/harryosmar/go-echo-core/context"
)

type Authenticator interface {
	Check(ctx context.Context, token string) (*ctx.JwtClaim, error)
}
