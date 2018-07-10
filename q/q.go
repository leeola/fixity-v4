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
	IncludeVersions bool
	LimitBy         int
	Constraint      Constraint
}

func New() Query {
	return Query{
		// set a default limit
		LimitBy: 10,
	}
}

func (q Query) WithVersions() Query {
	q.IncludeVersions = true
	return q
}

func (q Query) WithoutVersions() Query {
	q.IncludeVersions = false
	return q
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
