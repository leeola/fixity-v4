package value

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

func Int(i int) Value {
	return Value{
		Type:     TypeInt,
		IntValue: i,
	}
}
