package fixity

import (
	"fmt"

	"github.com/leeola/fixity/config"
	"github.com/leeola/fixity/q"
)

type Index interface {
	Indexer
	Querier
}

type Indexer interface {
	Index(mutRef Ref, m Mutation, d *DataSchema, v Values) error
}

// TODO(leeola): articulate a mechanism to query against unique ids or
// versions.
type Querier interface {
	Query(q.Query) ([]Match, error)
}

type Match struct {
	ID  string `json:"id"`
	Ref Ref    `json:"ref"`
}

func NewIndexFromConfig(name string, c config.Config) (Index, error) {
	if name == "" {
		return nil, fmt.Errorf("empty index name")
	}

	tc, ok := c.IndexConfigs[name]
	if !ok {
		return nil, fmt.Errorf("index name not found: %q", name)
	}

	constructor, ok := indexRegistry[tc.Type]
	if !ok {
		return nil, fmt.Errorf("index type not found: %q", tc.Type)
	}

	ix, err := constructor.New(name, c)
	if err != nil {
		return nil, fmt.Errorf("index constructor %s: %v", name, err)
	}

	return ix, nil
}
