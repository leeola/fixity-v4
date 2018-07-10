package q

import (
	"github.com/leeola/fixity/q/operator"
	"github.com/leeola/fixity/value"
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

func (q Query) Const(c Constraint) Query {
	q.Constraint = c
	return q
}

func (q Query) Eq(field string, value value.Value) Query {
	return q.Const(Eq(field, value))
}

func Eq(field string, value value.Value) Constraint {
	return Constraint{
		Operator: operator.Equal,
		Field:    &field,
		Value:    &value,
	}
}

func (q Query) And(c ...Constraint) Query {
	q.Const(And(c...))
	return q
}

// And requires that all given constraints are succeed.
//
// Note that if a single constraint is supplied, no AND constraint is
// returned, only the single constraint. In other words, with a
// constraint like AND(Eq(1,2)), AND is completely pointless.
func And(c ...Constraint) Constraint {
	if len(c) == 1 {
		return c[0]
	}

	return Constraint{
		Operator:       operator.And,
		SubConstraints: c,
	}
}
