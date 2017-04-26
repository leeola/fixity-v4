package storm

import (
	"github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	"github.com/fatih/structs"
	"github.com/leeola/errors"
	kq "github.com/leeola/kala/q"
	ops "github.com/leeola/kala/q/operators"
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
	stormQuery := db.Select(ms...)

	if kq.SortBy != "" {
		stormQuery.OrderBy(kq.SortBy)
	}
	if kq.SortDescending {
		stormQuery.Reverse()
	}
	if kq.SkipBy != 0 {
		stormQuery.Skip(kq.SkipBy)
	}
	if kq.LimitBy != 0 {
		stormQuery.Limit(kq.LimitBy)
	}

	return stormQuery, nil
}

func ConvertConstraint(c kq.Constraint) (sq.Matcher, error) {
	switch c.Operator {
	case ops.Equal:
		return sq.Eq(c.Field, c.Value), nil

	case ops.GreaterThan:
		return sq.Gt(c.Field, c.Value), nil

	case ops.GreaterThanOrEqual:
		return sq.Gte(c.Field, c.Value), nil

	case ops.In:
		return sq.In(c.Field, c.Value), nil

	case ops.LessThan:
		return sq.Lt(c.Field, c.Value), nil

	case ops.LessThanOrEqual:
		return sq.Lte(c.Field, c.Value), nil

	case ops.Regex:
		v, ok := c.Value.(string)
		if !ok {
			return nil, errors.New("unexpected Regex field value")
		}

		return sq.Re(c.Field, v), nil

	case ops.And, ops.Or, ops.Not:
		matchers := make([]sq.Matcher, len(c.Constraints))
		for i, subc := range c.Constraints {
			m, err := ConvertConstraint(subc)
			if err != nil {
				return nil, err
			}
			matchers[i] = m
		}

		switch c.Operator {
		case ops.And:
			return sq.And(matchers...), nil
		case ops.Or:
			return sq.Or(matchers...), nil
		case ops.Not:
			return sq.Not(matchers...), nil
		default:
			// should never happen
			return nil, errors.Errorf(
				"unexpected operator inside And/Or/Not case: %s", c.Operator)
		}

	default:
		return nil, errors.Errorf("unsupported operator: %s", c.Operator)
	}
}
