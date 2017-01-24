package contenttype

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

// WriteMeta is a helper to store and index an interface.
func WriteMetaAndVersion(s store.Store, i index.Indexer, v store.Version, m interface{}) (string, string, error) {
	// First write the meta, and then fill in the
	metaH, err := store.MarshalAndWrite(s, m)
	if err != nil {
		return "", "", errors.Stack(err)
	}

	// TODO(leeola): Fill in any missing fields of the version. Eg, timestamp,
	// previous Version, etc.

	// Set the meta hash to the hash we just stored.
	v.Meta = metaH

	// Now write the meta as well.
	versionH, err := store.MarshalAndWrite(s, v)
	if err != nil {
		return "", "", errors.Stack(err)
	}

	// Pass the changes as metadata to the indexer.
	if err := i.Version(versionH, v, m); err != nil {
		return "", "", errors.Stack(err)
	}

	return metaH, versionH, nil
}

func VersionFromChanges(s store.Store, q index.Queryer, c Changes) (store.Version, error) {
	var version store.Version

	// TODO(leeola): enable anchor lookups.
	anchor, _ := c.GetAnchor()
	newAnchor, _ := c.GetNewAnchor()
	if anchor != "" || newAnchor {
		return store.Version{}, errors.New("anchor versioning is not implemented")
	}

	prevVer, ok := c.GetPreviousVersion()
	if ok {
		v, err := store.ReadVersion(s, prevVer)
		if err != nil {
			return store.Version{}, errors.Stack(err)
		}
		version = v
	}
	version.PreviousVersion = prevVer
	version.PreviousVersionCount += 1

	meta, ok := c.GetMeta()
	if ok {
		version.Meta = meta
	}

	// We might be returning a zero value version, and that's intended. It depends
	// on if the user supplied change info to locate the version or not.
	return store.Version{}, nil
}
