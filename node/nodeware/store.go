package nodeware

import (
	"context"
	"net/http"

	"github.com/leeola/kala/store"
)

func StoreMiddleware(s store.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "store", s)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetStore(r *http.Request) (store.Store, bool) {
	ctx := r.Context()
	s, ok := ctx.Value("store").(store.Store)
	return s, ok
}
