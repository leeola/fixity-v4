package snail

import (
	"errors"

	"github.com/blevesearch/bleve"
	"github.com/leeola/fixity/q"
)

type Matcher interface {
	Match(Doc, q.Constraint) bool
}

type MatcherFunc func(Doc, q.Constraint) bool

func (f MatcherFunc) Match(d Doc, c q.Constraint) bool {
	return f(d, c)
}

func andMatcher(matchers []Matcher) Matcher {
	return MatcherFunc(func(d Doc, c q.Constraint) bool {
		for _, m := range matchers {
			if !m.Match(d, c) {
				return false
			}
		}
		return true
	})
}

func eqMatcher(d Doc, c q.Constraint) bool {
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

	return MatcherFunc(func(d Doc, _ q.Constraint) bool {
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
