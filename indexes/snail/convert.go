package snail

import (
	"github.com/leeola/errors"
	"github.com/leeola/fixity/q"
	ops "github.com/leeola/fixity/q/operators"
)

func (s *Snail) convertConstraint(c q.Constraint) (Matcher, error) {
	switch c.Operator {
	case ops.Equal:
		return eqMatcher(c), nil

	case ops.In:
		return inMatcher(c), nil

	case ops.FullTextSearch:
		// TODO(leeola): check if it's a ID search or a Ver search, and return
		// the correct fts for that.
		return ftsMatcher(s.bleveVer, c)

	case ops.And:
		ms, err := s.convertConstraints(c.Constraints)
		if err != nil {
			return nil, err
		}

		return andMatcher(ms, c), nil

	default:
		return nil, errors.Errorf("unsupported operator: %s", c.Operator)
	}
}

func (s *Snail) convertConstraints(cs []q.Constraint) ([]Matcher, error) {
	matchers := make([]Matcher, len(cs))
	for i, c := range cs {
		m, err := s.convertConstraint(c)
		if err != nil {
			return nil, err
		}
		matchers[i] = m
	}
	return matchers, nil
}
