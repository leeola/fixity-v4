package contenttype

import (
	"time"

	"github.com/leeola/kala/store"
)

// Version is a type alias to implement FromChanges for store.Version.
type Version store.Version

func (ver *Version) FromChanges(c Changes) {
	if v, ok := c.GetAnchor(); ok {
		ver.Anchor = v
	}
	if v, ok := c.GetContentType(); ok {
		ver.ContentType = v
	}
	if v, ok := c.GetChangeLog(); ok {
		ver.ChangeLog = v
	}
	if v, ok := c.GetPreviousVersion(); ok {
		ver.PreviousVersion = v
	}
	if v, ok := c.GetMeta(); ok {
		ver.Meta = v
	}

	// TODO(leeola): enabled uploaded at setting. Optional! If not supplied,
	// automatic.
	// if v, ok := c.GetUploadedAt(); ok {
	// 	ver.UploadedAt = v
	// }
	ver.UploadedAt = time.Now()
}

// StoreType is a helper to cast the Version back to a store.Version.
func (v Version) StoreType() store.Version {
	return store.Version(v)
}
