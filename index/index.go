package index

import "github.com/leeola/kala/q"

type Options map[string]interface{}

type Field struct {
	Field    string `json:"field"`
	Value    interface{} `json:"value"`
	Options  Options `json:"options"`
}

// Index implements indexing and searching functionality for a kala store.
type Index interface {
	Index(id string, []Field) error
	Search(q.Query) ([][]byte, error)
}
