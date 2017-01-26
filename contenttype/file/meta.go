package file

import ct "github.com/leeola/kala/contenttype"

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

	// apply common mutations to the various embedded structs.
	// these exist because some of the embedded values will be decided based
	// on other values, of which they may not be in scope. These methods
	// will do nothing as needed.
	m.Meta.NameFromFilename(m.Filename)
}

func (m *FileMeta) FromChanges(c ct.Changes) {
	if f, ok := c.GetString("filename"); ok {
		m.Filename = f
	}
}
