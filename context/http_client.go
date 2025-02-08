package ctx

import (
	"context"
	libraryhttpclientgo "github.com/harryosmar/http-client-go"
)

func NewHttpClientContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, libraryhttpclientgo.XRequestIdContext, NewContextBuilder(ctx).GetRequestId())
}
