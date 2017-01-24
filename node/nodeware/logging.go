package nodeware

import (
	"context"
	"net/http"

	"github.com/inconshreveable/log15"
)

func LoggingMiddleware(name string, log log15.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			rLog := log.New("name", name)
			ctx = context.WithValue(ctx, "log", rLog)
			rLog.Debug("http request", "method", r.Method, "url", r.URL)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLog(r *http.Request) log15.Logger {
	ctx := r.Context()
	if log, ok := ctx.Value("log").(log15.Logger); ok {
		return log
	}

	log := log15.New()
	log.Error("failed to get log from request context")
	return log
}
