package index

// NOTE(leeola): This interface will likely be in heavy flux. Do not trust it :)
type Index interface {
	// Query the index for a single result.
	QueryOne(Query) (Result, error)

	// Query the Index with the given fields.
	Query(Query) (Results, error)
}

type Query struct {
	// If supplied, an error (ErrIndexVersionsDoNotMatch) will be returned if the
	// expected IndexVersion and the current IndexVersion do not match.
	//
	// This is useful for query fields that change between index rebuilds
	// *(IndexEntry, mainly)*.
	IndexVersion string `json:"indexVersion"`

	// An unordered index of all the hashes in the database.
	//
	// Order is not gauranteed to be the same between re-indexes.
	FromEntry int `json:"fromEntry"`

	// Limit the results to this value.
	//
	// Note that this is ignored with QueryOne.
	Limit int `json:"limit"`
}

// PinQuery is a subset of a Query contaning fields logical to Peer pinning.
type PinQuery struct {
	// no query fields at the moment.
}

type Result struct {
	// See Index.Version() and Query.IndexVersion
	IndexVersion string `json:"indexVersion"`
	Hash         string `json:"hash"`
}

type Results struct {
	// See Index.Version() and Query.IndexVersion
	IndexVersion string   `json:"indexVersion"`
	Hashes       []string `json:"hashes"`
}

func (q Query) IsZero() bool {
	switch {
	case q.FromEntry != 0:
		return false
	case q.IndexVersion != "":
		return false
	case q.Limit != 0:
		return false
	default:
		return true
	}
}

func (pq PinQuery) IsZero() bool {
	switch {
	default:
		return true
	}
}

// CommaString returns a comma delimited string of all the fields.
//
// IMPORTANT: The order of the fields *must be the same* for each individual
// query.
//
// CommaString is used to store the paginate (indexEntry) data for a remote peer,
// and CommaString is used as a key for that data.
func (pq PinQuery) CommaString() string {
	// There are no fields to delimit currently
	return ""
}
