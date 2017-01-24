package contenttype

import (
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/util/strutil"
)

// Meta is a type alias to implement Indexable and From/To Changes methods.
//
// This mainly serves to move the method implementatons outside of the store, so
// that the store can remain low level as intended and not cross reference the
// packages where unneeded.
type Meta store.Meta

// FromChanges applies Changes to the Meta type as specified.
//
// BUG(leeola): only auto assign UploadedAt if not specified. This lets a user
// export/import multiple metas for the same anchor and preserve the original order
// of the metas. Before implementing this, make sure UploadedAt is purged when
// loading old meta to "mutate" it. Eg, if File contenttype loads the old meta
// but does not purge the uploadedAt, and then this method uses the defined
// uploadedAt, the *new* meta will have the same timestamp, which makes sorting
// the two entries impossible.
func (m *Meta) FromChanges(c Changes) {
	if v, ok := c.GetMultiHash(); ok {
		m.MultiHash = v
	}
	if v, ok := c.GetMultiPart(); ok {
		m.MultiPart = v
	}
	if v, ok := c.GetName(); ok {
		m.Name = v
	}
	if v, ok := c.GetTags(); ok {
		m.Tags = v
	}
	if v, ok := c["addTag"]; ok {
		for _, s := range v {
			m.Tags = strutil.AddUnique(m.Tags, s)
		}
	}
	if v, ok := c["delTag"]; ok {
		for _, s := range v {
			m.Tags = strutil.DelFromSlice(m.Tags, s)
		}
	}
}

// ToChanges implements ContentType.ToChanges(Changes)
func (m Meta) ToChanges(c Changes) {
	if m.MultiHash != "" {
		c.Set("multiHash", m.MultiHash)
	}
	if m.MultiPart != "" {
		c.Set("multiPart", m.MultiPart)
	}
	if m.Name != "" {
		c.Set("name", m.Name)
	}
	if m.Tags != nil {
		c["tags"] = m.Tags
	}
}

// ToMetadata implements index.ToMetadata()
func (m Meta) ToMetadata(im index.Metadata) {
	if m.MultiHash != "" {
		im["multiHash"] = m.MultiHash
	}
	if m.MultiPart != "" {
		im["multiPart"] = m.MultiPart
	}
	if m.Name != "" {
		im["name"] = m.Name
	}
	if m.Tags != nil {
		im["tags"] = m.Tags
	}
}
