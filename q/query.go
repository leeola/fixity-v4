package q

import "github.com/leeola/fixity/q/operators"

type Constraint struct {
	// Operator
	Operator    string
	Field       string
	Value       interface{}
	Constraints []Constraint
}

type Constraints []Constraint

type Sort struct {
	Field      string
	Descending bool
}

// Query is a generic method to construct a query in the index implementation.
//
// Eg, if mysql is the index, a Query would be constructed into a SQL string.
// Furthermore, the generic constraints allow for any type of keywords/etc
// to be used by the underlying indexer implementation. FullTextSearch for
// example is a more niche feature, and not supported by many indexers.
type Query struct {
	SortBy     []Sort
	SkipBy     int
	LimitBy    int
	Constraint Constraint
}

func New() *Query {
	return &Query{
		// set a default limit
		LimitBy: 10,
	}
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
	q.SortBy = append(q.SortBy, Sort{
		Field:      field,
		Descending: descending,
	})
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

// And requires that all given constraints are succeed.
//
// Note that if a single constraint is used, no And constraint is used.
// This is mainly to reduce needless abstractions, albeit not intuitively.
func And(c ...Constraint) Constraint {
	if len(c) == 1 {
		return c[0]
	}

	return Constraint{
		Operator:    operators.And,
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

func Fts(field string, value interface{}) Constraint {
	return Constraint{
		Operator: operators.FullTextSearch,
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

// In checks if the given value is found within a list of values.
//
// Details are up to the implementor, but usually a equality check is done
// on the values within the list.
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

func Not(c ...Constraint) Constraint {
	return Constraint{
		Operator:    operators.Not,
		Constraints: c,
	}
}

// Or requires that any of the given constraints succeed.
//
// Note that if a single constraint is used, no Or constraint is used.
// This is mainly to reduce needless abstractions, albeit not intuitively.
func Or(c ...Constraint) Constraint {
	if len(c) == 1 {
		return c[0]
	}

	return Constraint{
		Operator:    operators.Or,
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

// Add the given constraint to this slice.
func (cs *Constraints) Add(c Constraint) {
	*cs = append(*cs, c)
}

// Eq creates and adds an Eq constraint to this constraint slice.
func (cs *Constraints) Eq(field string, value interface{}) {
	cs.Add(Eq(field, value))
}

// Fts creates and adds an Fts constraint to this constraint slice.
func (cs *Constraints) Fts(field string, value interface{}) {
	cs.Add(Fts(field, value))
}

// In creates and adds an In constraint to this constraint slice.
func (cs *Constraints) In(field string, value interface{}) {
	cs.Add(In(field, value))
}
