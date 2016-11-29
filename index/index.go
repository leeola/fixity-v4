package index

// NOTE(leeola): This interface will likely be in heavy flux. Do not trust it :)
type Index interface {
	// Query the index for a single result.
	QueryOne(Query) (Result, error)

	// Query the Index with the given fields.
	Query(Query) (Results, error)

	// Return a unique identifier for the given **built index**.
	//
	// This is not a version of the code or software behind the index, but rather
	// each time an index is built a version must be generated that will persist
	// for the lifetime of the index.
	//
	// If the index is rebuilt for any reason, a new index version **must** be
	// generated.
	Version() string
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
	IndexEntry int `json:"indexEntry"`
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
