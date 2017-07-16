package fixity

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
// "type", such as the ContentBlock and DeleteBlock fields.
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
// Content, removed Content, etc.
//
// Despite the name, a Delete block does not actually remove anything from
// the blockchain. Instead, as the chain is being traversed from the Head(),
// blocks that have had their hash written in a DeleteBlock are simply
// skipped. This method of skipping previous blocks allows multiple
// Fixity nodes to remove and write content on the fly and decide on
// consensus in the future. No altering of the chain ever takes place.
type Blockchain interface {
	// // AppendBlocks locks the store and writes the given blocks in order.
	// //
	// // The field PreviousBlockHash's value of all blocks *must* be empty.
	// //
	// // The returned Block array will contain the new hashes of the given
	// // blocks.
	// AppendBlocks(appendTo Block, blocks []Block) ([]Block, error)

	// AppendContent creates a new block with the given content.
	//
	// If the block content is the same as the current Head() the blockchain
	// may not be progressed.
	AppendContent(Content) (Block, error)

	// Head returns the latest block in the blockchain.
	//
	// Head must return ErrNoPrev if the blockchain is empty.
	Head() (Block, error)

	// DeleteContent writes a delete block  blocks with the given content.
	//
	// It does this by writing a Delete block onto the chain. Iterating
	// through the chain will cause the given contents to be skipped
	// over. The main blockchain is never mutated or altered.
	//
	// This does not remove the content, that is handled by the Fixity
	// implementor.
	DeleteContent(...Content) (Block, error)
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

	// Content contains the content hashl for this content block.
	ContentBlock *ContentBlock `json:"contentBlock,omitempty"`

	// Delete contains data about block(s) that have been deleted.
	//
	// Technically no block is ever deleted on the main blockchain,
	// however any encountered delete block will cause previous
	// blocks to be skipped when navigating backwards via the Previous().
	DeleteBlock *DeleteBlock `json:"deleteBlock,omitempty"`

	// Hash is the hash of the Block itself, provided by Fixity.
	//
	// This value is not stored.
	Hash string `json:"-"`

	// Store allows block method(s) to load previous blocks and content.
	//
	// This value is not stored.
	Store Store `json:"-"`
}

// ContentBlock provides information about content on the blockchain.
type ContentBlock struct {
	// Hash of the Content.
	Hash string `json:"hash"`
}

// DeleteBlock provides information about the block that was deleted.
type DeleteBlock struct {
	// Hashes of the deleted blocks.
	Hashes []string `json:"deletedHashes"`
}
