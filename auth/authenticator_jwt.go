package auth

import (
	"context"
	coreContext "github.com/harryosmar/go-echo-core/context"
	coreError "github.com/harryosmar/go-echo-core/error"
	signaturego "github.com/harryosmar/hash-go"
)

type AuthenticatorJwt struct {
	signer signaturego.JwtSign
}

func NewAuthenticatorJwt(signer signaturego.JwtSign) *AuthenticatorJwt {
	return &AuthenticatorJwt{signer: signer}
}

func (a AuthenticatorJwt) Check(ctx context.Context, token string) (*coreContext.JwtClaim, error) {
	mapClaims, err := a.signer.Validate(ctx, token)
	if err != nil {
		return nil, coreError.ErrUnauthorizedAccess.WithError(err)
	}

	sub, found := mapClaims["sub"]
	if !found {
		return nil, coreError.ErrUnauthorizedAccessSub404
	}

	jti, found := mapClaims["jti"]
	if !found {
		return nil, coreError.ErrUnauthorizedAccessJti404
	}

	privilegesInterface, found := mapClaims["privileges"]
	if !found {
		return nil, coreError.ErrUnauthorizedAccessPrivileges404
	}
	privilegesListInterface, ok := privilegesInterface.([]interface{})
	if !ok {
		return nil, coreError.ErrUnauthorizedAccessPrivilegesInvalidFormat
	}

	privileges := []string{}
	for _, pi := range privilegesListInterface {
		privileges = append(privileges, pi.(string))
	}

	// @todo implement ACL privileges
	return &coreContext.JwtClaim{
		Sub:        sub.(string),
		Jti:        jti.(string),
		Privileges: privileges,
	}, nil
}
