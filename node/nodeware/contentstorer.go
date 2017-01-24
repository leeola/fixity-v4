package nodeware

import (
	"context"
	"net/http"

	"github.com/leeola/kala/contenttype"
)

func ContentStorersMiddleware(cs map[string]contenttype.ContentType) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "contentStorers", cs)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetContentStorers(r *http.Request) (map[string]contenttype.ContentType, bool) {
	ctx := r.Context()
	css, ok := ctx.Value("contentStorers").(map[string]contenttype.ContentType)
	return css, ok
}
