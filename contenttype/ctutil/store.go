package ctutil

import (
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

func WritePartRoller(s store.Store, i index.Indexer, r store.PartRoller) ([]string, error) {
	var hashes []string
	for {
		c, err := r.Roll()
		if err != nil && err != io.EOF {
			return nil, errors.Stack(err)
		}

		if err == io.EOF {
			break
		}

		h, err := store.WritePart(s, c)
		if err != nil {
			return nil, errors.Stack(err)
		}

		if err := i.Entry(h); err != nil {
			return nil, errors.Stack(err)
		}

		hashes = append(hashes, h)
	}

	return hashes, nil
}
