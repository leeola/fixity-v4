package storm

import (
	"github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	"github.com/leeola/errors"
	kq "github.com/leeola/kala/q"
	"github.com/leeola/kala/q/operators"
)

func ConvertQuery(db storm.DB, kq kq.Query) (storm.Query, error) {
	var ms []sq.Matcher

	if !structs.IsZero(kq.Constraint) {
		m, err := ConvertConstraint(kq.Constraint)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	// TODO(leeola): Storm may not accept an empty set of constraints, so we may need to
	// reject kala queries without constraints.
	stormQuery := storm.Select(ms...)

	if kq.SortBy != "" {
		stormQuery.OrderBy(kq.SortBy)
	}
	if kq.SortDescending {
		stormQuery.Reverse()
	}
	if kq.Skip != 0 {
		stormQuery.Skip(kq.Skip)
	}
	if kq.Limit != 0 {
		stormQuery.Limit(kq.Limit)
	}

	return stormQuery, nil
}

func ConvertConstraint(c kq.Constraint) (sq.Matcher, error) {
	switch c.Operator {
	case operators.Equal:
		return sq.Eq(c.Field, c.Value)

	case operators.GreaterThan:
		return sq.Gt(c.Field, c.Value)

	case operators.GreaterThanOrEqual:
		return sq.Gte(c.Field, c.Value)

	case operators.In:
		return sq.In(c.Field, c.Value)

	case operators.LessThan:
		return sq.Lt(c.Field, c.Value)

	case operators.LessThanOrEqual:
		return sq.Lte(c.Field, c.Value)

	case operators.Not:
		return sq.Not(c.Field, c.Value)

	case operators.Regex:
		return sq.Re(c.Field, c.Value)

	case operators.And, operators.Or:
		matchers := make([]sq.Matcher, len(c.Constraints))
		for i, subc := range c.Constraints {
			m, err := ConvertConstraint(subc)
			if err != nil {
				return err
			}
			matchers[i] = m
		}

		if c.Operator == operators.And {
			return sq.And(matchers...)
		} else {
			return sq.Or(matchers...)
		}
	default:
		return errors.Errorf("unsupported operator: %s", c.Operator)
	}
}
