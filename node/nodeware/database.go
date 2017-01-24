package nodeware

import (
	"context"
	"net/http"

	"github.com/leeola/kala/database"
)

func DatabaseMiddleware(db database.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "database", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDatabase(r *http.Request) (database.Database, bool) {
	ctx := r.Context()
	db, ok := ctx.Value("database").(database.Database)
	return db, ok
}
