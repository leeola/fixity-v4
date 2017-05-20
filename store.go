package fixity

import "io"

type Store interface {
	// Check if the given hash exists in the Store
	Exists(string) (bool, error)

	// Takes a hex string of the content hash, and returns a reader for the content
	Read(string) (io.ReadCloser, error)

	// Write raw data to the store.
	//
	// Return the hash of the written data.
	Write([]byte) (string, error)

	// Write the given data to the store only if it matches the given hash.
	//
	// Note that this must compute the hash to ensure the bytes match the given hex
	// hash.
	WriteHash(string, []byte) error

	// List records in the store.
	//
	// IMPORTANT: Listing may not be deterministic and does not ensure that new records
	// or removed records are included in the listing. Therefor Listing should be done
	// before before a store is being actively served.
	List() (<-chan string, error)

	// TODO(leeola): Enable a close method to shutdown any
	//
	// // Close shuts down any connections that may need to be closed.
	// Close() error
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
	Commit

	// MultiJsonHash is a map of JsonHashWithMeta values.
	//
	// Each stored JsonHash is paired with an optional JsonMeta field describing
	// indexing metadata for the stored Json.
	//
	// See MultiJsonHash docstring for further explanation.
	MultiJsonHash MultiJsonHash `json:"multiJsonHash,omitempty"`

	// MultiBlobHash is the hash address of any blob data stored for this version.
	//
	// This is stored by address (hash) rather than embedded as MultiJsonHash is,
	// because MultiBlob is significantly bigger, and can grow basically without
	// limit. The MultiJson and MultiJsonHash structs are expected to store far
	// less data.
	//
	// See MultiBlob and Blob docstrings for further explanation of the MultiBlob.
	MultiBlobHash string `json:"multiBlobHash,omitempty"`

	// PreviousVersionCount stores a count of all previous versions.
	//
	// This serves to provide a more human friendly method of knowing how many
	// modifications there were, without having to run through the entire
	// PreviousVersion chain.
	PreviousVersionCount int `json:"previousVersionCount,omitempty"`

	// MultiBlob is the read contents of the MultiBlobHash.
	//
	// This is loaded for convenience during Fixity.Read methods. It is not
	// stored within the marshalled value of JsonWithMeta.
	//
	// This must be closed if not nil!
	MultiBlob io.ReadCloser `json:"-"`

	// MultiJson is the read contents of the MultiJsonHash hashes.
	//
	// This is loaded for convenience during Fixity.Read methods. It is not
	// stored within the marshalled value of JsonWithMeta.
	MultiJson MultiJson `json:"-"`
}

// MultiJsonHash is a JsonHashWithMetas map, keyed for unordered unmarshalling.
type MultiJsonHash map[string]JsonHashWithMeta

// JsonWithMeta stores the hash and meta of a Json struct.
type JsonHashWithMeta struct {
	JsonWithMeta

	// JsonHash is the hash address of of the json data.
	//
	// See Json docstring for further explanation of Json.
	JsonHash string `json:"jsonHash,omitempty"`

	// JsonBytes hides the JsonBytes field from the embedded JsonWithMeta field.
	//
	// This serves to prevent it from being written in the store. Note that
	// it is a pointer because a struct{} alone would still cause an empty
	// object to be written.
	JsonBytes *struct{} `json:"jsonBytes,omitempty"`
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
	BlobBytes []byte `json:"blob"`
}
