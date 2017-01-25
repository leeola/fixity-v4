package contenttype

import (
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/kala/contenttype/ctutil"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/roller/camli"
)

// WriteMeta is a helper to store and index an interface.
func WriteMetaAndVersion(s store.Store, i index.Indexer, v Version, m interface{}) (string, string, error) {
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
	if err := i.Version(versionH, v.StoreType(), m); err != nil {
		return "", "", errors.Stack(err)
	}

	return metaH, versionH, nil
}

func ReadVersionFromChanges(s store.Store, q index.Queryer, c Changes) (Version, error) {
	anchor, _ := c.GetAnchor()
	previousVersion, _ := c.GetPreviousVersion()
	// If there is no previous version to base this mutation off of, then query the
	// indexer for the most recent hash for this anchor.
	if previousVersion == "" && anchor != "" {
		qu := index.Query{
			SearchVersions: false,
			Metadata: index.Metadata{
				// NOTE: Putting the hash in quotes because the querystring in bleve
				// has issues with a hyphenated hashstring. This is annoying, and
				// should be fixed somehow...
				"anchor": `"` + anchor + `"`,
			},
		}

		result, err := q.QueryOne(qu)
		if err != nil {
			return Version{}, errors.Stack(err)
		}

		if result.Hash.Hash != "" {
			c.SetPreviousVersion(result.Hash.Hash)
		}
	}

	var version Version
	// if a previous version is specified, read it from the store.
	// This could be from a the caller or queried from the anchor.
	if previousVersion, ok := c.GetPreviousVersion(); ok && previousVersion != "" {
		if err := store.ReadAndUnmarshal(s, previousVersion, &version); err != nil {
			return Version{}, errors.Stack(err)
		}
		version.PreviousVersionCount += 1
	}

	// modify the version struct with the fields defined in the changes map.
	// Note that this will include the PreviousVersion from above.
	version.FromChanges(c)

	if newAnchor, _ := c.GetNewAnchor(); newAnchor {
		a, err := store.NewAnchor()
		if err != nil {
			return Version{}, errors.Stack(err)
		}
		version.Anchor = a
	}

	// We might be returning a zero value version, and that's intended. It depends
	// on if the user supplied change info to locate the version or not.
	return version, nil
}

// WriteContent stores and indexes the given readcloser, returning the hashes.
//
// Note that this function returns both the multiPartHash and *all* of the hashes
// combined. This is because the multiPart hash address will almost always need
// to be recorded somewhere, and likewise the slice of hashes needs to be printed
// all combined.
//
// This function takes care of the normal use case. It is a helper, afterall.
func WriteContent(s store.Store, i index.Indexer, rc io.ReadCloser) (string, []string, error) {
	if rc == nil {
		return "", nil, errors.New("missing ReadCloser")
	}
	defer rc.Close()

	roller, err := camli.New(rc)
	if err != nil {
		return "", nil, errors.Stack(err)
	}

	// write the actual content
	hashes, err := ctutil.WritePartRoller(s, i, roller)
	if err != nil {
		return "", nil, errors.Stack(err)
	}

	// write the multipart
	multiPartHash, err := store.WriteMultiPart(s, store.MultiPart{
		Parts: hashes,
	})
	if err != nil {
		return "", nil, errors.Stack(err)
	}
	hashes = append(hashes, multiPartHash)

	// Write the entries, not including the final metadata hash
	// The last hash is metadata, and we'll add that manually.
	for _, h := range hashes {
		if err := i.Entry(h); err != nil {
			return "", nil, errors.Stack(err)
		}
	}

	return multiPartHash, hashes, nil
}
