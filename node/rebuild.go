package node

import (
	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
)

func (n *Node) IsNewIndex() bool {
	// Ignoring the error, as it does not matter. Comparing an empty value
	// will mean we need to rebuild the index, and an error results in the same.
	v, _ := n.db.GetNodeIndexVersion()
	isOldIndex := v == n.index.IndexVersion()
	return !isOldIndex
}

func (n *Node) RebuildIndex() error {
	n.log.Info("rebuilding index")

	if err := n.index.Reset(); err != nil {
		return errors.Wrap(err, "failed to reset index")
	}

	ch, err := n.store.List()
	if err != nil {
		return errors.Stack(err)
	}

	for h := range ch {
		isVersion, v, _, err := store.IsVersionWithBytes(n.store, h)
		if err != nil {
			return errors.Wrapf(err, "failed to read version from hash: %s", h)
		}

		// If it's not a version struct, index just the entry - no metadata.
		if !isVersion {
			if err := n.index.Entry(h); err != nil {
				return errors.Wrapf(err, "failed to index entry: %s", h)
			}
			continue
		}

		metaB, err := store.ReadAll(n.store, v.Meta)
		if err != nil {
			return errors.Wrapf(err, "failed to read meta from hash: %s", h)
		}

		mu, ok := n.contentTypes[v.ContentType]

		if !ok {
			n.log.Warn("couldn't index content type", "type", v.ContentType)
			if err := n.index.Entry(h); err != nil {
				return errors.Wrapf(err, "failed to index entry:%s", h)
			}
			continue
		}

		m, err := mu.UnmarshalMeta(metaB)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal metadata from:%s", h)
		}

		if err := n.index.Version(h, v, m); err != nil {
			return errors.Wrapf(err, "failed to index entry:%s", h)
		}
	}

	// Now that we're done indexing, store the new index version in the node db.
	if err := n.db.SetNodeIndexVersion(n.index.IndexVersion()); err != nil {
		return errors.Wrap(err, "failed to store new index version")
	}

	return nil
}
