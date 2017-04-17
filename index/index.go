package index

import "github.com/leeola/kala/q"

type Options map[string]interface{}

// Index implements indexing and searching functionality for a kala store.
type Index interface {
	// Index the given data with Options describing how to index the data.
	//
	// Note that Options may or may not be required, it depends on the indexer
	// and the desired features. For example, if FullTextSearch is desired
	// one may specify `Options["fullTextSearch"] = "all"` and etc.
	Index(data interface{}, opts Options) error
	Search(q.Query) ([][]byte, error)
}
