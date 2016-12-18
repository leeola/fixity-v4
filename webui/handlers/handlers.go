package handlers

import "net/http"

func GetRoot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
