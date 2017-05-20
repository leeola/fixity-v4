package fixity

import (
	"encoding/json"

	"github.com/leeola/errors"
)

// JsonBytesWithFields is a helper to add Json and IndexFields to a MultiJson.
//
// The data structre is a bit verbose for MultiJson and  JsonWithMeta,
// hence this helper.
func (m MultiJson) JsonBytesWithFields(key string, b []byte, f Fields) {
	jwm := JsonWithMeta{
		Json: Json{
			json.RawMessage(b),
		},
	}

	if len(f) > 0 {
		jwm.JsonMeta = &JsonMeta{
			IndexedFields: f,
		}
	}

	m[key] = jwm
}

// MarshalWithFields marshals and adds the interface with fields to MultiJson.
func (m MultiJson) MarshalWithFields(key string, v interface{}, f Fields) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	m[key] = JsonWithMeta{
		Json: Json{
			JsonBytes: json.RawMessage(b),
		},
		JsonMeta: &JsonMeta{
			IndexedFields: f,
		},
	}

	return nil
}

// Unmarshal unmarshals the given key's Json to the given interface.
func (m MultiJson) Unmarshal(key string, v interface{}) error {
	jwm, ok := m[key]
	if !ok {
		return errors.Errorf("multijson key not found: %s", key)
	}

	return json.Unmarshal([]byte(jwm.JsonBytes), v)
}
