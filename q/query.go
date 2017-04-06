package q

import "github.com/leeola/kala/q/operators"

type Constraint struct {
	// Operator
	Operator    string
	Field       string
	Value       interface{}
	Constraints []Constraint
}

// Query is a generic method to construct a query in the index implementation.
//
// Eg, if mysql is the index, a Query would be constructed into a SQL string.
// Furthermore, the generic constraints allow for any type of keywords/etc
// to be used by the underlying indexer implementation. FullTextSearch for
// example is a more niche feature, and not supported by many indexers.
type Query struct {
	SortField      string
	SortDescending string
	SkipBy         int
	LimitBy        int
	Constraint     Constraint
}

func New() *Query {
    return &Query{}
}

func (q *Query) Limit(l int) *Query {
	q.LimitBy = l
	return q
}

func (q *Query) Skip(s int) *Query {
	q.SkipBy = s
	return q
}

func (q *Query) Sort(field string, descending bool) *Query {
	q.SortField = field
	q.SortDescending = descending
	return q
}

func (q *Query) And(c ...Constraint) *Query {
	q.Const(And(c...))
	return q
}

func (q *Query) Const(c Constraint) *Query {
	q.Constraint = c
	return q
}

func (q *Query) Or(c ...Constraint) *Query {
	q.Const(Or(c...))
	return q
}

func And(c ..Constraint) Constraint {
	return Constraint{
		Operator: operators.Equal,
		Constraints: c,
	}
}

func Eq(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.Equal,
		Field:    field,
		Value:    value,
	}
}

func Gt(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.GreaterThan,
		Field:    field,
		Value:    value,
	}
}

func Gte(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.GreaterThanOrEqual,
		Field:    field,
		Value:    value,
	}
}

func In(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.In,
		Field:    field,
		Value:    value,
	}
}

func Lt(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.GreaterThan,
		Field:    field,
		Value:    value,
	}
}

func Lte(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.GreaterThan,
		Field:    field,
		Value:    value,
	}
}

func Not(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.Not,
		Field:    field,
		Value:    value,
	}
}

func Or(c ...Constraint) Constraint {
	return Constraint{
		Operator: operators.Or,
		Constraints: c,
	}
}

func Re(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.Regex,
		Field:    field,
		Value:    value,
	}
}
