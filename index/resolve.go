package index

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
)

// ResolveHashOrAnchor returns a hash for the given anchor.
//
// If the given string is a hash, not an anchor, it will be returned directly.
// This works by checking if the address exists, and returns it immediately if
// it does.
func ResolveHashOrAnchor(s store.Store, i Queryer, a string) (string, error) {
	exists, err := s.Exists(a)
	if err != nil {
		return "", errors.Stack(err)
	}

	if exists {
		return a, nil
	}

	q := Query{
		Metadata: Metadata{
			"anchor": a,
		},
	}

	result, err := i.QueryOne(q)
	if err != nil {
		return "", errors.Wrap(err, "failed to query anchor")
	}

	if result.Hash.Hash == "" {
		return "", errors.New("index: anchor not found")
	}

	return result.Hash.Hash, nil
}
