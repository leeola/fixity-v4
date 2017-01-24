package webware

import (
	"context"
	"net/http"
)

func TemplatersMiddleware(t map[string]interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "contentTemplaters", t)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetContentTemplater(r *http.Request, ctype string) (interface{}, bool) {
	ctx := r.Context()
	cMap, ok := ctx.Value("contentTemplaters").(map[string]interface{})
	if !ok {
		return nil, false
	}

	c, ok := cMap[ctype]
	if !ok {
		return nil, false
	}

	return c, true
}
