package fixity

import (
	"strings"

	base58 "github.com/jbenet/go-base58"
)

func NewRef(hasher string, b []byte) Ref {
	return Ref(hasher + "-" + base58.Encode(b))
}

func (r Ref) Decode() []byte {
	h := nthSplit("-", string(r), 2, 1)
	return base58.Decode(h)
}

func (r Ref) HasherName() string {
	return nthSplit("-", string(r), 2, 0)
}

func nthSplit(char, s string, n, i int) string {
	split := strings.SplitN(char, s, n)
	if len(split) <= i {
		return ""
	}

	return split[i]
}
