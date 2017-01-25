package index

import (
	"io"

	"github.com/leeola/kala/store"
)

// Index is the full Index interface.
type Index interface {
	Indexer
	Queryer

	// Reset deletes the entire index so it can be written from a fresh state.
	//
	// This is, of course, destructive. The index should be fully rebuilt after
	// calling this.
	Reset() error
}

// Indexer is used to index content being written and is primarily called from
// the ContentType Uploaders as they are the furthest downstream and know the
// most about the data being uploaded.
type Indexer interface {
	// Version indexes the given Version and Metadata together.
	//
	// This will be the main use case of
	Version(h string, v store.Version, m interface{}) error

	// Entry indexes the given hash with no additional metadata.
	Entry(h string) error

	// Version returns the unique index version.
	//
	// The returned value must be unique each time the index is reset or created.
	IndexVersion() string
}

type Queryer interface {
	// Query the index for a single result.
	//
	// SortBy must be optional.
	QueryOne(Query, ...SortBy) (Result, error)

	// Query the Index with the given fields.
	Query(Query, ...SortBy) (Results, error)
}

// Lister is responsible for listing the hashes to rebuild an index from.
//
// This is a package-local interface of the Store interface.
type Lister interface {
	Read(string) (io.ReadCloser, error)
	List() (<-chan string, error)
}

type Query struct {
	// If supplied, an error (ErrIndexVersionsDoNotMatch) will be returned if the
	// expected IndexVersion and the current IndexVersion do not match.
	//
	// This is useful for query fields that change between index rebuilds
	// *(IndexEntry, mainly)*.
	IndexVersion string `json:"indexVersion"`

	// By default an Indexer should not include all the different versions of an
	// anchor in a query result. It muddies up the the results, as the same anchor
	// can be in there dozens of times making it difficult to find other anchors.
	//
	// If however you are trying to sort through the revisions of an anchor you can
	// enable this flag to include all uploaded metas. These abide by the same
	// query flags as non versions search, and thus can span multiple anchors if
	// the query is not constrained to a single anchor.
	SearchVersions bool `json:"searchVersions"`

	// An unordered index of all the hashes in the database.
	//
	// Order is not gauranteed to be the same between re-indexes.
	FromEntry int `json:"fromEntry"`

	// Limit the results to this value.
	//
	// Note that this is ignored with QueryOne.
	Limit int `json:"limit"`

	Metadata Metadata `json:"metadata"`
}

// SortBy is a basic struct for passing sorting information to a Query.
type SortBy struct {
	Field      string `json:"field"`
	Descending bool   `json:"descending"`
}

type Metadata map[string]interface{}

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
	Hashes       Hashes `json:"hashes"`
}

type Hash struct {
	Entry int    `json:"entry"`
	Hash  string `json:"hash"`
}

type Hashes []Hash

func (hs Hashes) Len() int {
	return len(hs)
}

func (hs Hashes) Less(i, j int) bool {
	return hs[i].Entry < hs[j].Entry
}

func (hs Hashes) Swap(i, j int) {
	h := hs[i]
	hs[i] = hs[j]
	hs[j] = h
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
