package fixity

import (
	"fmt"

	base58 "github.com/jbenet/go-base58"
	multihash "github.com/multiformats/go-multihash"
)

func NewRef(b []byte) Ref {
	return Ref(base58.Encode(b))
}

func (r Ref) Decode() []byte {
	return base58.Decode(string(r))
}

func (r Ref) HashName() (string, error) {
	mh, err := multihash.FromB58String(string(r))
	if err != nil {
		return "", fmt.Errorf("fromb58string: %v", err)
	}

	decoded, err := multihash.Decode(mh)
	if err != nil {
		return "", fmt.Errorf("fromb58string: %v", err)
	}

	return decoded.Name, nil
}
