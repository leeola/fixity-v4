package fixityutil

import (
	"encoding/json"

	"github.com/leeola/fixity"
)

// MarshalJson marshals to a fixity.Json from the given interface.
func MarshalJson(v interface{}) (fixity.Json, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return fixity.Json{}, err
	}

	return fixity.Json{
		Json: json.RawMessage(b),
	}, nil
}

// MustMarshalJson panics if the Marshal fails.
func MustMarshalJson(v interface{}) fixity.Json {
	j, err := MarshalJson(v)
	if err != nil {
		panic(err.Error())
	}
	return j
}

// UnmarshalJson unmarshals the given Json struct into the given interface.
func UnmarshalJson(j fixity.Json, v interface{}) error {
	return json.Unmarshal([]byte(j.Json), v)
}
