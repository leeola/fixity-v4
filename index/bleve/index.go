package bleve

import (
	"fmt"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/index"
	"github.com/leeola/fixity/value"
)

func (ix *Index) Index(ref fixity.Ref, m fixity.Mutation, d *fixity.DataSchema, v fixity.Values) error {

	indexedValues := map[string]interface{}{}

	if v != nil {
		for k, v := range v {
			switch v.Type {
			case value.TypeInt:
				indexedValues[k] = v.IntValue
			case value.TypeString:
				indexedValues[k] = v.StringValue
			default:
				return fmt.Errorf("unhandled value type: %s", v.Type)
			}
		}
	}

	indexedValues[index.FIDKey] = m.ID
	indexedValues[index.FRefKey] = string(ref)
	if d != nil {
		indexedValues[index.FSizeKey] = d.Size
		indexedValues[index.FChecksumKey] = d.Checksum
	}

	if err := ix.idIndex.Index(m.ID, indexedValues); err != nil {
		return fmt.Errorf("bleve id index: %v", err)
	}

	if err := ix.refIndex.Index(string(ref), indexedValues); err != nil {
		return fmt.Errorf("bleve ref index: %v", err)
	}

	return nil
}
