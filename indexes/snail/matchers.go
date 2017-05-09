package snail

import "github.com/leeola/fixity/q"

type Matcher func(q.Constraint, DocFields) bool

func andMatcher(matchers []Matcher) Matcher {
	return func(c q.Constraint, d DocFields) bool {
		for _, m := range matchers {
			if !m(c, d) {
				return false
			}
		}
		return true
	}
}

func eqMatcher(c q.Constraint, d DocFields) bool {
	dv, ok := d[c.Field]
	if !ok {
		return false
	}

	return c.Value == dv
}
