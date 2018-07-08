package bleve

import (
	"errors"

	"github.com/leeola/fixity/q"
)

func (ix *Index) Query(q.Query) error {
	return errors.New("not implemented")
}
