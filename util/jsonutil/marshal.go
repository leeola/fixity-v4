package jsonutil

import (
	"encoding/json"
	"net/http"
)

func MarshalToWriter(w http.ResponseWriter, v interface{}) (int, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(b)
	if err != nil {
		return 0, err
	}

	return n, nil
}
