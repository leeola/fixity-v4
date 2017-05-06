package fixity

import "github.com/leeola/fixity/q"

type Field struct {
	Field   string       `json:"field"`
	Value   interface{}  `json:"value,omitempty"`
	Options FieldOptions `json:"options,omitempty"`
}

// Index implements indexing and searching functionality for a fixity store.
type Index interface {
	// Index the given hash and id to the given fields.
	//
	// The hash and id serve to conceptually index two different things.
	// The hash will allow a search to query all versions, and the id will allow
	// a search to query the latest version of each id.
	Index(hash, id string, fields []Field) error

	Search(*q.Query) ([]string, error)

	// TODO(leeola): Enable a close method to shutdown any
	//
	// // Close shuts down any connections that may need to be closed.
	// Close() error
}

// Fields is a helper type for convenient mutation of a []Field.
type Fields []Field

// Append the given field to this slice.
func (fs *Fields) Append(f Field) {
	*fs = append(*fs, f)
}

// FieldUnmarshaller is responsible for unmarshalling fields from a []byte slice.
//
// This is used to implement Fixity's ability to unmarshal the data and return
// the requested field, thereby indexing the Field.Value even if it was not
// specified.
//
// Fixity will likely call the field unmarshaller many times, so the
// unmarshalled value should be cached between method calls, and lazily unmarshalled
// since it may not be used at all.
type FieldUnmarshaller interface {
	Unmarshal(field string) (value interface{}, err error)
}
