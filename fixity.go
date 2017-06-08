package fixity

import "io"

// Fixity implements writing, indexing and reading with a Fixity store.
//
// This interface will be implemented for multiple stores, such as a local on
// disk store and a remote over network store.
type Fixity interface {
	// Blob returns a raw blob of the given hash.
	//
	// Mainly useful for inspecting the underlying data structure.
	Blob(hash string) ([]byte, error)

	// Create a block for the given reader and index fields.
	Create(r io.Reader, f ...[]Field) ([]string, error)

	// Delete marks the given BlobMeta hash to be garbage collected in time.
	//
	// Each Content, Blob and Chunk will be deleted if no other block in the
	// blockchain depends on it. This is a slow process.
	Delete(h string) error

	Update(h string, r io.Reader, f ...[]Field) ([]string, error)

	// // Search for documents matching the given query.
	// Search(*q.Query) ([]string, error)

	// TODO(leeola): Enable a close method to shutdown any
	//
	// // Close shuts down any connections that may need to be closed.
	// Close() error
}

type Block struct {
	Block             int    `json:"block"`
	PreviousBlockHash string `json:"previousBlockHash`

	//Deletion  *Deletion  `json:"deletion,omitempty"`
	//Deletions  *Deletion  `json:"deletion,omitempty"`

	ContentHash string `json:"cotentHash"`
}

type Deletion struct {
	BlockHash   string `json:"blockHash"`
	ContentHash string `json:"contentHash"`
}

type Deletions []Deletion

type Content struct {
	PreviousContentHash string `json:"previousContentHash"`
	BlobHash            string `json:"blobHash"`
	IndexedFields       Fields `json:"indexedFields"`
}

type Blob struct {
	ChunkHashes []string `json:"chunkHashes"`
	RollSize    int      `json:"rollSize"`
}

type Chunk struct {
	ChunkBytes []byte `json:"chunkBytes"`
}
