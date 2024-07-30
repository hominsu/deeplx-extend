package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

const ContextKeyRemoteAddr ContextKey = "x-md-global-remote-addr"

func Info() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(*http.Transport); ok {
					ctx = metadata.AppendToClientContext(ctx, string(ContextKeyRemoteAddr), ht.Request().RemoteAddr)
				}
			}
			return handler(ctx, req)
		}
	}
}
