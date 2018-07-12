package bleve

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/config"
)

const configType = "bleve"

func init() {
	fixity.RegisterIndex(configType, fixity.IndexCreatorFunc(Constructor))
}

func Constructor(n string, c config.Config) (fixity.QueryIndexer, error) {
	return New(n, c)
}
