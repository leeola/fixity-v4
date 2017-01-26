package image

import (
	"path/filepath"

	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/contenttype/file"
)

type Meta struct {
	ct.Meta
	file.FileMeta
	ImageMeta
}

type ImageMeta struct {
	Format string `json:"format"`
}

func (m *Meta) FromChanges(c ct.Changes) {
	m.Meta.FromChanges(c)
	m.FileMeta.FromChanges(c)
	m.ImageMeta.FromChanges(c)

	// apply common mutations to the various embedded structs.
	// these exist because some of the embedded values will be decided based
	// on other values, of which they may not be in scope. These methods
	// will do nothing as needed.
	m.Meta.NameFromFilename(m.Filename)
	m.ImageMeta.FormatFromExt(m.Filename)
}

func (m *ImageMeta) FromChanges(c ct.Changes) {
	if f, ok := c.GetString("format"); ok {
		m.Format = f
	}
}

func (m *ImageMeta) FormatFromExt(filename string) {
	if m.Format == "" && filename != "" {
		m.Format = filepath.Ext(filename)

		// remove the dot from the extension
		if m.Format != "" {
			m.Format = m.Format[1:]
		}
	}
}
