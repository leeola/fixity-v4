package blev

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

func (b *Bleve) Version(h string, v store.Version, m interface{}) error {
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

	// Convert the meta interface to a map so we can combine Version and the Meta.
	// This is done because we're indexing all values of Version and Meta to a single
	// hash.
	//
	// NOTE(leeola): Other indexers may do this differently, but this is the solution
	// i'm using for bleve.. improvements welcome.
	indexable := map[string]interface{}{}

	mapFields(indexable, structs.Fields(m))

	if _, ok := indexable["anchor"]; ok {
		return errors.New("Version and Meta field overlap: anchor")
	}
	indexable["anchor"] = v.Anchor
	if _, ok := indexable["contentType"]; ok {
		return errors.New("Version and Meta field overlap: contentType")
	}
	indexable["contentType"] = v.ContentType
	if _, ok := indexable["meta"]; ok {
		return errors.New("Version and Meta field overlap: meta")
	}
	indexable["meta"] = v.Meta
	if _, ok := indexable["previousVersion"]; ok {
		return errors.New("Version and Meta field overlap: previousVersion")
	}
	indexable["previousVersion"] = v.PreviousVersion
	if _, ok := indexable["previousVersionCount"]; ok {
		return errors.New("Version and Meta field overlap: previousVersionCount")
	}
	indexable["previousVersionCount"] = v.PreviousVersionCount
	if _, ok := indexable["uploadedAt"]; ok {
		return errors.New("Version and Meta field overlap: uploadedAt")
	}
	indexable["uploadedAt"] = v.UploadedAt

	// Store the docCount as indexEntry, this allows the user to query specific
	// entries. Very important!
	if _, ok := indexable["indexEntry"]; ok {
		return errors.New("reserved Meta field in use: indexEntry")
	}
	indexable["indexEntry"] = docCount

	// if the version has an anchor set, index the indexable to the unique anchor
	//
	// This will noop if there is no anchor in the version.
	if err := b.indexUniqueAnchor(h, v, indexable); err != nil {
		return errors.Stack(err)
	}

	return b.entryIndex.Index(h, indexable)
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

func mapFields(m map[string]interface{}, fs []*structs.Field) {
	for _, f := range fs {
		mapField(m, f)
	}
}

func mapField(m map[string]interface{}, f *structs.Field) {
	if f.IsEmbedded() {
		mapFields(m, f.Fields())
		return
	}

	// trim the omitempty suffix, as in:
	// `json:"foo,omitempty"` or `json:",omitempty"`
	key := strings.TrimSuffix(f.Tag("json"), ",omitempty")
	if key == "" {
		key = f.Name()
	}

	m[key] = f.Value()
}
