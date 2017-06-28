package fixity

import "github.com/leeola/fixity/q"

// Index implements indexing and searching functionality for a fixity store.
//
// An index can be implemented by anything with basic querying abilities,
// usually databases, to retrieve the associated key.
//
// Generally, an index must:
//
// 1. Translate the Fixity *Query into a format that makes sense in the
//    indexer, such as translating *Query to a SQL statement as might be
//    done by sqlite.
// 2. Return hashes of the documents matching the given query.
// 3. Paginate and limit the searches.
// 4. Support searching all fields, if a query constraint has a field value
//    of "*".
//
// It's acceptable that an indexer only implement a subset of these features,
// as the user can choose which indexer to use. Eg, some indexers, especially
// sql based indexers, may have difficulty searching all fields (point 4).
type Index interface {
	// Index the given hash and id to the given fields.
	//
	// The hash and id serve to conceptually index two different things.
	// The hash will allow a search to query all versions, and the id will allow
	// a search to query the latest version of each id.
	Index(hash, id string, fields []Field) error

	// Search the index for the given Query, returning the matching hashes or ids.
	//
	// Note that within the query, a field can be specified as `"*"` to signify
	// all fields are to be checked. Some operators may not support "*", and
	// behavior can differ between Index implementations.
	Search(*q.Query) ([]string, error)

	// TODO(leeola): Enable a close method to shutdown any
	//
	// // Close shuts down any connections that may need to be closed.
	// Close() error
}

// Field of a document to be indexed via Index.Index.
type Field struct {
	Field   string       `json:"field"`
	Value   interface{}  `json:"value,omitempty"`
	Options FieldOptions `json:"options,omitempty"`
}

func (a *Field) Equal(b Field) bool {
	switch {
	case a.Field != b.Field:
		return false
	case a.Value != b.Value:
		return false
	}
	for aK, aV := range a.Options {
		bV, ok := b.Options[aK]
		if !ok {
			return false
		}
		if aV != bV {
			return false
		}
	}
	return true
}

// Fields is a helper type for convenient mutation of a []Field.
type Fields []Field

// Append the given field to this slice.
func (fs *Fields) Append(f Field) {
	*fs = append(*fs, f)
}

func (a Fields) Equal(b Fields) bool {
	if len(a) != len(b) {
		return false
	}
	for i, aF := range a {
		if !aF.Equal(b[i]) {
			return false
		}
	}
	return true
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
