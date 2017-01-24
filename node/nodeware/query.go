package nodeware

import (
	"context"
	"net/http"

	"github.com/leeola/kala/index"
)

func QueryMiddleware(q index.Queryer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "query", q)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetQuery(r *http.Request) (index.Queryer, bool) {
	ctx := r.Context()
	v, ok := ctx.Value("query").(index.Queryer)
	return v, ok
}
