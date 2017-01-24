package webware

import (
	"context"
	"net/http"

	"github.com/leeola/kala/client"
)

func ClientMiddleware(c *client.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "client", c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClient(r *http.Request) (*client.Client, bool) {
	ctx := r.Context()
	c, ok := ctx.Value("client").(*client.Client)
	return c, ok
}
