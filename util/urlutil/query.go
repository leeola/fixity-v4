package urlutil

import (
	"net/http"
	"strconv"
)

// GetQueryString exists just for completeness sake.
func GetQueryString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func GetQueryBool(r *http.Request, key string) bool {
	s := r.URL.Query().Get(key)
	v, _ := strconv.ParseBool(s)
	return v
}
