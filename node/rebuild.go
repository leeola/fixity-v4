package node

import "github.com/leeola/errors"

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

	n.log.Warn("rebuilding disabled for refactoring!")
	return nil

	// ch, err := n.store.List()
	// if err != nil {
	// 	return errors.Stack(err)
	// }

	// for h := range ch {
	// 	ctype, b, err := store.GetContentTypeWithBytes(n.store, h)
	// 	if err != nil {
	// 		return errors.Wrapf(err, "failed to get contentType from hash:%s", h)
	// 	}

	// 	// If it has no type, index the entry
	// 	if ctype == "" {
	// 		if err := n.index.Entry(h); err != nil {
	// 			return errors.Wrapf(err, "failed to index entry:%s", h)
	// 		}
	// 		continue
	// 	}

	// 	mu, ok := n.metadataUnmarshallers[ctype]

	// 	if !ok {
	// 		if err := n.index.Entry(h); err != nil {
	// 			return errors.Wrapf(err, "failed to index entry:%s", h)
	// 		}
	// 		continue
	// 	}

	// 	m, err := mu.UnmarshalMetadata(b)
	// 	if err != nil {
	// 		return errors.Wrapf(err, "failed to unmarshal metadata from:%s", h)
	// 	}

	// 	if err := n.index.Metadata(h, m); err != nil {
	// 		return errors.Wrapf(err, "failed to index entry:%s", h)
	// 	}
	// }

	// // Now that we're done indexing, store the new index version in the node db.
	// if err := n.db.SetNodeIndexVersion(n.index.Version()); err != nil {
	// 	return errors.Wrap(err, "failed to store new index version")
	// }

	// return nil
}
