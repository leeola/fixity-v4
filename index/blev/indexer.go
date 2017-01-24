package blev

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

func (b *Bleve) Version(h string, v store.Version, m interface{}) error {
	return errors.New("not implemented")
}

func (b *Bleve) Meta(h string, m interface{}) error {
	return errors.New("not implemented")

	// exists, err := b.keyExists(h)
	// if err != nil {
	// 	return errors.Stack(err)
	// }
	// // If the hash exists, no need to index it again, we already did.
	// if exists {
	// 	return nil
	// }

	// docCount, err := b.entryIndex.DocCount()
	// if err != nil {
	// 	return errors.Stack(err)
	// }

	// // Increasing the doc count is important!
	// // Due to go zero values and how kala is designed, zero values are ignored for
	// // many things, including entries. So the 0th entry is the 1st index.
	// docCount += 1

	// m := index.Metadata{}
	// i.ToMetadata(m)
	// m["indexEntry"] = docCount

	// // if the metadata has an anchor set, index the metadata to the unique anchor
	// if ahI, ok := m["anchor"]; ok {
	// 	if ah, ok := ahI.(string); ok {
	// 		if err := b.indexUniqueAnchor(h, ah, m); err != nil {
	// 			return errors.Stack(err)
	// 		}
	// 	}
	// }

	// return b.entryIndex.Index(h, m)
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

	docCount, err := b.entryIndex.DocCount()
	if err != nil {
		return errors.Stack(err)
	}

	// Increasing the doc count is important!
	// Due to go zero values and how kala is designed, zero values are ignored for
	// many things, including entries. So the 0th entry is the 1st index.
	docCount += 1

	m := index.Metadata{
		"indexEntry": docCount,
	}

	return b.entryIndex.Index(h, m)
}
