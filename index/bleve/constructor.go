package bleve

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/config"
)

const configType = "bleve"

func init() {
	fixity.RegisterIndex(configType, fixity.IndexConstructorFunc(Constructor))
}

func Constructor(n string, c config.Config) (fixity.Index, error) {
	return New(n, c)
}
