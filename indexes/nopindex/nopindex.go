package nopindex

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/q"
)

type Index struct{}

func New() *Index {
	return &Index{}
}

func (i *Index) Index(string, string, []fixity.Field) error {
	return nil
}

func (i *Index) Search(*q.Query) ([]string, error) {
	return nil, nil
}
