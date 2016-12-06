package jsonutil

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, s string, code int) {
	w.WriteHeader(code)
	if s == "" {
		return
	}
	b, _ := json.Marshal(struct {
		Error string
	}{
		Error: s,
	})
	w.Write(b)
}
