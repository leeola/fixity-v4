package disk

import (
	"bytes"
	"encoding/hex"
	"path/filepath"

	base58 "github.com/jbenet/go-base58"
)

func (s *Blobstore) pathHash(h string) string {
	// use hex paths on osx because it does not support case sensitive paths.
	h = hex.EncodeToString(base58.Decode(h))

	var buffer bytes.Buffer
	last := len(h) - 1
	for i, char := range h {
		buffer.WriteRune(char)
		if (i+1)%2 == 0 && i != last {
			buffer.WriteRune('/')
		}
	}

	p := buffer.String()
	return filepath.Join(s.path, p)
}
