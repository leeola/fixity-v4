package blev

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
)

func (b *Bleve) Metadata(h string, i index.Indexable) error {
	exists, err := b.keyExists(h)
	if err != nil {
		return errors.Stack(err)
	}
	// If the hash exists, no need to index it again, we already did.
	if exists {
		return nil
	}

	docCount, err := b.index.DocCount()
	if err != nil {
		return errors.Stack(err)
	}

	// Increasing the doc count is important!
	// Due to go zero values and how kala is designed, zero values are ignored for
	// many things, including entries. So the 0th entry is the 1st index.
	docCount += 1

	m := i.ToMetadata()
	m["index"] = docCount

	return b.index.Index(h, m)
}

func (b *Bleve) Entry(h string) error {
	exists, err := b.keyExists(h)
	if err != nil {
		return errors.Stack(err)
	}
	// If the hash exists, no need to index it again, we already did.
	if exists {
		return nil
	}

	docCount, err := b.index.DocCount()
	if err != nil {
		return errors.Stack(err)
	}

	// Increasing the doc count is important!
	// Due to go zero values and how kala is designed, zero values are ignored for
	// many things, including entries. So the 0th entry is the 1st index.
	docCount += 1

	m := index.Metadata{
		"index": docCount,
	}

	return b.index.Index(h, m)
}
