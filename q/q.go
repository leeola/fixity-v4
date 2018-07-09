package q

import (
	"github.com/leeola/fixity/value"
	"github.com/leeola/fixity/value/operator"
)

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

func New() Query {
	return Query{
		// set a default limit
		LimitBy: 10,
	}
}

func (q Query) Eq(field string, value value.Value) Query {
	q.Constraint = Eq(field, value)
	return q
}

func Eq(field string, value value.Value) Constraint {
	return Constraint{
		Operator: operator.Equal,
		Field:    &field,
		Value:    &value,
	}
}
