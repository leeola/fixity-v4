package nosign

import (
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/config"
)

const configType = "nosign"

func init() {
	fixity.RegisterStore(configType, fixity.StoreCreatorFunc(Constructor))
}

func Constructor(name string, c config.Config) (fixity.Store, error) {
	return New(name, c)
}
