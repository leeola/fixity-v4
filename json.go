package fixity

import "encoding/json"

// AddJsonWithFields is a helper to add Json and IndexFields to a MultiJson.
//
// The data structre is a bit verbose for MultiJson and  JsonWithMeta,
// hence this helper.
func (m MultiJson) AddJsonWithFields(key string, b []byte, f Fields) {
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
