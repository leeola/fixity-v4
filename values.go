package fixity

type Values map[string]Value

type Value struct {
	Type        ValueType
	IntValue    int
	StringValue string
}

type ValueType int

const (
	ValueTypeInt    ValueType = 1
	ValueTypeString ValueType = 2
)

func (m Values) SetInt(key string, v int) {
	m[key] = Value{
		Type:     ValueTypeInt,
		IntValue: v,
	}
}

func (m Values) Int(key string) (int, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}

	if v.Type != ValueTypeInt {
		return 0, false
	}

	return v.IntValue, true
}
