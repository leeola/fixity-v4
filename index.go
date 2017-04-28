package kala

import "github.com/leeola/kala/q"

type Field struct {
	Field   string       `json:"field"`
	Value   interface{}  `json:"value,omitempty"`
	Options FieldOptions `json:"options,omitempty"`
}

type FieldOptions map[string]interface{}

// Index implements indexing and searching functionality for a kala store.
type Index interface {
	// Index the given hash and id to the given fields.
	//
	// Note that the hash and id serve to conceptually index two different things.
	// The hash will allow a search to query all versions
	Index(hash, id string, fields []Field) error

	Search(q.Query) ([]string, error)
}

// Fields is a helper type for convenient mutation of a []Field.
type Fields []Field

// Append the given field to this slice.
func (fs *Fields) Append(f Field) {
	*fs = append(*fs, f)
}
