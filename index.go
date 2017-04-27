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
	Index([]Field) error
	Search(q.Query) ([]string, error)
}

type Fields []Field
