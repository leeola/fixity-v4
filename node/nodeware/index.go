package nodeware

import (
	"context"
	"net/http"

	"github.com/leeola/kala/index"
)

func IndexMiddleware(q index.Indexer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "index", q)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetIndex(r *http.Request) (index.Indexer, bool) {
	ctx := r.Context()
	v, ok := ctx.Value("index").(index.Indexer)
	return v, ok
}
