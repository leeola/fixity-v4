package index

import (
	"errors"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/config"
)

const configType = "bleve"

func init() {
	fixity.RegisterIndex(configType, fixity.IndexCreatorFunc(Constructor))
}

func Constructor(name string, c config.Config) (fixity.QueryIndexer, error) {
	return nil, errors.New("not implemented")
}
