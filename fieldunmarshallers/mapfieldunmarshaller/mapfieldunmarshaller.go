package mapfieldunmarshaller

import (
	"encoding/json"

	"github.com/leeola/kala"
)

// TODO(leeola): Make the unmarshaller configurable so it's not just Json.
type MapFieldUnmarshaller struct {
	B []byte
	M map[string]interface{}
}

func New(b []byte) *MapFieldUnmarshaller {
	return &MapFieldUnmarshaller{
		B: b,
	}
}

func (m *MapFieldUnmarshaller) Unmarshal(field string) (interface{}, error) {
	if m.M == nil {
		if err := json.Unmarshal(m.B, &m.M); err != nil {
			return nil, err
		}
		// let GC clear the bytes
		m.B = nil
	}

	v, ok := m.M[field]
	if !ok {
		return nil, kala.ErrFieldNotFound
	}

	return v, nil
}
