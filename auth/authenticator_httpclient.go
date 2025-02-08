package auth

import (
	"context"
	"fmt"
	coreCtx "github.com/harryosmar/go-echo-core/context"
	libraryhttpclientgo "github.com/harryosmar/http-client-go"
	v2 "github.com/harryosmar/http-client-go/v2"
)

type (
	AuthenticatorHttpClient struct {
		httpClient                libraryhttpclientgo.HttpClientRepository
		authServiceBaseUrl        string
		tokenValidateEndpointPath string
	}

	ApiAuthCheckRequest struct {
		Token string `json:"token"`
	}

	ApiAuthCheckResponse struct {
		Sub string `json:"sub"`
	}
)

func NewAuthenticatorHttpClient(
	httpClient libraryhttpclientgo.HttpClientRepository,
	authServiceBaseUrl string,
	tokenValidateEndpointPath string,
) *AuthenticatorHttpClient {
	return &AuthenticatorHttpClient{
		httpClient:                httpClient,
		authServiceBaseUrl:        authServiceBaseUrl,
		tokenValidateEndpointPath: tokenValidateEndpointPath,
	}
}

func (a AuthenticatorHttpClient) Check(ctx context.Context, token string) (*coreCtx.JwtClaim, error) {
	response, err := v2.Post[*ApiAuthCheckRequest, ApiAuthCheckResponse](
		context.WithValue(ctx, libraryhttpclientgo.XRequestIdContext, coreCtx.NewContextBuilder(ctx).GetRequestId()),
		a.httpClient,
		fmt.Sprintf("%s/%s", a.authServiceBaseUrl, a.tokenValidateEndpointPath),
		&ApiAuthCheckRequest{Token: token},
		map[string]string{
			"Content-Type": "application/json",
		},
	)
	if err != nil {
		return nil, err
	}

	// @todo implement ACL privileges
	// use variable "response" to set dto.JwtClaim fields value
	return &coreCtx.JwtClaim{
		Iss: "",
		Sub: response.Content.Sub,
		Exp: 0,
		Iat: 0,
		Jti: "",
	}, err
}
