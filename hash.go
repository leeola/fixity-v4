package fixity

import (
	"fmt"
	"hash"

	"github.com/dchest/blake2b"
	multihash "github.com/multiformats/go-multihash"
)

const (
	blake2b256 = "blake2b-256"

	// DefaultMultihashName is the hasher function name from the multihash
	// library that is being used for new fixity hashes.
	//
	// Older fixity Refs may be a different value.
	DefaultMultihashName = blake2b256
)

var (
	multihashCode uint64
)

func init() {
	c, ok := multihash.Names[DefaultMultihashName]
	if !ok {
		panic(fmt.Sprintf("multihash name not found: %q", DefaultMultihashName))
	}
	multihashCode = c

}

// Hash provides a central hash function for any stores.
//
// For convenience, this function returns the Ref not the hashed bytes.
func Hash(b []byte) (Ref, error) {
	// -1 uses the default size for the given code.
	mh, err := multihash.Sum(b, multihashCode, -1)
	if err != nil {
		return "", fmt.Errorf("sum: %v", err)
	}

	return NewRef(mh), nil
}

// Hasher returns a *non-multihash* hash.Hash interface allowing incremental
// writes to generate a sum.
func Hasher(multihashName string) (hash.Hash, error) {
	switch multihashName {
	case blake2b256:
		return blake2b.New256(), nil
	default:
		return nil, fmt.Errorf("unexpected multihash name: %q", multihashName)
	}
}
