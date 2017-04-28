package kalautil

import (
	"encoding/json"

	"github.com/leeola/kala"
)

// MarshalJson marshals to a kala.Json from the given interface.
func MarshalJson(v interface{}) (kala.Json, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return kala.Json{}, err
	}

	return kala.Json{
		Json: json.RawMessage(b),
	}, nil
}
