package file

import (
	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/index"
)

type Meta struct {
	ct.Meta
	FileMeta
}

type FileMeta struct {
	Filename string `json:"filename"`
}

func (m *Meta) FromChanges(c ct.Changes) {
	m.Meta.FromChanges(c)
	m.FileMeta.FromChanges(c)

	// If the name wasn't explicitly set (ie, to an empty value), and
	// the name is empty and the filename is not empty, default the name to
	// the filename.
	if _, ok := c.GetString("name"); !ok && m.Name == "" && m.Filename != "" {
		m.Name = m.Filename
	}
}

func (m *FileMeta) FromChanges(c ct.Changes) {
	if f, ok := c.GetString("filename"); ok {
		m.Filename = f
	}
}

func (m FileMeta) ToMetadata(im index.Metadata) {
	if m.Filename != "" {
		im["filename"] = m.Filename
	}
}
