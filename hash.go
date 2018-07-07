package fixity

import (
	"fmt"

	multihash "github.com/multiformats/go-multihash"
)

const (
	// DefaultMultihashName is the hasher function name from the multihash
	// library that is being used for new fixity hashes.
	//
	// Older fixity Refs may be a different value.
	DefaultMultihashName = "blake2b-256"
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
