package snail

import (
	"errors"

	"github.com/blevesearch/bleve"
	"github.com/leeola/fixity/q"
)

type Matcher interface {
	Match(Doc) bool
}

type MatcherFunc func(Doc) bool

func (f MatcherFunc) Match(d Doc) bool {
	return f(d)
}

func andMatcher(matchers []Matcher, c q.Constraint) Matcher {
	return MatcherFunc(func(d Doc) bool {
		for _, m := range matchers {
			if !m.Match(d) {
				return false
			}
		}
		return true
	})
}

func inMatcher(c q.Constraint) Matcher {
	return MatcherFunc(func(d Doc) bool {
		// if it's a wildcard constraint, check all fields.
		if c.Field == "*" {
			for k, _ := range d.Fields {
				c := c
				c.Field = k
				if in(d, c) {
					return true
				}
				return false
			}
		}

		return in(d, c)
	})
}

// in checks of the doc matches the constraint, ignoring wildcard functionality.
//
// The lack of wildcard supports prevents infinite recursion while still
// supporting fields with the name "*"
func in(d Doc, c q.Constraint) bool {
	i, ok := d.Fields[c.Field]
	// if the value doesn't exist, return false
	if !ok {
		return false
	}

	vs, ok := i.([]interface{})

	if !ok {
		return false
	}

	for _, v := range vs {
		if c.Value == v {
			return true
		}
	}
	return false
}

func eqMatcher(c q.Constraint) Matcher {
	return MatcherFunc(func(d Doc) bool {
		// if it's a wildcard constraint, check all fields.
		if c.Field == "*" {
			for _, v := range d.Fields {
				if v == c.Value {
					return true
				}
			}
			return false
		}

		v, ok := d.Fields[c.Field]
		// if the value doesn't exist, return false
		if !ok {
			return false
		}

		return v == c.Value
	})
}

// ftsMatcher checks the bleve index the constraint and compares the match key.
//
// This matcher ignores the supplied fieldValue entirely, and strictly compares
// the key that bleve returns and the key that the DB is checking against.
// This matcher also caches the bleve results, only performing the constraint
// value query once.
//
// IMPORTANT: ftsMatcher will query bleve before the returned matcher is used.
// The main reason for this is to return an error if needed.
func ftsMatcher(b bleve.Index, c q.Constraint) (Matcher, error) {
	queryStr, ok := c.Value.(string)
	if !ok {
		return nil, errors.New("FullTextSearch constraint value must be a string")
	}

	// NOTE: We're querying *all* of the index to cache the results. If this proves
	// too big, we can query on demand and only the specific id in question. This
	// will be slower, but it's the only choice we have, since matchers cannot be
	// delayed.
	bq := bleve.NewMatchPhraseQuery(queryStr)

	if c.Field != "*" {
		bq.SetField(c.Field)
	}

	// TODO(leeola): How many does the search request limit by default? How can
	// we make it return all, no matter the size?
	search := bleve.NewSearchRequest(bq)
	searchResults, err := b.Search(search)
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(searchResults.Hits))
	for i, documentMatch := range searchResults.Hits {
		keys[i] = documentMatch.ID
	}

	return MatcherFunc(func(d Doc) bool {
		// Check if the current document's key is in the bleve list.
		// If it is, we have a match.
		for _, k := range keys {
			if d.Key == k {
				return true
			}
		}
		return false
	}), nil
}
