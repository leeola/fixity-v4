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

	// TODO(leeola): implement once Nodes are being worked on.
	//
	// // ReadHashes reads each hash in order until one returns content.
	// //
	// // This method exists primarily to allow Fixity nodes to request a hash,
	// // and any mapped hashes (encrypted hashes/etc) from subsequent nodes.
	// ReadHashes(hashes []string) (Version, error)

	// ReadId unmarshals the given id into a Version struct.
	//
	// Included in the Version is the Json and MultiBlob, if any exist. If no
	// Json exists the Json struct will be zero value, and if no MultiBlob
	// exists the ReadCloser will be nil.
	ReadId(id string) (Version, error)

	// Search for documents matching the given query.
	Search(*q.Query) ([]string, error)

	// Write the given Commit, MultiJson, and Reader to the Fixity store.
	//
	// A single write can support an arbitrary number of Json documents
	// via the MultiJson map. The reasoning behind this is documented in
	// the MultiJson docstring.
	Write(Commit, MultiJson, io.Reader) ([]string, error)

	// TODO(leeola): Writeblob is disabled until the mapping API/Schema is figured
	// out. Specifically, it's clear how we can store maps of one hash to another,
	// but how the requesting of the original hash from the node network is still
	// up for debate.
	//
	// Eg, if a hash is requested from a node and not found, the
	// Eg, a network of nodes could contain both the hash and the alternate-hash
	// (encrypted, alternate storage like ipfs, etc), how do we request it from
	// the network. Making multiple requests, eg for the hash and then alternate,
	// to each node is bad.
	//
	// Currently i'm thinking requests will propagate via cascading ReadHashes.
	// See ReadHashes docstring for more thoughts on this.
	//
	// // WriteBlob writes the given bytes to the store, returning the hash address.
	// //
	// // This is a lower level function allowing two Fixity nodes to write content
	// // to eachother.
	// //
	// // It is expected that the resulting hash **may not match** the hash
	// // the writer expected. If this happens, a map should be created.
	// WriteBlob([]byte) (string, error)

	// TODO(leeola): Enable a close method to shutdown any
	//
	// // Close shuts down any connections that may need to be closed.
	// Close() error
}

// Commit contains ordering and mutation info for the data being written.
//
// Eg, the Id to group writes together, the PreviousVersionHash to load
// mutations and/or order, and the CreatedAt  to represent timed order.
//
// Most fields are optional, depending on the Fixity and Index implementations.
type Commit struct {
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

	// ChangeLog is a simple human friendly message about this Version.
	ChangeLog string `json:"changeLog,omitempty"`
}

// MultiJson is a JsonWithMetas map, keyed for unordered unmarshalling.
//
// MultiJson differs from MultiJsonHash in that MultiJson is supplied by
// users, and contains the JsonBytes. MultiJsonHash is stored within
// the Fixity.Store, and does *not* contain the JsonBytes. The Bytes are
// stored separately, as to separate the Meta from the actual Content.
//
// MultiJson and MultiJsonHash allow a writer to store multiple json structs
// together, within a single Commit.
//
// A single Commit Write can support an arbitrary number of Json documents via
// the MultiJson map. Each Json value within the JsonWithMeta is stored as
// it's own content address.
//
// This allows the caller to optimize how the data is stored. Ensuring that
// frequently changing data is not stored with infrequently changing data,
// effectively manually deduplicating the json.
//
// This method of deduplication, vs rolling checksums as seen in Blobs,
// is chosen because the caller of Write is able to effectively choose
// the rolling splits by seperating Json out into separate objects.
// Furthermore, for rolling checksums to be effective with smaller documents
// the rolling algorithm would need to chunk at very small intervals,
// introducing a lot of extra documents in the store with little gain.
//
// Finally, and most importantly, storing Json as chunked bytes would cause
// the json to effectively be encoded. No longer is the content "just json",
// but rather you need to join bytes together to construct your actual data,
// as is the case with binary blobs. Blobs don't have a choice on this, as
// Binary isn't Json, but Json does. Keeping the storage model easy to reason
// about and easy to migrate away from, analyze with external tools, etc,
// is a core philosophy of Fixity.
type MultiJson map[string]JsonWithMeta

// JsonWithMeta stores the bytes and meta of a Json struct.
type JsonWithMeta struct {
	Json

	// JsonMeta stores information about the raw Json being stored.
	//
	// This is primarily used to provide insights on how to index and unmarshal
	// the Json struct.
	//
	// See JsonMeta docstring for further details.
	JsonMeta *JsonMeta `json:"jsonMeta,omitempty"`
}

// Json is a struct which stores text data in Json form.
//
// This data is often indexed, and is the method by which Blob data stores
// and indexes metadata about that blob data. It does not require or imply
// that blob data exists with the given Json, as the Json may be the primary
// data being stored. As is the case with a Wiki, etc.
type Json struct {
	// JsonBytes is the actual json data being stored.
	JsonBytes json.RawMessage `json:"jsonBytes"`
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
