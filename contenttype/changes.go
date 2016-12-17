package contenttype

import (
	"net/url"
	"strconv"
	"time"

	"github.com/leeola/kala/store"
)

// Changes is map of generic mutations to be handled by ContentType interfaces.
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
// 		POST /upload?contentType=file&filename=foo
//
// Might construct `changes["filename"] = "foo"`, and then pass it to the
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
type Changes map[string][]string

func (c Changes) GetString(k string) (string, bool) {
	v, ok := c[k]
	if !ok {
		return "", false
	}

	var s string
	if len(v) >= 1 {
		s = v[0]
	}

	return s, true
}

func (c Changes) GetBool(k string) (bool, bool) {
	s, ok := c.GetString(k)
	if !ok {
		return false, false
	}

	b, err := strconv.ParseBool(s)
	return b, err == nil
}

func (c Changes) Set(k, v string) {
	c[k] = []string{v}
}

func (c Changes) Add(k, v string) {
	vs, _ := c[k]
	c[k] = append(vs, v)
}

func (c Changes) SetAnchor(h string) {
	c.Set("anchor", h)
}

func (c Changes) SetMultiHash(h string) {
	c.Set("multiHash", h)
}

func (c Changes) SetMultiPart(h string) {
	c.Set("multiPart", h)
}

func (c Changes) SetPreviousMeta(h string) {
	c.Set("previousMeta", h)
}

func (c Changes) SetContentType(t string) {
	c.Set("contentType", t)
}

func (c Changes) GetAnchor() (string, bool) {
	return c.GetString("anchor")
}

func (c Changes) GetNewAnchor() (bool, bool) {
	return c.GetBool("newAnchor")
}

func (c Changes) GetMultiHash() (string, bool) {
	return c.GetString("multiHash")
}

func (c Changes) GetMultiPart() (string, bool) {
	return c.GetString("multiPart")
}

func (c Changes) GetPreviousMeta() (string, bool) {
	return c.GetString("previousMeta")
}

func (c Changes) GetContentType() (string, bool) {
	return c.GetString("contentType")
}

// Note: this func may not belong here (as it is specific to the http api),
// but it's here because it's important to keep the switch statement accurate
// to the state of this package. I'm open to moving this if a better solution
// presents itself.
func NewChangesFromValues(m url.Values) Changes {
	c := Changes{}
	for k, v := range m {
		if len(v) == 0 {
			continue
		}
		switch k {
		case "anchor":
			c.SetAnchor(v[0])
		case "multiHash":
			c.SetMultiHash(v[0])
		case "multiPart":
			c.SetMultiPart(v[0])
		case "previousMeta":
			c.SetPreviousMeta(v[0])
		case "contentType":
			c.SetContentType(v[0])
		default:
			c[k] = v
		}
	}
	return c
}

func ApplyCommonChanges(m *store.Meta, c Changes) {
	// Always set the timestamp
	m.UploadedAt = time.Now()

	if v, ok := c.GetContentType(); ok {
		m.ContentType = v
	}

	if v, ok := c.GetAnchor(); ok {
		m.Anchor = v
	}

	if v, ok := c.GetMultiHash(); ok {
		m.MultiHash = v
	}

	if v, ok := c.GetMultiPart(); ok {
		m.MultiPart = v
	}

	if v, ok := c.GetPreviousMeta(); ok {
		m.PreviousMeta = v
	}
}
