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
		JsonBytes: json.RawMessage(b),
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
	return json.Unmarshal([]byte(j.JsonBytes), v)
}

// MarshalJsonWithMeta marshals to a fixity.JsonWithMeta from the given interface.
func MarshalJsonWithMeta(v interface{}) (fixity.JsonWithMeta, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return fixity.JsonWithMeta{}, err
	}

	return fixity.JsonWithMeta{
		Json: fixity.Json{
			JsonBytes: json.RawMessage(b),
		},
	}, nil
}

// MustMarshalJsonWithMeta panics if the Marshal fails.
func MustMarshalJsonWithMeta(v interface{}) fixity.JsonWithMeta {
	j, err := MarshalJsonWithMeta(v)
	if err != nil {
		panic(err.Error())
	}
	return j
}
