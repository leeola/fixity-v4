package blockchaindb

import (
	"io"

	"github.com/leeola/fixity"
)

type Db interface {
	io.Closer
	Head() (string, error)
	HeadContentBlock(id string) (string, error)
	Update(fixity.Block) error
}

// Block
//
// This is based off of the immutable blockchain and must always
// be kept in sync.
type Block struct {
	Block                int    `json:"block"`
	Hash                 string `json:"hash"`
	ContentId            string `json:"contentId"`
	PreviousContentBlock int    `json:"previousContentBlock"`
}
