package q

import "github.com/leeola/fixity/value"

type Constraint struct {
	Operator       string       `json:"operator"`
	Field          *string      `json:"field,omitempty"`
	Value          *value.Value `json:"value,omitempty"`
	SubConstraints []Constraint `json:"subConstraints,omitempty"`
}

type Query struct {
	LimitBy    int
	Constraint Constraint
}
