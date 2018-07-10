package index

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/q"
)

type QueryIndexer interface {
	Indexer
	Querier
}

type Indexer interface {
	Index(mutRef fixity.Ref, m fixity.Mutation, d *fixity.DataSchema, v fixity.Values) error
}

// TODO(leeola): articulate a mechanism to query against unique ids or
// versions.
type Querier interface {
	Query(q.Query) ([]fixity.Match, error)
}

const (
	FIDKey       string = "fid"
	FRefKey      string = "fref"
	FSizeKey     string = "fsize"
	FChecksumKey string = "fchecksum"
)
