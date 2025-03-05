package auth

import (
	"context"
	"fmt"
	cache "github.com/harryosmar/cache-go"
	coreContext "github.com/harryosmar/go-echo-core/context"
	coreError "github.com/harryosmar/go-echo-core/error"
	"time"
)

type AuthenticatorJwtSingleSession struct {
	authenticator *AuthenticatorJwt
	cache         cache.CacheRepo
}

func NewAuthenticatorJwtSingleSession(authenticator *AuthenticatorJwt, cache cache.CacheRepo) *AuthenticatorJwtSingleSession {
	return &AuthenticatorJwtSingleSession{authenticator: authenticator, cache: cache}
}

func (a AuthenticatorJwtSingleSession) Persist(ctx context.Context, claim *coreContext.JwtClaim, now time.Time) error {
	return a.cache.Store(
		ctx,
		a.generateKey(claim.Sub),
		[]byte(claim.Jti),
		time.Duration(claim.Exp-now.Unix())*time.Second,
	)
}

func (a AuthenticatorJwtSingleSession) isValidSession(ctx context.Context, claim *coreContext.JwtClaim) error {
	cacheKey := a.generateKey(claim.Sub)
	bytes, found, err := a.cache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}
	if !found {
		return coreError.ErrForbiddenSessionNotFound
	}
	actualJti := string(bytes)
	if actualJti != claim.Jti {
		return coreError.ErrForbiddenMultipleSessionDetected
	}
	return nil
}

func (a AuthenticatorJwtSingleSession) generateKey(sub string) string {
	return fmt.Sprintf("login:session:%s", sub)
}

func (a AuthenticatorJwtSingleSession) Check(ctx context.Context, token string) (*coreContext.JwtClaim, error) {
	jwtClaim, err := a.authenticator.Check(ctx, token)
	if err != nil {
		return nil, err
	}

	err = a.isValidSession(ctx, jwtClaim)
	if err != nil {
		return nil, err
	}

	return jwtClaim, err
}
