package fixity

import "github.com/leeola/fixity/value"

type Values map[string]value.Value

func (m Values) Int(key string) (int, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}

	if v.Type != value.TypeInt {
		return 0, false
	}

	return v.IntValue, true
}
