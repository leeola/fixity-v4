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
	Index(mutRef fixity.Ref, m fixity.Mutation, d *fixity.Data, v fixity.Values) error
}

// TODO(leeola): articulate a mechanism to query against unique ids or
// versions.
type Querier interface {
	Query(q.Query) ([]fixity.Ref, error)
}
