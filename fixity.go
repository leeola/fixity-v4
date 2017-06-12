package fixity

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/leeola/errors"
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
	Blob(hash string) (io.ReadCloser, error)

	Head() (Block, error)

	Read(id string) (Content, error)

	ReadHash(hash string) (Content, error)

	// Remove marks the given block's content to be garbage collected eventually.
	//
	// Each Content, Blob and Chunk will be deleted if no other block in the
	// blockchain depends on it. This is a slow process.
	//
	// If the block is not a content block, an error will be returned.
	Remove(id string) error

	// Write a block for the given reader and index fields.
	Write(id string, r io.Reader, f ...Field) ([]string, error)

	// // Search for documents matching the given query.
	Search(*q.Query) ([]string, error)

	// TODO(leeola): Enable a close method to shutdown any
	//
	// // Close shuts down any connections that may need to be closed.
	// Close() error
}

type Block struct {
	Block             int    `json:"block"`
	PreviousBlockHash string `json:"previousBlockHash"`

	//Deletion  *Deletion  `json:"deletion,omitempty"`
	//Deletions  *Deletion  `json:"deletion,omitempty"`
	//Append  *Append  `json:"append,omitempty"`

	ContentHash string `json:"cotentHash,omitempty"`

	// BlockHash is the hash of the Block itself, provided by Fixity.
	BlockHash string `json:"-"`

	// Store allows block method(s) to load previous blocks and content.
	Store Store `json:"-"`
}

type Deletion struct {
	BlockHash   string `json:"blockHash"`
	ContentHash string `json:"contentHash,omitempty"`
}

type Deletions []Deletion

type Content struct {
	Id                  string `json:"id,omitempty"`
	PreviousContentHash string `json:"previousContentHash,omitempty"`
	BlobHash            string `json:"blobHash"`
	IndexedFields       Fields `json:"indexedFields,omitempty"`

	// ReadCloser allows the Content to be read from directly.
	io.ReadCloser `json:"-"`

	// ContentHash is the hash of the Content itself, provided by Fixity.
	ContentHash string `json:"-"`

	// Store allows block method(s) to load previous content.
	Store Store `json:"-"`
}

type Blob struct {
	ChunkHashes []string `json:"chunkHashes"`
	Size        int64    `json:"size"`
	RollSize    int      `json:"rollSize"`
}

type Chunk struct {
	ChunkBytes []byte `json:"chunkBytes"`
	Size       int64  `json:"size"`
}

func (b *Block) PreviousBlock() (Block, error) {
	if b.Store == nil {
		return Block{}, errors.New("Store not set")
	}

	if b.PreviousBlockHash == "" {
		return Block{}, nil
	}

	rc, err := b.Store.Read(b.PreviousBlockHash)
	if err != nil {
		return Block{}, err
	}
	defer rc.Close()

	bB, err := ioutil.ReadAll(rc)
	if err != nil {
		return Block{}, err
	}

	var previousBlock Block
	if err := json.Unmarshal(bB, &previousBlock); err != nil {
		return Block{}, err
	}

	previousBlock.BlockHash = b.PreviousBlockHash
	previousBlock.Store = b.Store

	return previousBlock, nil
}
