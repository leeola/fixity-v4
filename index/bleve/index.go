package bleve

import (
	"errors"

	"github.com/leeola/fixity"
)

func (ix *Index) Index(ref fixity.Ref, m fixity.Mutation, v *fixity.Values, d *fixity.Data) error {
	return errors.New("not implemented")
}
