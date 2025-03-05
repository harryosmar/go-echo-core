package auth

import (
	"context"
	"github.com/harryosmar/go-echo-core/context"
	"time"
)

type Authenticator interface {
	Check(ctx context.Context, token string) (*ctx.JwtClaim, error)
	Persist(ctx context.Context, claim *ctx.JwtClaim, now time.Time) error
}
