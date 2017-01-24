package blev

import (
	"github.com/blevesearch/bleve/document"
	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
)

func (b *Bleve) keyExists(k string) (bool, error) {
	doc, err := b.entryIndex.Document(k)
	return doc != nil, err
}

func (b *Bleve) indexUniqueAnchor(h string, v store.Version, m map[string]interface{}) error {
	if v.Anchor == "" {
		return nil
	}

	doc, err := b.anchorIndex.Document(v.Anchor)
	if err != nil {
		return err
	}

	newTime := v.UploadedAt

	if doc != nil {
		for _, f := range doc.Fields {
			if f.Name() == "uploadedAt" {
				tf, ok := f.(*document.DateTimeField)
				if !ok {
					return errors.New("uploadedAt field is not valid DateTimeField")
				}

				previousTime, err := tf.DateTime()
				if err != nil {
					return errors.Wrap(err, "unable to get Time from DateTimeField")
				}

				// If the given metadata uploaded at time is not newer than the time
				// currently stored for the anchor, return and do nothing.
				//
				// If it is new, the new anchor metadata will be indexed below.
				if newTime.Before(previousTime) {
					return nil
				}

				break
			}
		}
	}

	// Copy the map, because we need to add a special metadata field to it so we
	// can retrieve the original hash.
	//
	// This is because the entry index uses the hash itself as the key, and so
	// a match is a query match can return the key/id directly. The unique index
	// uses the anchor hash, not the metadata hash, which means it cannot return
	// the key/id as that is not the hash the queryer wants.
	uniqueM := map[string]interface{}{}
	for k, v := range m {
		uniqueM[k] = v
	}

	// Now add our unique key. This value will be retrieved when a unique query is
	// done and returned to the user.
	uniqueM["_metaHash"] = h

	return errors.Stack(b.anchorIndex.Index(v.Anchor, uniqueM))
}
