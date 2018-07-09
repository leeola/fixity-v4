package bleve

import (
	"errors"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/q"
)

func (ix *Index) Query(q.Query) ([]fixity.Ref, error) {
	return nil, errors.New("not implemented")
}
