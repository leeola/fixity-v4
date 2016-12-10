package index

// Indexer is used to index content being written and is primarily called from
// the ContentType Uploaders as they are the furthest downstream and know the
// most about the data being uploaded.
type Indexer interface {
	Metadata(h string, m Indexable) error
}

// EntryIndexer is a special interface that the Kala Node uses to index each blob.
//
// It is separated from the Indexer interface because ContentType Uploaders do not
// and should not need to index entries themselves.
type EntryIndexer interface {
	// Entry indexes the given hash with no additional metadata.
	//
	// IMPORTANT: Kala expects the index entry to be auto incrementing, as that is
	// how an entry is determined to be new or old by its peers.
	Entry(h string) error
}

type Queryer interface {
	// Query the index for a single result.
	QueryOne(Query) (Result, error)

	// Query the Index with the given fields.
	Query(Query) (Results, error)
}

// Indexable converts any data type that can return a metadata type to be indexed.
type Indexable interface {
	ToMetadata() Metadata
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

type Metadata map[string]string

// PinQuery is a subset of a Query contaning fields logical to Peer pinning.
type PinQuery struct {
	// no query fields at the moment.
}

type Result struct {
	// See Index.Version() and Query.IndexVersion
	IndexVersion string `json:"indexVersion"`
	Hash         Hash   `json:"hash"`
}

type Results struct {
	// See Index.Version() and Query.IndexVersion
	IndexVersion string `json:"indexVersion"`
	Hashes       []Hash `json:"hashes"`
}

type Hash struct {
	Entry int    `json:"entry"`
	Hash  string `json:"hash"`
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
