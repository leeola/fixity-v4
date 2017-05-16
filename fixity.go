package fixity

import (
	"encoding/json"
	"io"
	"time"

	"github.com/leeola/fixity/q"
)

// Fixity implements writing, indexing and reading with a Fixity store.
//
// This interface will be implemented for multiple stores, such as a local on
// disk store and a remote over network store.
type Fixity interface {
	// Blob returns a raw blob of the given hash.
	//
	// Mainly useful for inspecting the underlying data structure.
	Blob(hash string) ([]byte, error)

	// ReadHash unmarshals the given hash contents into a Version.
	//
	// Included in the Version is the Json and MultiBlob, if any exist. If no
	// Json exists the Json struct will be zero value, and if no MultiBlob
	// exists the ReadCloser will be nil.
	//
	// ReadHash will return ErrNotVersion if the given hash is not a valid hash.
	ReadHash(hash string) (Version, error)

	// ReadId unmarshals the given id into a Version struct.
	//
	// Included in the Version is the Json and MultiBlob, if any exist. If no
	// Json exists the Json struct will be zero value, and if no MultiBlob
	// exists the ReadCloser will be nil.
	ReadId(id string) (Version, error)

	// Search for documents matching the given query.
	Search(*q.Query) ([]string, error)

	// Write the given Commit, Reader, and JsonWithMeta to the Fixity store.
	//
	// A single Write can support an arbitrary number of JsonWithMeta fields.
	// Each Json value within the JsonWithMeta is stored as it's own content
	// address.
	//
	// This allows the caller to optimize how the data is stored. Ensuring that
	// frequently changing data is not stored with infrequently changing data,
	// effectively manually deduplicating the json.
	//
	// This method of deduplication, vs rolling checksums as seen in Blobs,
	// is chosen because the caller of Write is able to effectively choose
	// the rolling splits by seperating Json out into separate objects.
	Write(Commit, io.Reader, MultiJson) ([]string, error)

	// TODO(leeola): Enable a close method to shutdown any
	//
	// // Close shuts down any connections that may need to be closed.
	// Close() error
}

// Commit is a higher level Version, allowing simple and high level writes.
//
// Many or all fields may be duplicated from the Version struct. See Version
// for documentation on them.
type Commit struct {
	Id                  string     `json:"id,omitempty"`
	PreviousVersionHash string     `json:"previousVersion,omitempty"`
	UploadedAt          *time.Time `json:"uploadedAt,omitempty"`
	ChangeLog           string     `json:"changeLog,omitempty"`
}

// Version of json and blob data tracked through history and time.
//
// This is the root method for tracking mutation in Fixity. Each write to Fixity
// writes the json and blob data and records their addresses here in this
// struct along with some additional metadata.
//
// Note that many of these fields are optional, and it is up to the Fixity
// implementation to enforce reasonable requirements.
type Version struct {
	// MultiJson is a map of JsonWithMeta values.
	//
	// Each stored JsonHash is paired with an optional JsonMeta field describing
	// indexing metadata for the stored Json.
	MultiJson MultiJson `json:"multiJson,omitempty"`

	// MultiBlobHash is the hash address of any blob data stored for this version.
	//
	// See MultiBlob and Blob docstrings for further explanation of the MultiBlob.
	MultiBlobHash string `json:"multiBlobHash,omitempty"`

	// Id is a unique string which allows Versions to be linked.
	//
	// Since Fixity is immutable, Versions allow a single piece of data to be
	// mutated over time and history. Each version represents a single state
	// of mutation for the given Json and Blob hash. The Id, allows each
	// version of, say, a single File or Wiki page to have the same identifier
	// and represent the same item.
	//
	// Ids can be random or contain meaning, the usage is entirely up to the
	// user.
	Id string `json:"id,omitempty"`

	// UploadedAt is used to track the Version over time, and sort the most recent.
	//
	// This is important, as many versions of a single id have to be sorted somehow.
	// Sorting them by PreviousVersionCount and PreviousVersionHash is possible,
	// but that leads itself to conflicts which then have to be resolved, merged,
	// etc.
	//
	// Sorting by time allows for automatic resolution of any conflict, and is
	// the most hands-free method of conflict resolution. Not guaranteed to be
	// correct, but guaranteed to be easy.
	UploadedAt *time.Time `json:"uploadedAt,omitempty"`

	// PreviousVersionHash stores the Version preceeding this Version, if any.
	//
	// This not only provides a historical record of each mutation, but it can
	// help identify version forks. A fork in this case, is when multiple
	// writers write based off of the same PreviousVersionHash. Since Fixity
	// stores data by content address, forks and "conflicts" are not
	// problematic, but can cause confusion to the actual writer of the data.
	PreviousVersionHash string `json:"previousVersion,omitempty"`

	// PreviousVersionCount stores a count of all previous versions.
	//
	// This serves to provide a more human friendly method of knowing how many
	// modifications there were, without having to run through the entire
	// PreviousVersion chain.
	PreviousVersionCount int `json:"previousVersionCount,omitempty"`

	// ChangeLog is a simple human friendly message about this Version.
	ChangeLog string `json:"changeLog,omitempty"`

	// MultiBlob is the read contents of the MultiBlobHash.
	//
	// This is loaded for convenience during Fixity.Read methods. It is not
	// stored within the marshalled value of JsonWithMeta.
	//
	// This must be closed if not nil!
	MultiBlob io.ReadCloser `json:"-"`
}

// MultiJson is a map of JsonWithMetas, keyed for unordered unmarshalling.
type MultiJson map[string]JsonWithMeta

// JsonWithMeta stores the hash and meta of a Json struct.
type JsonWithMeta struct {
	// JsonHash is the hash address of any json data stored for this version.
	//
	// See Json docstring for further explanation of Json.
	JsonHash string `json:"jsonHash,omitempty"`

	// JsonMeta stores information about the raw Json being stored.
	//
	// This is primarily used to provide insights on how to index and unmarshal
	// the Json struct.
	//
	// See JsonMeta docstring for further details.
	JsonMeta *JsonMeta `json:"jsonMeta,omitempty"`

	// Json is the read contents of the JsonHash.
	//
	// This is loaded for convenience during Fixity.Read methods. It is not
	// stored within the marshalled value of JsonWithMeta.
	Json Json `json:"-"`
}

// Json is a struct which stores text data in Json form.
//
// This data is often indexed, and is the method by which Blob data stores
// and indexes metadata about that blob data. It does not require or imply
// that blob data exists with the given Json, as the Json may be the primary
// data being stored. As is the case with a Wiki, etc.
type Json struct {
	// Json is the actual json data being stored.
	//
	// Note that Fixity provides some helpers to marshal/unmarshal the Json
	// struct into an interface as well as automatic index field inspecting,
	// which assumes valid Json, but if those are not used this Json then the
	// []byte slice can be anything.
	Json json.RawMessage `json:"json"`
}

// JsonMeta stores information about the raw Json being stored.
//
// This serves, for example, to ensure that if the Index is rebuilt,
// it always knows which fields of the Json data need to be indexed.
// As well as mappings for json fields, etc.
//
// Without Metadata about the Json data, Json data would become a black box
// with no information to help Fixity rebuild indexes and etc.
type JsonMeta struct {
	// IndexFields are the fields of the Json data to be indexed.
	//
	// These can include the value if the indexer cannot assert the real
	// value to be indexed from the Json.Json []byte slice.
	IndexedFields Fields `json:"indexedFields"`
}

// MultiBlob stores the Blob addresses of a piece of data.
//
// The data, say an Image, is split up into multiple Blobs as to allow
// for the content to be dedupicated.
//
// TODO(leeola): add a TotalSize field.
type MultiBlob struct {
	BlobHashes []string `json:"blobHashes"`
}

// Blob is a chunk of MultiBlob data, serving to deduplicate large content.
//
// TODO(leeola): add a Size field.
type Blob struct {
	Blob []byte `json:"blob"`
}
