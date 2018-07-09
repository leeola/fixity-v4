package value

import (
	"fmt"
	"strconv"
)

//go:generate stringer -type=Type -output=value_string.go

type Value struct {
	Type        Type   `json:"type"`
	IntValue    int    `json:"intValue,omitempty"`
	StringValue string `json:"stringValue,omitempty"`
}

type Type int

const (
	TypeInt    Type = 1
	TypeString Type = 2
)

func Int(v int) Value {
	return Value{
		Type:     TypeInt,
		IntValue: v,
	}
}

func String(v string) Value {
	return Value{
		Type:        TypeString,
		StringValue: v,
	}
}

// Value returns an untyped value of whatever value field is defined
// by Value.Type.
//
// This should not be used unless the type is already being thrown away,
// such as in a map[string]interface{} or etc.
func (v Value) UntypedValue() (interface{}, error) {
	switch v.Type {
	case TypeInt:
		return v.IntValue, nil
	case TypeString:
		return v.StringValue, nil
	default:
		return nil, fmt.Errorf("unexpected value type: %s", v.Type)
	}
}

// ToString returns a string representation of the Value struct's typed value.
func (v Value) ToString() (string, error) {
	switch v.Type {
	case TypeInt:
		return strconv.Itoa(v.IntValue), nil
	case TypeString:
		return v.StringValue, nil
	default:
		return "", fmt.Errorf("unexpected value type: %s", v.Type)
	}
}

func (v Value) String() string {
	switch v.Type {
	case TypeInt:
		return fmt.Sprintf("IntValue(%d)", v.IntValue)
	case TypeString:
		return fmt.Sprintf("StringValue(%s)", v.StringValue)
	default:
		return "UnknownValue"
	}
}
