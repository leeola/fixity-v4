package store

import (
	"net/url"
	"strconv"
	"time"
)

// MetaChanges is map of generic mutations to be handled by ContentType interfaces.
//
// It's generic because only a ContentType interface can know how to handle
// an addition of certain metadata. For example, filenames are metadata that
// is stored in the Meta blob, but is not a field on the actual store.Meta struct.
//
// Because of this, if you want to set a filename the actual ContentType implementor
// will need to set the filename. So, from the http api (or other future sources),
// will need to ferry this change through the generic ContentType interface.
//
// For a concrete example, the api request:
// 		POST /upload/file?filename=foo
//
// Might construct `metaChanges["filename"] = "foo"`, and then pass it to the
// `file` ContentType interface implementor. How it handles the actually setting
// and storing the filename within the metadata is up to it. This type just
// facilitates that data exchange, with some nicities for common change operations
// (like anchor and previousMeta changes).
//
// NOTE: The map[string]string type is used because the outer most interface
// typically accepts data in the form of strings, eg the http api or cli, and
// after that the only code that is informed about the actual type is the
// ContentType implementor.
//
// Eg, if you upload a file with a custom Size field, the http API has no idea the
// value of that field is a number. It would have to attempt to parse meaningful
// values at random, for no reason. The File ContentType will know the key Size is
// an int, and can convert as needed.
type MetaChanges map[string]string

func (c MetaChanges) GetString(k string) (string, bool) {
	s, ok := c[k]
	if !ok {
		return "", false
	}

	return s, true
}

func (c MetaChanges) GetBool(k string) (bool, bool) {
	s, ok := c[k]
	if !ok {
		return false, false
	}

	b, err := strconv.ParseBool(s)
	return b, err == nil
}

func (c MetaChanges) SetAnchor(h string) {
	c["anchor"] = h
}

func (c MetaChanges) SetMulti(h string) {
	c["multi"] = h
}

func (c MetaChanges) SetPreviousMeta(h string) {
	c["previousMeta"] = h
}

func (c MetaChanges) SetContentType(t string) {
	c["contentType"] = t
}

func (c MetaChanges) GetAnchor() (string, bool) {
	return c.GetString("anchor")
}

func (c MetaChanges) GetNewAnchor() (bool, bool) {
	return c.GetBool("newAnchor")
}

func (c MetaChanges) GetMulti() (string, bool) {
	return c.GetString("multi")
}

func (c MetaChanges) GetPreviousMeta() (string, bool) {
	return c.GetString("previousMeta")
}

func (c MetaChanges) GetContentType() (string, bool) {
	return c.GetString("contentType")
}

// Note: this func may not belong here (as it is specific to the http api),
// but it's here because it's important to keep the switch statement accurate
// to the state of this package. I'm open to moving this if a better solution
// presents itself.
func NewMetaChangesFromValues(m url.Values) MetaChanges {
	c := MetaChanges{}
	for k, v := range m {
		if len(v) == 0 {
			continue
		}
		switch k {
		case "anchor":
			c.SetAnchor(v[0])
		case "multi":
			c.SetMulti(v[0])
		case "previousMeta":
			c.SetPreviousMeta(v[0])
		case "contentType":
			c.SetContentType(v[0])
		default:
			c[k] = v[0]
		}
	}
	return c
}

func ApplyCommonChanges(m *Meta, c MetaChanges) {
	// Always set the timestamp
	m.UploadedAt = time.Now()

	if v, ok := c.GetAnchor(); ok {
		m.Anchor = v
	}

	if v, ok := c.GetMulti(); ok {
		m.Multi = v
	}

	if v, ok := c.GetPreviousMeta(); ok {
		m.PreviousMeta = v
	}

	if v, ok := c.GetContentType(); ok {
		m.ContentType = v
	}
}
