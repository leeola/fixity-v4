package fixity

import (
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/fixity/q"
)

// Fixity implements user focused writing and reading of data.
//
// This interface will be implemented for multiple stores, such as a local on
// disk store and a remote over network store.
type Fixity interface {
	// Blob returns a raw blob of the given hash.
	//
	// Mainly useful for inspecting the underlying data structure.
	Blob(hash string) (io.ReadCloser, error)

	// Blockchain allows one to manage and inspect the Fixity Blockchain.
	//
	// The blockchain is low level and should be used with care. See Blockchain
	// docstring for further details.
	Blockchain() Blockchain

	// Delete the given id's content from the fixity store.
	//
	// Each Content, Blob and Chunk will be deleted if no other block in the
	// blockchain depends on it. Verifying this is done by the garbage
	// collector and is a slow process.
	//
	// All blocks for the given id will be removed from the blockchain.
	Delete(id string) error

	// Read the latest Content with the given id.
	Read(id string) (Content, error)

	// Read the Content with the given hash.
	ReadHash(hash string) (Content, error)

	// Write the given reader to the fixity store and index fields.
	Write(id string, r io.Reader, f ...Field) ([]string, error)

	// Search for documents matching the given query.
	Search(*q.Query) ([]string, error)

	// Close shuts down any connections that may need to be closed.
	Close() error
}

// Blockchain implements low level block management methods for Fixity.
//
// The fixity blockchain does not contain any traditional proof of work
// and does not have anything related to crypto or crypto currencies.
// The name was chosen in hopes to clearly express the ledger side of
// blockchains. While there will be many similarities in how the Fixity
// blockchain uses immutability and history, there will also be many
// differences from popular blockchains.
//
// A blockchain in Fixity serves as a ledger for what content exists
// on a distributed Fixity network. If a hash address cannot be found
// within one of the Blocks on the blockchain, such as the hash of
// a Content, Blob or Chunk, then it is considered available for
// garbage collection and will be removed.
//
// The Fixity blockchain consists of three main parts: The Block number,
// the PreviousBlockHash and some data effectively giving the Block a
// "type".
//
// The Block number is an ever incrementing value and with the help of
// PreviousBlockHash it provides a way for distributed nodes to achieve
// eventual consensus.
//
// The PreviousBlockHash provides a way to track the entire blockchain
// from the Head() Block. It also provides a way for the blockchain itself
// to be mutable, in an immutable environment. More on mutability soon.
//
// The Block type is the reason why the block exists. It may have added
// Content, removed Content, mutated the chain, etc.
//
// Mutability of the blockchain is achieved by appending new blocks that
// skip one or more blocks in their PreviousBlockHash. For example, with
// a blockchain of 5 blocks, Block 6 could be written with a
// PreviousBlockHash set to the hash of Block 4. Since Block 6 is the head,
// traversing the blockchain would look like: Block 6 -> Block 4 -> Block 3
// and so on. Note that Blocks start with 0 index, but for these examples
// we're not using zero index.
//
// To achieve mutability on blocks that aren't currently the Head()
// as Block 5 was in the previous example, all Blocks from the target
// Block to the Head() must be rewritten to the blockchain. In order.
// This means that removing old blocks can be costly and slow.
//
// Fixity strives to keep the ledger as a trustable and easy to verify
// chain. An alternative to block skipping would be to write a content
// deletion block, essentially writing to the ledger that content is
// to be garbage collected. However, verifying the ever growing blockchain
// would mean needing to reference these deletion blocks frequently to know
// what content should and shouldn't be looked into. In otherwords,
// content on the blockchain may have a deletion block for it further up
// the chain, so verifying content of the ledger becomes difficult.
//
// The chosen method of content skipping does most of the difficult work
// up front and results in a very clean ledger. It does this at the cost
// of needing to complicate the removal/skipping process.
//
// This interface focuses on all of the above functionality.
type Blockchain interface {
	// AppendBlocks locks the store and writes the given blocks in order.
	//
	// The field PreviousBlockHash's value of all blocks *must* be empty.
	//
	// The returned Block array will contain the new hashes of the given
	// blocks.
	AppendBlocks(appendTo Block, blocks []Block) ([]Block, error)

	// Head returns the latest block in the blockchain.
	Head() (Block, error)

	// SkipBlock removes the given block from the blockchain.
	SkipBlock(Block) ([]Block, error)
}

// Block serves as a ledger for mutations of the fixity datastore.
//
// Each block stores an always incrementing Block number and a hash of
// the previous block in the chain. These two fields allow a fixity
// store to be iterated through the always appending history.
//
// While the history is always appending, previous blocks may be skipped,
// effectively removing them from the history of the blockchain. This is
// done by writing a new block whose PreviousBlockHash value skips one or
// more previous blocks in the chain.
type Block struct {
	// Block is the ever incrementing block number for this block.
	//
	// Each block will be incremented from the previous block.
	Block int `json:"block"`

	// PreviousBlockHash is the hash of the block that came before this.
	//
	// Note that the blockchain itself is mutable, such that the
	// PreviousBlockHash isn't guaranteed to have the block number of Block-1.
	// If a block was skipped, the block numbers may differ.
	//
	// See FixityBlockchain.SkipBlock for more information on block skipping
	// and implications of that.
	PreviousBlockHash string `json:"previousBlockHash"`

	// Skip contains Skip data and makes this Block a Skip Block.
	Skip *Skip `json:"skip,omitempty"`

	// ContentHash contains the ContentHash and makes this block a Content block.
	ContentHash string `json:"cotentHash,omitempty"`

	// BlockHash is the hash of the Block itself, provided by Fixity.
	BlockHash string `json:"-"`

	// Store allows block method(s) to load previous blocks and content.
	Store Store `json:"-"`
}

// Skip blocks provide information about the block that was skipped.
type Skip struct {
	// BlockHash of the block to be skipped.
	BlockHash string `json:"blockHash"`
}

// Content stores blob, index and history information for Fixity content.
type Content struct {
	// Id provides a user friendly way to reference a chain of Contents.
	//
	// History of Content is tracked through the PreviousContentHash chain,
	// however that does not provide a clear single identity for users.
	// The id field allows this, can be indexed and assocoated and is
	// easy to conceptualize.
	Id string `json:"id,omitempty"`

	// PreviousContentHash stores the previous Content for this Content.
	//
	// This allows a single entity, such as a file or a database "record"
	// to be mutated through time. To reference this history of contents,
	// the Id is used.
	PreviousContentHash string `json:"previousContentHash,omitempty"`

	// BlobHash is the hash of the  Blob containing this content's data.
	BlobHash string `json:"blobHash"`

	// IndexedFields contains the indexed metadata for this content.
	//
	// This allows the content to be searched for and can be used to
	// store basic metadata about the content.
	IndexedFields Fields `json:"indexedFields,omitempty"`

	// ReadCloser allows the Content to be read from directly.
	//
	// TODO(leeola): Remove this in favor of a Blob() method which
	// returns the Blob & embedded reader. This saves creation of the
	// reader until it's been explicitly requested.
	//
	// This value is not stored.
	io.ReadCloser `json:"-"`

	// ContentHash is the hash of the Content itself, provided by Fixity.
	//
	// This value is not stored.
	ContentHash string `json:"-"`

	// Store allows block method(s) to load previous content.
	//
	// This value is not stored.
	Store Store `json:"-"`
}

// Blob stores a series of ordered ChunkHashes
type Blob struct {
	ChunkHashes []string `json:"chunkHashes"`
	Size        int64    `json:"size,omitempty"`
	RollSize    int      `json:"rollSize,omitempty"`

	// NextBlobHash is not currently supported / implemented anywhere, but
	// is required for very large storage. Eg, if there are so many chunks
	// for a given dataset that it cannot be stored in memory during writing
	// and reading, then we will need to split them up via NextBlobHash.
	//
	// // NextBlobHash stores another blob which is to be appended to this blob.
	// //
	// // This serves to allow very large blobs that cannot be loaded entirely
	// // into to memory to be split up into many parts.
	// NextBlobHash string `json:"nextBlobHash,omitempty"`

	// ReadCloser allows the Blob to be read from directly.
	//
	// This value is not stored.
	io.ReadCloser `json:"-"`
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

	var previousBlock Block
	err := readAndUnmarshal(b.Store, b.PreviousBlockHash, &previousBlock)
	if err != nil {
		return Block{}, err
	}

	previousBlock.BlockHash = b.PreviousBlockHash
	previousBlock.Store = b.Store

	return previousBlock, nil
}

func (b *Block) Content() (Content, error) {
	if b.Store == nil {
		return Content{}, errors.New("Store not set")
	}

	if b.ContentHash == "" {
		return Content{}, errors.New("contentHash is empty")
	}

	var c Content
	err := readAndUnmarshal(b.Store, b.ContentHash, &c)
	if err != nil {
		return Content{}, err
	}

	return c, nil
}
