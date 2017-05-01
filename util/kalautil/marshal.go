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

// MustMarshalJson panics if the Marshal fails.
func MustMarshalJson(v interface{}) kala.Json {
	j, err := MarshalJson(v)
	if err != nil {
		panic(err.Error())
	}
	return j
}

// UnmarshalJson unmarshals the given Json struct into the given interface.
func UnmarshalJson(j kala.Json, v interface{}) error {
	return json.Unmarshal([]byte(j.Json), v)
}
