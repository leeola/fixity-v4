package kalautil

import (
	"encoding/json"

	"github.com/leeola/kala"
)

// ToJson creates a kala.Json from the given interface.
func ToJson(v interface{}) (kala.Json, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return kala.Json{}, err
	}

	return kala.Json{
		Json: json.RawMessage(b),
	}, nil
}
